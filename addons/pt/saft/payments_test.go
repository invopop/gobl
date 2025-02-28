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

func validPayment() *bill.Payment {
	return &bill.Payment{
		Type: bill.PaymentTypeReceipt,
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
			saft.ExtKeyPaymentType: saft.PaymentTypeOther,
		},
		Series:    "RG SERIES-A",
		Code:      "123",
		IssueDate: cal.MakeDate(2024, 3, 10),
		Lines: []*bill.PaymentLine{
			{
				Document: &org.DocumentRef{
					IssueDate: cal.NewDate(2024, 3, 1),
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
				Debit: num.NewAmount(100, 2),
			},
		},
		Method: &pay.Instructions{
			Key: "credit-transfer",
		},
	}
}

func TestPaymentValidation(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)

	t.Run("valid payment", func(t *testing.T) {
		pmt := validPayment()
		assert.NoError(t, addon.Validator(pmt))
	})

	t.Run("invalid series", func(t *testing.T) {
		pmt := validPayment()

		pmt.Series = "SERIES-A"
		assert.ErrorContains(t, addon.Validator(pmt), "series: must start with 'RG '")
	})

	t.Run("invalid code", func(t *testing.T) {
		pmt := validPayment()

		pmt.Code = "ABCD"
		assert.ErrorContains(t, addon.Validator(pmt), "code: must be in a valid format")
	})

	t.Run("valid full code", func(t *testing.T) {
		pmt := validPayment()

		pmt.Series = ""
		pmt.Code = "RG SERIES-A/123"
		assert.NoError(t, addon.Validator(pmt))
	})

	t.Run("missing extension", func(t *testing.T) {
		pmt := validPayment()
		pmt.Ext = nil

		assert.ErrorContains(t, addon.Validator(pmt), "ext: (pt-saft-payment-type: required")
	})

	t.Run("missing supplier tax ID code", func(t *testing.T) {
		pmt := validPayment()
		pmt.Supplier.TaxID.Code = cbc.CodeEmpty

		assert.ErrorContains(t, addon.Validator(pmt), "supplier: (tax_id: (code: cannot be blank")

		pmt.Supplier.TaxID = nil
		assert.ErrorContains(t, addon.Validator(pmt), "supplier: (tax_id: cannot be blank.")

		pmt.Supplier = nil
		assert.NoError(t, addon.Validator(pmt))
	})

	t.Run("missing customer name", func(t *testing.T) {
		pmt := validPayment()
		pmt.Customer.Name = ""

		assert.ErrorContains(t, addon.Validator(pmt), "customer: (name: cannot be blank")

		pmt.Customer.TaxID.Code = ""
		assert.NoError(t, addon.Validator(pmt))

		pmt.Customer.TaxID = nil
		assert.NoError(t, addon.Validator(pmt))

		pmt.Customer = nil
		assert.NoError(t, addon.Validator(pmt))
	})

	t.Run("missing line document", func(t *testing.T) {
		pmt := validPayment()
		pmt.Lines[0].Document = nil

		assert.ErrorContains(t, addon.Validator(pmt), "lines: (0: (document: cannot be blank.).)")

		pmt.Lines[0] = nil
		assert.NoError(t, addon.Validator(pmt))
	})

	t.Run("missing line document issue date", func(t *testing.T) {
		pmt := validPayment()
		pmt.Lines[0].Document.IssueDate = nil

		assert.ErrorContains(t, addon.Validator(pmt), "lines: (0: (document: (issue_date: cannot be blank")
	})

	t.Run("missing VAT category in line tax", func(t *testing.T) {
		pmt := validPayment()
		pmt.Lines[0].Document.Tax = nil

		assert.ErrorContains(t, addon.Validator(pmt), "lines: (0: (document: (tax: cannot be blank.).).).")

		pmt.Lines[0].Document.Tax = new(tax.Total)
		assert.ErrorContains(t, addon.Validator(pmt), "lines: (0: (document: (tax: missing category VAT.).).).")
	})

	t.Run("missing line tax required extensions", func(t *testing.T) {
		pmt := validPayment()
		pmt.Lines[0].Document.Tax.Categories[0].Rates[0].Ext = nil

		err := addon.Validator(pmt)
		assert.ErrorContains(t, err, "pt-region: required")
		assert.ErrorContains(t, err, "pt-saft-tax-rate: required")

		pmt.Lines[0].Document.Tax.Categories[0].Rates[0] = nil
		assert.NoError(t, addon.Validator(pmt))
	})

	t.Run("negative amounts", func(t *testing.T) {
		pmt := validPayment()
		pmt.Lines[0].Debit = num.NewAmount(-100, 2)

		assert.ErrorContains(t, addon.Validator(pmt), "lines: (0: (debit: must be no less than 0")

		pmt.Lines[0].Credit = num.NewAmount(-100, 2)
		assert.ErrorContains(t, addon.Validator(pmt), "lines: (0: (credit: must be no less than 0")

		pmt.Lines[0].Debit = &num.AmountZero
		pmt.Lines[0].Credit = &num.AmountZero
		assert.NoError(t, addon.Validator(pmt))
	})
}

func TestPaymentNormalization(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)

	t.Run("general", func(t *testing.T) {
		pmt := validPayment()
		pmt.Ext = nil
		addon.Normalizer(pmt)
		assert.Equal(t, "RG", pmt.Ext[saft.ExtKeyPaymentType].String())
	})

	t.Run("VAT cash", func(t *testing.T) {
		pmt := validPayment()
		pmt.SetTags("vat-cash")
		addon.Normalizer(pmt)
		assert.Equal(t, "RC", pmt.Ext[saft.ExtKeyPaymentType].String())
	})
}
