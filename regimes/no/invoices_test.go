package no_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
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

func TestInvoiceValidation(t *testing.T) {
	t.Parallel()

	t.Run("valid standard invoice", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("missing supplier tax ID", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (tax_id: cannot be blank.)")
	})

	t.Run("supplier tax ID with empty code", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID = &tax.Identity{Country: l10n.TaxCountryCode(l10n.NO), Code: ""}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "tax_id")
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

func TestInvoiceCalculation(t *testing.T) {
	t.Parallel()

	t.Run("standard invoice amounts", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())

		assert.Equal(t, "1500.00", inv.Totals.Sum.String())
		assert.Equal(t, "1500.00", inv.Totals.Total.String())
		assert.Equal(t, "375.00", inv.Totals.Tax.String())
		assert.Equal(t, "1875.00", inv.Totals.TotalWithTax.String())
		assert.Equal(t, "1875.00", inv.Totals.Payable.String())

		vat := inv.Totals.Taxes.Category(tax.CategoryVAT)
		require.NotNil(t, vat)
		assert.Equal(t, "375.00", vat.Amount.String())
		require.Len(t, vat.Rates, 1)
		assert.Equal(t, "25.0%", vat.Rates[0].Percent.String())
	})

	t.Run("prices include tax", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.Tax = &bill.Tax{
			PricesInclude: tax.CategoryVAT,
		}
		require.NoError(t, inv.Calculate())

		// Price 1500 NOK inclusive of 25% VAT
		// Net = 1500 / 1.25 = 1200
		// Tax = 1500 - 1200 = 300
		assert.Equal(t, "1500.00", inv.Totals.Sum.String())
		assert.Equal(t, "1200.00", inv.Totals.Total.String())
		assert.Equal(t, "300.00", inv.Totals.Tax.String())
		assert.Equal(t, "1500.00", inv.Totals.Payable.String())
	})

	t.Run("multi-rate invoice", func(t *testing.T) {
		t.Parallel()
		inv := testInvoiceStandard(t)
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(2, 0),
				Item: &org.Item{
					Name:  "Konsulenttjenester",
					Price: num.NewAmount(100000, 2),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     tax.RateGeneral,
					},
				},
			},
			{
				Quantity: num.MakeAmount(3, 0),
				Item: &org.Item{
					Name:  "Matlevering",
					Price: num.NewAmount(20000, 2),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     tax.RateReduced,
					},
				},
			},
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Togbillett",
					Price: num.NewAmount(50000, 2),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     tax.RateSuperReduced,
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())

		// Line 1: 2 × 1000 = 2000, VAT 25% = 500
		// Line 2: 3 × 200 = 600, VAT 15% = 90
		// Line 3: 1 × 500 = 500, VAT 12% = 60
		// Total: 3100, Tax: 650, Payable: 3750
		assert.Equal(t, "3100.00", inv.Totals.Sum.String())
		assert.Equal(t, "3100.00", inv.Totals.Total.String())
		assert.Equal(t, "650.00", inv.Totals.Tax.String())
		assert.Equal(t, "3750.00", inv.Totals.Payable.String())

		vat := inv.Totals.Taxes.Category(tax.CategoryVAT)
		require.NotNil(t, vat)
		require.Len(t, vat.Rates, 3)

		assert.Equal(t, "2000.00", vat.Rates[0].Base.String())
		assert.Equal(t, "25.0%", vat.Rates[0].Percent.String())
		assert.Equal(t, "500.00", vat.Rates[0].Amount.String())

		assert.Equal(t, "600.00", vat.Rates[1].Base.String())
		assert.Equal(t, "15.0%", vat.Rates[1].Percent.String())
		assert.Equal(t, "90.00", vat.Rates[1].Amount.String())

		assert.Equal(t, "500.00", vat.Rates[2].Base.String())
		assert.Equal(t, "12.0%", vat.Rates[2].Percent.String())
		assert.Equal(t, "60.00", vat.Rates[2].Amount.String())
	})
}
