package mx_test

import (
	"strings"
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/mx"
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
			Ext: tax.Extensions{
				mx.ExtKeyCFDIIssuePlace: "21000",
			},
		},
		Supplier: &org.Party{
			Name: "Test Supplier",
			Ext: tax.Extensions{
				mx.ExtKeyCFDIPostCode:     "21000",
				mx.ExtKeyCFDIFiscalRegime: "601",
			},
			TaxID: &tax.Identity{
				Country: l10n.MX,
				Code:    "AAA010101AAA",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			Ext: tax.Extensions{
				mx.ExtKeyCFDIPostCode:     "65000",
				mx.ExtKeyCFDIFiscalRegime: "608",
				mx.ExtKeyCFDIUse:          "G01",
			},
			TaxID: &tax.Identity{
				Country: l10n.MX,
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
						mx.ExtKeyCFDIProdServ: "01010101",
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
		inv.Tax = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		require.NotNil(t, inv.Tax)
		assert.Equal(t, tax.ExtValue("21000"), inv.Tax.Ext[mx.ExtKeyCFDIIssuePlace])
	})
	t.Run("with supplier address code", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax = nil
		delete(inv.Supplier.Ext, mx.ExtKeyCFDIPostCode)
		inv.Supplier.Addresses = append(inv.Supplier.Addresses,
			&org.Address{
				Locality: "Mexico",
				Code:     "21000",
			},
		)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
		require.NotNil(t, inv.Tax)
		assert.Equal(t, tax.ExtValue("21000"), inv.Tax.Ext[mx.ExtKeyCFDIIssuePlace])
	})
}

func TestCustomerValidation(t *testing.T) {
	inv := validInvoice()

	inv.Customer.TaxID = nil
	assertValidationError(t, inv, "customer: (tax_id: cannot be blank.)")

	inv.Customer = nil
	assertValidationError(t, inv, "customer: cannot be blank")
}

func TestLineValidation(t *testing.T) {
	inv := validInvoice()

	inv.Lines[0].Quantity = num.MakeAmount(0, 0)
	assertValidationError(t, inv, "lines: (0: (quantity: must be greater than 0; total: must be greater than 0.).)")

	inv.Lines[0].Quantity = num.MakeAmount(-1, 0)
	assertValidationError(t, inv, "lines: (0: (quantity: must be greater than 0; total: must be greater than 0.).)")

	inv = validInvoice()

	inv.Lines[0].Item.Price = num.MakeAmount(0, 0)
	assertValidationError(t, inv, "lines: (0: (total: must be greater than 0.).)")

	inv.Lines[0].Item.Price = num.MakeAmount(-1, 0)
	assertValidationError(t, inv, "lines: (0: (total: must be greater than 0.).)")

	inv = validInvoice()

	inv.Lines[0].Taxes = nil
	assertValidationError(t, inv, "lines: (0: (taxes: cannot be blank.).)")
}

func TestPaymentInstructionsValidation(t *testing.T) {
	inv := validInvoice()
	inv.Payment = &bill.Payment{
		Instructions: &pay.Instructions{},
	}

	inv.Payment.Instructions.Key = "direct-debit"
	assertValidationError(t, inv, "payment: (instructions: (key: must be a valid value.).)")

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
	assertValidationError(t, inv, "payment: (advances: (0: (key: must be a valid value.).).)")

	inv.Payment.Advances[0].Key = "unexisting"
	assertValidationError(t, inv, "payment: (advances: (0: (key: must be or start with a valid key.).).)")

	inv.Payment.Advances[0].Key = ""
	assertValidationError(t, inv, "payment: (advances: (0: (key: cannot be blank.).).)")
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
		mx.ExtKeyCFDIFiscalRegime: "601",
		mx.ExtKeyCFDIPostCode:     "21000",
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
	assertValidationError(t, inv, "preceding: cannot be mapped from a `standard` type invoice")

	inv.Type = bill.InvoiceTypeCreditNote
	assertValidationError(t, inv, "preceding: (0: must have a `sat-uuid` stamp.)")

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
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), expected)
}
