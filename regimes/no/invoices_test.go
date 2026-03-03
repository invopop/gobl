package no_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Regime:   tax.WithRegime("NO"),
		Series:   "TEST",
		Code:     "0001",
		Currency: currency.NOK,
		Supplier: &org.Party{
			Name: "Eksempel AS",
			TaxID: &tax.Identity{
				Country: l10n.TaxCountryCode(l10n.NO),
				Code:    "923456783",
			},
			Addresses: []*org.Address{
				{
					Street:   "Eksempelveien 1",
					Locality: "Oslo",
					Code:     "0150",
					Country:  l10n.NO.ISO(),
				},
			},
		},
		Customer: &org.Party{
			Name: "Mottaker AS",
			TaxID: &tax.Identity{
				Country: l10n.TaxCountryCode(l10n.NO),
				Code:    "889640782",
			},
			Addresses: []*org.Address{
				{
					Street:   "Kundeveien 42",
					Locality: "Bergen",
					Code:     "5003",
					Country:  l10n.NO.ISO(),
				},
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Konsulenttjenester",
					Price: num.NewAmount(150000, 2),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     tax.RateGeneral,
					},
				},
			},
		},
	}
}

func testInvoiceSimplified(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Regime:   tax.WithRegime("NO"),
		Series:   "TEST",
		Code:     "0002",
		Currency: currency.NOK,
		Tags: tax.Tags{
			List: []cbc.Key{tax.TagSimplified},
		},
		Supplier: &org.Party{
			Name: "Eksempel AS",
			TaxID: &tax.Identity{
				Country: l10n.TaxCountryCode(l10n.NO),
				Code:    "923456783",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Dagligvarer",
					Price: num.NewAmount(50000, 2),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     tax.RateGeneral,
					},
				},
			},
		},
	}
}

func TestInvoiceValidation(t *testing.T) {
	t.Parallel()

	t.Run("valid standard invoice", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("missing supplier name", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.Supplier.Name = ""
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (name: cannot be blank.)")
	})

	t.Run("missing supplier address", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.Supplier.Addresses = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (addresses: cannot be blank.)")
	})

	t.Run("supplier without tax ID is valid", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("missing customer", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: cannot be blank.")
	})

	t.Run("missing customer name", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.Customer.Name = ""
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: (name: cannot be blank.)")
	})

	t.Run("customer without address is valid", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.Customer.Addresses = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("credit note without preceding", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeCreditNote
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "preceding: cannot be blank.")
	})

	t.Run("debit note without preceding", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeDebitNote
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "preceding: cannot be blank.")
	})

	t.Run("standard invoice with preceding is allowed", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.Preceding = []*org.DocumentRef{
			{Code: "TEST-0001"},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})
}

func TestSimplifiedInvoiceValidation(t *testing.T) {
	t.Parallel()

	t.Run("valid simplified invoice", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceSimplified(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("simplified invoice allows customer", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceSimplified(t)
		inv.Customer = &org.Party{
			Name: "Optional Kunde AS",
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("simplified invoice does not require supplier address", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceSimplified(t)
		inv.Supplier.Addresses = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("simplified invoice does not require supplier tax ID", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceSimplified(t)
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("simplified invoice still requires supplier name", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceSimplified(t)
		inv.Supplier.Name = ""
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (name: cannot be blank.)")
	})
}
