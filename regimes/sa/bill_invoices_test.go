package sa_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	_ "github.com/invopop/gobl/regimes/sa"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime: tax.WithRegime("SA"),
		Code:   "SAMPLE-001",
		Supplier: &org.Party{
			Name: "Acme Corp Saudi",
			TaxID: &tax.Identity{
				Country: "SA",
				Code:    "312345678912343",
			},
		},
		Customer: &org.Party{
			Name: "Sample Consumer LLC",
			TaxID: &tax.Identity{
				Country: "SA",
				Code:    "399999999900003",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Development services",
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "general",
					},
				},
			},
		},
	}
}

func calculatedInvoice(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := validInvoice()
	require.NoError(t, inv.Calculate())
	return inv
}

func TestValidInvoice(t *testing.T) {
	inv := calculatedInvoice(t)
	assert.NoError(t, rules.Validate(inv))
}

func TestInvoiceSupplierValidation(t *testing.T) {
	t.Run("missing supplier tax ID", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-SA-BILL-INVOICE-01]")
	})

	t.Run("nil supplier", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Supplier = nil
		assert.NotPanics(t, func() {
			_ = rules.Validate(inv)
		})
	})

	t.Run("supplier with empty tax ID code", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = &tax.Identity{
			Country: "SA",
			Code:    "",
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("supplier with invalid tax ID format", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = &tax.Identity{
			Country: "SA",
			Code:    "12345",
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "IDENTITY-01")
	})
}

func TestInvoiceSimplified(t *testing.T) {
	t.Run("simplified without customer", func(t *testing.T) {
		inv := validInvoice()
		inv.SetTags(tax.TagSimplified)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("simplified with customer", func(t *testing.T) {
		inv := validInvoice()
		inv.SetTags(tax.TagSimplified)
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})
}

func TestInvoiceReverseCharge(t *testing.T) {
	t.Run("reverse charge invoice", func(t *testing.T) {
		inv := validInvoice()
		inv.SetTags(tax.TagReverseCharge)
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})
}

func TestInvoiceCorrectionTypes(t *testing.T) {
	t.Run("credit note", func(t *testing.T) {
		inv := validInvoice()
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Preceding = []*org.DocumentRef{
			{Code: "INV/001"},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("debit note", func(t *testing.T) {
		inv := validInvoice()
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Preceding = []*org.DocumentRef{
			{Code: "INV/001"},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})
}

func TestInvoiceMultipleLines(t *testing.T) {
	t.Run("standard and zero-rated lines", func(t *testing.T) {
		inv := validInvoice()
		inv.Lines = append(inv.Lines, &bill.Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "Financial service",
				Price: num.NewAmount(1000, 2),
			},
			Taxes: tax.Set{
				{
					Category: "VAT",
					Percent:  num.NewPercentage(0, 2),
				},
			},
		})
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})
}

func TestInvoiceTaxCalculation(t *testing.T) {
	t.Run("standard VAT rate is 15%", func(t *testing.T) {
		inv := calculatedInvoice(t)
		require.NotNil(t, inv.Totals)
		require.NotNil(t, inv.Totals.Taxes)
		require.Len(t, inv.Totals.Taxes.Categories, 1)

		vatCat := inv.Totals.Taxes.Categories[0]
		assert.Equal(t, "VAT", vatCat.Code.String())
		require.Len(t, vatCat.Rates, 1)

		rate := vatCat.Rates[0]
		assert.Equal(t, "15.0%", rate.Percent.String())

		// 10 * 100.00 = 1000.00, 15% of 1000 = 150.00
		assert.Equal(t, "150.00", rate.Amount.String())
	})
}
