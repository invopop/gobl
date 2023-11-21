package mx_test

import (
	"strings"
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
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
		Supplier: &org.Party{
			Name: "Test Supplier",
			Ext: cbc.CodeMap{
				mx.ExtKeyCFDIFiscalRegime: "601",
			},
			TaxID: &tax.Identity{
				Country: l10n.MX,
				Code:    "AAA010101AAA",
				Zone:    "21000",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			Ext: cbc.CodeMap{
				mx.ExtKeyCFDIFiscalRegime: "608",
				mx.ExtKeyCFDIUse:          "G01",
			},
			TaxID: &tax.Identity{
				Country: l10n.MX,
				Code:    "ZZZ010101ZZZ",
				Zone:    "65000",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "bogus",
					Price: num.MakeAmount(10000, 2),
					Unit:  "mutual",
					Ext: cbc.CodeMap{
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
		Payment: &bill.Payment{
			Instructions: &pay.Instructions{
				Key: "online+wallet",
			},
		},
	}
}

func TestValidInvoice(t *testing.T) {
	inv := validInvoice()
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
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

	inv.Payment.Instructions.Key = "direct-debit"
	assertValidationError(t, inv, "payment: (instructions: (key: must be a valid value.).)")

	inv.Payment.Instructions.Key = "unexisting"
	assertValidationError(t, inv, "payment: (instructions: (key: must be or start with a valid key.).)")

	inv.Payment.Instructions.Key = ""
	assertValidationError(t, inv, "payment: (instructions: (key: cannot be blank.).)")

	inv.Payment.Instructions = nil
	assertValidationError(t, inv, "payment: (instructions: cannot be blank.)")

	inv.Payment = nil
	assertValidationError(t, inv, "payment: cannot be blank")
}

func TestPaymentTermsValidation(t *testing.T) {
	inv := validInvoice()

	inv.Payment.Terms = &pay.Terms{}

	inv.Payment.Terms.Notes = strings.Repeat("x", 1001)
	assertValidationError(t, inv, "payment: (terms: (notes: the length must be no more than 1000.).)")

	inv.Payment.Terms.Notes = strings.Repeat("x", 1000)
	require.NoError(t, inv.Validate())
}

func TestUsoCFDIScenarioValidation(t *testing.T) {
	inv := validInvoice()

	inv.Customer.Ext = cbc.CodeMap{
		mx.ExtKeyCFDIFiscalRegime: "601",
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

func TestMabeValidation(t *testing.T) {
	inv := validInvoice()

	inv.Supplier.Ext[mx.ExtKeyMabeProviderCode] = "12345" // Marks it as a Mabe supplier

	assertValidationError(t, inv, "delivery: cannot be blank")

	inv.Delivery = &bill.Delivery{}

	assertValidationError(t, inv, "delivery: (receiver: cannot be blank")

	inv.Delivery.Receiver = &org.Party{
		Name: "Test Receiver",
	}

	assertValidationError(t, inv, "delivery: (receiver: (ext: (mx-mabe-delivery-plant: required")

	inv.Delivery.Receiver.Ext = cbc.CodeMap{
		mx.ExtKeyMabeDeliveryPlant: "S001",
	}

	assertValidationError(t, inv, "lines: (0: (item: (ext: (mx-mabe-item-code: required")

	inv.Lines[0].Item.Ext[mx.ExtKeyMabeItemCode] = "12345"

	assertValidationError(t, inv, "ordering: cannot be blank")

	inv.Ordering = &bill.Ordering{}

	assertValidationError(t, inv, "ordering: (code: cannot be blank")

	inv.Ordering.Code = "12345"

	assertValidationError(t, inv, "ext: (mx-mabe-reference-1: required")

	inv.Ext = cbc.CodeMap{
		mx.ExtKeyMabeReference1: "12345",
	}

	require.NoError(t, inv.Validate())
}

func assertValidationError(t *testing.T, inv *bill.Invoice, expected string) {
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), expected)
}
