package cfdi_test

import (
	"strings"
	"testing"

	"github.com/invopop/gobl/addons/mx/cfdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	_ "github.com/invopop/gobl/regimes/mx"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Code:      "123",
		Currency:  "MXN",
		IssueDate: cal.MakeDate(2023, 1, 1),
		Tax: &bill.Tax{
			Addons: []cbc.Key{cfdi.KeyV4},
			Ext: tax.Extensions{
				cfdi.ExtKeyIssuePlace: "21000",
			},
		},
		Supplier: &org.Party{
			Name: "Test Supplier",
			Ext: tax.Extensions{
				cfdi.ExtKeyPostCode:     "21000",
				cfdi.ExtKeyFiscalRegime: "601",
			},
			TaxID: &tax.Identity{
				Country: "MX",
				Code:    "AAA010101AAA",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			Ext: tax.Extensions{
				cfdi.ExtKeyPostCode:     "65000",
				cfdi.ExtKeyFiscalRegime: "608",
				cfdi.ExtKeyUse:          "G01",
			},
			TaxID: &tax.Identity{
				Country: "MX",
				Code:    "ZZZ010101ZZZ",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "bogus",
					Price: num.MakeAmount(10000, 2),
					Unit:  org.UnitPackage,
					Ext: tax.Extensions{
						cfdi.ExtKeyProdServ: "01010101",
					},
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
					},
				},
			},
		},
	}
}

func TestValidInvoice(t *testing.T) {
	inv := validInvoice()
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestNormalizeInvoice(t *testing.T) {
	t.Run("no tax", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax = &bill.Tax{
			Addons: []cbc.Key{cfdi.KeyV4},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		require.NotNil(t, inv.Tax)
		assert.Equal(t, tax.ExtValue("21000"), inv.Tax.Ext[cfdi.ExtKeyIssuePlace])
	})
	t.Run("with supplier address code", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax = &bill.Tax{
			Addons: []cbc.Key{cfdi.KeyV4},
		}
		delete(inv.Supplier.Ext, cfdi.ExtKeyPostCode)
		inv.Supplier.Addresses = append(inv.Supplier.Addresses,
			&org.Address{
				Locality: "Mexico",
				Code:     "21000",
			},
		)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		require.NotNil(t, inv.Tax)
		assert.Equal(t, tax.ExtValue("21000"), inv.Tax.Ext[cfdi.ExtKeyIssuePlace])
	})
}

func TestCustomerValidation(t *testing.T) {
	inv := validInvoice()

	inv.Customer.TaxID = nil
	assertValidationError(t, inv, "customer: (tax_id: cannot be blank.)")

	inv.Customer = nil
	require.NoError(t, inv.Calculate())
	assert.NoError(t, inv.Validate())
}

func TestLineValidation(t *testing.T) {
	inv := validInvoice()

	inv.Lines[0].Quantity = num.MakeAmount(0, 0)
	assertValidationError(t, inv, "lines: (0: (quantity: must be greater than 0.).)")

	inv.Lines[0].Quantity = num.MakeAmount(-1, 0)
	assertValidationError(t, inv, "lines: (0: (quantity: must be greater than 0; total: must be no less than 0.).)")

	inv = validInvoice()

	inv.Lines[0].Item.Price = num.MakeAmount(-1, 0)
	assertValidationError(t, inv, "lines: (0: (total: must be no less than 0.).)")
}

func TestPaymentInstructionsValidation(t *testing.T) {
	inv := validInvoice()
	inv.Payment = &bill.Payment{
		Instructions: &pay.Instructions{},
	}

	inv.Payment.Instructions.Key = "direct-debit"
	assertValidationError(t, inv, "payment: (instructions: (ext: (mx-cfdi-payment-means: required.).).)")

	inv.Payment.Instructions.Key = "unexisting"
	assertValidationError(t, inv, "payment: (instructions: (key: must be or start with a valid key.).)")
}

func TestPaymentAdvancesValidation(t *testing.T) {
	inv := validInvoice()
	inv.Payment = &bill.Payment{
		Advances: []*pay.Advance{
			{
				Description: "A prepayment",
			},
		},
	}

	inv.Payment.Advances[0].Key = "direct-debit"
	assertValidationError(t, inv, "payment: (advances: (0: (ext: (mx-cfdi-payment-means: required.).).).)")

	inv.Payment.Advances[0].Key = "unexisting"
	assertValidationError(t, inv, "payment: (advances: (0: (key: must be or start with a valid key.).).)")

	inv.Payment.Advances[0].Key = ""
	assertValidationError(t, inv, "payment: (advances: (0: (ext: (mx-cfdi-payment-means: required.).).).)")
}

func TestPaymentTermsValidation(t *testing.T) {
	inv := validInvoice()
	inv.Payment = &bill.Payment{
		Terms: &pay.Terms{},
	}

	inv.Payment.Terms.Notes = strings.Repeat("x", 1001)
	assertValidationError(t, inv, "payment: (terms: (notes: the length must be no more than 1000.).)")

	inv.Payment.Terms.Notes = strings.Repeat("x", 1000)
	require.NoError(t, inv.Validate())
}

func TestUsoCFDIScenarioValidation(t *testing.T) {
	inv := validInvoice()

	inv.Customer.Ext = tax.Extensions{
		cfdi.ExtKeyFiscalRegime: "601",
		cfdi.ExtKeyPostCode:     "21000",
	}
	assertValidationError(t, inv, "ext: (mx-cfdi-use: required.)")
}

func TestPrecedingValidation(t *testing.T) {
	inv := validInvoice()

	inv.Preceding = []*bill.Preceding{
		{
			Code: "123",
			Stamps: []*head.Stamp{
				{
					Provider: "unexpected",
					Value:    "1234",
				},
			},
		},
	}
	assertValidationError(t, inv, "preceding: (0: (stamps: missing sat-uuid stamp.).); tax: (ext: (mx-cfdi-rel-type: required.).)")

	inv.Type = bill.InvoiceTypeCreditNote
	assertValidationError(t, inv, "preceding: (0: (stamps: missing sat-uuid stamp.).)")

	inv.Preceding[0].Stamps[0].Provider = "sat-uuid"
	require.NoError(t, inv.Validate())
}

func TestInvoiceDiscountValidation(t *testing.T) {
	inv := validInvoice()

	inv.Discounts = []*bill.Discount{
		{
			Percent: num.NewPercentage(20, 2),
		},
	}
	assertValidationError(t, inv, "discounts: the SAT doesn't allow discounts at invoice level")
}

func assertValidationError(t *testing.T, inv *bill.Invoice, expected string) {
	t.Helper()
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), expected)
}
