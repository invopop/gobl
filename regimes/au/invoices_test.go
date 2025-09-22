package au_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Series:   "TEST",
		Code:     "0002",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "AU",
				Code:    "51824753556",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "AU",
				Code:    "53004085616",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "bogus",
					Price: num.NewAmount(50000, 2),
					Unit:  org.UnitPackage,
				},
				Taxes: tax.Set{
					{
						Category: "GST",
						Rate:     "general",
					},
				},
			},
		},
	}
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("valid invoice under $1000", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("valid invoice over $1000 with customer ABN", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines[0].Item.Price = num.NewAmount(120000, 2)
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("valid invoice over $1000 with customer name only", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines[0].Item.Price = num.NewAmount(120000, 2)
		inv.Customer.TaxID = nil
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("invalid invoice over $1000 without customer identity", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines[0].Item.Price = num.NewAmount(120000, 2)
		inv.Customer.Name = ""
		inv.Customer.TaxID = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, inv.Validate(), "customer identity or ABN required for invoices $1,000 AUD or more")
	})

	t.Run("missing supplier ABN", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID.Code = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, inv.Validate(), "supplier: (tax_id: (code: cannot be blank.).)")
	})
}
