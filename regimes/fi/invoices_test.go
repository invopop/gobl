package fi_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	_ "github.com/invopop/gobl/regimes/fi" // the regime needs to be registered
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Series:    "TEST",
		Code:      "0002",
		IssueDate: cal.MakeDate(2026, 3, 1),
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "FI",
				Code:    "12345671", // Valid Y-tunnus
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "FI",
				Code:    "07375462", // Valid Y-tunnus
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(100, 0),
					Unit:  org.UnitPackage,
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

func TestInvoiceValidation(t *testing.T) {
	inv := validInvoice()

	t.Run("valid invoice", func(t *testing.T) {
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("missing supplier tax ID code", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID.Code = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, inv.Validate(), "supplier: (tax_id: (code: cannot be blank.).).")
	})

}

func TestVATRatesByDate(t *testing.T) {
	tests := []struct {
		name        string
		issueDate   cal.Date
		rate        cbc.Key
		expectedTax string
	}{
		{
			name:        "standard current 25.5%",
			issueDate:   cal.MakeDate(2026, 1, 15),
			rate:        "standard",
			expectedTax: "25.50",
		},
		{
			name:        "standard pre-2024-09 24%",
			issueDate:   cal.MakeDate(2024, 8, 31),
			rate:        "standard",
			expectedTax: "24.00",
		},
		{
			name:        "reduced current 13.5%",
			issueDate:   cal.MakeDate(2026, 1, 15),
			rate:        "reduced",
			expectedTax: "13.50",
		},
		{
			name:        "reduced pre-2026 14%",
			issueDate:   cal.MakeDate(2025, 12, 31),
			rate:        "reduced",
			expectedTax: "14.00",
		},
		{
			name:        "super-reduced 10%",
			issueDate:   cal.MakeDate(2026, 1, 15),
			rate:        "super-reduced",
			expectedTax: "10.00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := validInvoice()
			inv.IssueDate = tt.issueDate
			inv.Lines[0].Taxes[0].Rate = tt.rate
			require.NoError(t, inv.Calculate())
			assert.Equal(t, tt.expectedTax, inv.Totals.Tax.String())
		})
	}
}
