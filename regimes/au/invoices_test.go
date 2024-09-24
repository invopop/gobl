package au_test

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
		Series:    "TEST",
		Code:      "TEST",
		Type:      bill.InvoiceTypeStandard,
		IssueDate: cal.MakeDate(2024, 1, 1),

		// Addresses and contact information are optional in Australia
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: l10n.AU.Tax(),
				Code:    "12004044937",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: l10n.AU.Tax(),
				Code:    "53004085616",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "example",
					Price: num.MakeAmount(500, 0),
					Unit:  org.UnitPackage,
				},
				Taxes: tax.Set{
					{
						Category: "GST",
						Rate:     "standard",
					},
				},
			},
		},
	}
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("Valid Invoice", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("Empty Supplier ID", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID.Code = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, inv.Validate(), "supplier: (tax_id: (code: invalid format.).).")
	})
	t.Run("Empty Supplier Name", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Name = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, inv.Validate(), "supplier: (name: cannot be blank.).")
	})

	t.Run("Empty Customer - Under 1000", func(t *testing.T) {
		inv := validInvoice()
		inv.Customer.TaxID = nil
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("Empty Customer - Over 1000", func(t *testing.T) {
		inv := validInvoice()
		inv.Customer.TaxID = nil
		inv.Lines[0].Item.Price = num.MakeAmount(50000, 0)
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, inv.Validate(), "customer: (tax_id: cannot be blank.).")
	})
}
