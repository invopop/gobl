package au_test

import (
	"testing"
	"time"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/au"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime:    tax.WithRegime("AU"),
		Series:    "2026",
		Code:      "AU0001",
		IssueDate: cal.MakeDate(2026, time.April, 3),
		Currency:  currency.AUD,
		Supplier: &org.Party{
			Name: "Example Supplier Pty Ltd",
			TaxID: &tax.Identity{
				Country: l10n.AU.Tax(),
				Code:    "51824753556",
			},
			Addresses: []*org.Address{
				{
					Street:   "George Street",
					Number:   "100",
					Locality: "Sydney",
					State:    "NSW",
					Code:     "2000",
					Country:  l10n.AU.ISO(),
				},
			},
		},
		Customer: &org.Party{
			Name: "Example Customer Pty Ltd",
			Addresses: []*org.Address{
				{
					Street:   "Collins Street",
					Number:   "200",
					Locality: "Melbourne",
					State:    "VIC",
					Code:     "3000",
					Country:  l10n.AU.ISO(),
				},
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Software engineering services",
					Price: num.NewAmount(90000, 2),
					Unit:  org.UnitHour,
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryGST,
						Rate:     tax.RateGeneral,
					},
				},
			},
		},
	}
}

func TestInvoiceValidation(t *testing.T) {
	t.Parallel()

	t.Run("nil invoice", func(t *testing.T) {
		t.Parallel()
		var inv *bill.Invoice
		require.NoError(t, au.Validate(inv))
	})

	t.Run("valid invoice under threshold", func(t *testing.T) {
		t.Parallel()
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("valid invoice at threshold with customer name only", func(t *testing.T) {
		t.Parallel()
		inv := validInvoice()
		inv.Lines[0].Item.Price = num.NewAmount(100000, 2)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("invoice at threshold without customer", func(t *testing.T) {
		t.Parallel()
		inv := validInvoice()
		inv.Lines[0].Item.Price = num.NewAmount(100000, 2)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer")
	})

	t.Run("self billed invoice under threshold with customer name", func(t *testing.T) {
		t.Parallel()
		inv := validInvoice()
		inv.SetTags(tax.TagSelfBilled)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("self billed invoice under threshold without customer", func(t *testing.T) {
		t.Parallel()
		inv := validInvoice()
		inv.SetTags(tax.TagSelfBilled)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer")
	})

	t.Run("nil supplier", func(t *testing.T) {
		t.Parallel()
		inv := validInvoice()
		inv.Supplier = nil
		require.NoError(t, inv.Calculate())
		require.Error(t, inv.Validate())
	})

	t.Run("missing supplier tax ID", func(t *testing.T) {
		t.Parallel()
		inv := validInvoice()
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "tax_id")
	})

	t.Run("missing supplier name", func(t *testing.T) {
		t.Parallel()
		inv := validInvoice()
		inv.Supplier.Name = ""
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "name")
	})

	t.Run("supplier tax ID must be australian", func(t *testing.T) {
		t.Parallel()
		inv := validInvoice()
		inv.Supplier.TaxID.Country = "US"
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "country")
	})

	t.Run("supplier tax ID must be valid ABN", func(t *testing.T) {
		t.Parallel()
		inv := validInvoice()
		inv.Supplier.TaxID.Code = "11111111111"
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "invalid checksum")
	})
}
