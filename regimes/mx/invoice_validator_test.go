package mx_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
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
	assertValidationError(t, inv, "lines: (0: (quantity: must be positive; total: must be positive.).)")

	inv.Lines[0].Quantity = num.MakeAmount(-1, 0)
	assertValidationError(t, inv, "lines: (0: (quantity: must be positive; total: must be positive.).)")

	inv = validInvoice()

	inv.Lines[0].Item.Price = num.MakeAmount(0, 0)
	assertValidationError(t, inv, "lines: (0: (total: must be positive.).)")

	inv.Lines[0].Item.Price = num.MakeAmount(-1, 0)
	assertValidationError(t, inv, "lines: (0: (total: must be positive.).)")

	inv = validInvoice()

	inv.Lines[0].Taxes = nil
	assertValidationError(t, inv, "lines: (0: (taxes: cannot be blank.).)")
}

func assertValidationError(t *testing.T, inv *bill.Invoice, expected string) {
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), expected)
}
