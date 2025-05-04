package se_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/se"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	i := &bill.Invoice{
		Series:   "TEST",
		Code:     "0002",
		Currency: currency.SEK,
		Supplier: &org.Party{
			Name: "Test Supplier",
			Addresses: []*org.Address{
				{
					Country:  l10n.SE.ISO(),
					Code:     "12345",
					State:    "Stockholm",
					Locality: "Stockholm",
					Street:   "Test Street",
					Number:   "12",
				},
			},
			Identities: []*org.Identity{
				{
					Label:   "Org Number",
					Type:    se.IdentityTypeOrgNr,
					Country: l10n.SE.ISO(),
					Code:    "5560360793",
				},
			},
			TaxID: &tax.Identity{
				Country: l10n.TaxCountryCode(l10n.SE),
				Code:    "556036079301",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			Addresses: []*org.Address{
				{
					Country:  l10n.SE.ISO(),
					Code:     "54321",
					State:    "Stockholm",
					Locality: "Stockholm",
					Street:   "Test Street",
					Number:   "34",
				},
			},
			Identities: []*org.Identity{
				{
					Label:   "Personal Number",
					Type:    se.IdentityTypePersonNr,
					Country: l10n.SE.ISO(),
					Code:    "800101-0019",
				},
			},
			TaxID: &tax.Identity{
				Country: l10n.TaxCountryCode(l10n.SE),
				Code:    "800101001901",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Software Engineering Services",
					Price: num.NewAmount(110110, 2),
					Unit:  org.UnitHour,
				},
			},
		},
	}
	return i
}

func testInvoiceSimplified(t *testing.T) *bill.Invoice {
	t.Helper()
	i := &bill.Invoice{
		Series:   "TEST",
		Code:     "0002",
		Currency: currency.SEK,
		Supplier: &org.Party{
			// This is required due to bill.Invoice.validateSupplier only.
			// In Sweden, simplified invoices only require a supplier tax ID.
			Name: "Simplified Supplier",
			TaxID: &tax.Identity{
				Country: l10n.TaxCountryCode(l10n.SE),
				Code:    "556036079301",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Software Engineering Services",
					Price: num.NewAmount(40000, 2),
					Unit:  org.UnitHour,
				},
			},
		},
		Tags: tax.Tags{
			List: []cbc.Key{
				tax.TagSimplified,
			},
		},
	}
	return i
}

func TestInvoiceValidation(t *testing.T) {
	t.Parallel()
	t.Run("standard invoice", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("missing supplier tax ID", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.SetRegime("SE")
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (tax_id: cannot be blank.)")
	})

	t.Run("missing customer tax ID", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.SetRegime("SE")
		inv.Customer.TaxID = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("missing supplier name", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.SetRegime("SE")
		inv.Supplier.Name = ""
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (name: cannot be blank.)")
	})

	t.Run("missing customer name", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.SetRegime("SE")
		inv.Customer.Name = ""
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: (name: cannot be blank.)")
	})

	t.Run("missing supplier address", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.SetRegime("SE")
		inv.Supplier.Addresses = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (addresses: cannot be blank.)")
	})

	t.Run("missing customer address", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.SetRegime("SE")
		inv.Customer.Addresses = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: (addresses: cannot be blank.)")
	})

	t.Run("missing supplier identity", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.SetRegime("SE")
		inv.Supplier.Identities = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("missing customer identity", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.SetRegime("SE")
		inv.Customer.Identities = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("simplified invoice", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceSimplified(t)
		inv.SetRegime("SE")
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})
}
