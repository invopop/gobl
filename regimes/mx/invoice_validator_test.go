package mx_test

import (
	"context"
	"strings"
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
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
			Tags: []cbc.Key{
				"use+goods-acquisition",
			},
		},
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: l10n.MX,
				Code:    "AAA010101AAA",
				Zone:    "21000",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
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
					Identities: []*org.Identity{
						{
							Type: "SAT",
							Code: "01010101",
						},
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
	ctx := context.Background()
	require.NoError(t, inv.Calculate(ctx))
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

	inv.Tax.Tags = make([]cbc.Key, 0)
	assertValidationError(t, inv, "'use' tax tags is required")

	inv.Tax.Tags = nil
	assertValidationError(t, inv, "'use' tax tags is required")

	inv.Tax = nil
	assertValidationError(t, inv, "'use' tax tags is required")
}

func assertValidationError(t *testing.T, inv *bill.Invoice, expected string) {
	ctx := context.Background()
	require.NoError(t, inv.Calculate(ctx))
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), expected)
}
