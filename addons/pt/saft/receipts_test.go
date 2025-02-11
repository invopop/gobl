package saft_test

import (
	"testing"

	"github.com/invopop/gobl/addons/pt/saft"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func validReceipt() *bill.Receipt {
	return &bill.Receipt{
		Type: bill.ReceiptTypePayment,
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: "PT",
				Code:    "123456789",
			},
		},
		Customer: &org.Party{
			Name: "Customer Name",
			TaxID: &tax.Identity{
				Country: "PT",
				Code:    "987654321",
			},
		},
		Ext: tax.Extensions{
			saft.ExtKeyReceiptType: saft.ReceiptTypeOther,
		},
		Series:    "RG SERIES-A",
		Code:      "123",
		IssueDate: cal.MakeDate(2024, 3, 10),
		Lines: []*bill.ReceiptLine{
			{
				Document: &org.DocumentRef{
					IssueDate: cal.NewDate(2024, 3, 1),
				},
				Debit: num.NewAmount(100, 2),
				Tax: &tax.Total{
					Categories: []*tax.CategoryTotal{
						{
							Code: tax.CategoryVAT,
							Rates: []*tax.RateTotal{
								{
									Ext: tax.Extensions{
										pt.ExtKeyRegion:    "PT",
										saft.ExtKeyTaxRate: "NOR",
									},
								},
							},
						},
					},
				},
			},
		},
		Method: &pay.Instructions{
			Key: "credit-transfer",
		},
	}
}

func TestReceiptValidation(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)

	t.Run("valid receipt", func(t *testing.T) {
		rct := validReceipt()
		assert.NoError(t, addon.Validator(rct))
	})

	t.Run("invalid series", func(t *testing.T) {
		rct := validReceipt()

		rct.Series = "SERIES-A"
		assert.ErrorContains(t, addon.Validator(rct), "series: must start with 'RG '")
	})

	t.Run("invalid code", func(t *testing.T) {
		rct := validReceipt()

		rct.Code = "ABCD"
		assert.ErrorContains(t, addon.Validator(rct), "code: must be in a valid format")
	})

	t.Run("valid full code", func(t *testing.T) {
		rct := validReceipt()

		rct.Series = ""
		rct.Code = "RG SERIES-A/123"
		assert.NoError(t, addon.Validator(rct))
	})

	t.Run("missing extension", func(t *testing.T) {
		rct := validReceipt()
		rct.Ext = nil

		assert.ErrorContains(t, addon.Validator(rct), "ext: (pt-saft-receipt-type: required")
	})

	t.Run("missing supplier tax ID code", func(t *testing.T) {
		rct := validReceipt()
		rct.Supplier.TaxID.Code = cbc.CodeEmpty

		assert.ErrorContains(t, addon.Validator(rct), "supplier: (tax_id: (code: cannot be blank")

		rct.Supplier.TaxID = nil
		assert.ErrorContains(t, addon.Validator(rct), "supplier: (tax_id: cannot be blank.")
	})

	t.Run("missing customer name", func(t *testing.T) {
		rct := validReceipt()
		rct.Customer.Name = ""

		assert.ErrorContains(t, addon.Validator(rct), "customer: (name: cannot be blank")

		rct.Customer.TaxID.Code = ""
		assert.NoError(t, addon.Validator(rct))

		rct.Customer.TaxID = nil
		assert.NoError(t, addon.Validator(rct))
	})

	t.Run("missing line document", func(t *testing.T) {
		rct := validReceipt()
		rct.Lines[0].Document = nil

		assert.ErrorContains(t, addon.Validator(rct), "lines: (0: (document: cannot be blank")
	})

	t.Run("missing line document issue date", func(t *testing.T) {
		rct := validReceipt()
		rct.Lines[0].Document.IssueDate = nil

		assert.ErrorContains(t, addon.Validator(rct), "lines: (0: (document: (issue_date: cannot be blank")
	})

	t.Run("missing VAT category in line tax", func(t *testing.T) {
		rct := validReceipt()
		rct.Lines[0].Tax = nil

		assert.ErrorContains(t, addon.Validator(rct), "lines: (0: (tax: cannot be blank")

		rct.Lines[0].Tax = new(tax.Total)
		assert.ErrorContains(t, addon.Validator(rct), "lines: (0: (tax: missing category VAT")
	})

	t.Run("missing line tax required extensions", func(t *testing.T) {
		rct := validReceipt()
		rct.Lines[0].Tax.Categories[0].Rates[0].Ext = nil

		err := addon.Validator(rct)
		assert.ErrorContains(t, err, "pt-region: required")
		assert.ErrorContains(t, err, "pt-saft-tax-rate: required")
	})

	t.Run("negative amounts", func(t *testing.T) {
		rct := validReceipt()
		rct.Lines[0].Debit = num.NewAmount(-100, 2)

		assert.ErrorContains(t, addon.Validator(rct), "lines: (0: (debit: must be no less than 0")

		rct.Lines[0].Credit = num.NewAmount(-100, 2)
		assert.ErrorContains(t, addon.Validator(rct), "lines: (0: (credit: must be no less than 0")

		rct.Lines[0].Debit = &num.AmountZero
		rct.Lines[0].Credit = &num.AmountZero
		assert.NoError(t, addon.Validator(rct))
	})
}

func TestReceiptNormalization(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)

	t.Run("general", func(t *testing.T) {
		rct := validReceipt()
		rct.Ext = nil
		addon.Normalizer(rct)
		assert.Equal(t, "RG", rct.Ext[saft.ExtKeyReceiptType].String())
	})

	t.Run("VAT cash", func(t *testing.T) {
		rct := validReceipt()
		rct.SetTags("vat-cash")
		addon.Normalizer(rct)
		assert.Equal(t, "RC", rct.Ext[saft.ExtKeyReceiptType].String())
	})
}
