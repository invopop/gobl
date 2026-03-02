package co_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime:    tax.WithRegime("CO"),
		Series:    "TEST",
		Code:      "0001",
		IssueDate: cal.MakeDate(2024, 6, 15),
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "CO",
				Code:    "412615332",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "CO",
				Code:    "8110079918",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(10000, 2),
					Unit:  org.UnitPackage,
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

// TestVATRatesByDate verifies that the correct VAT rate is applied at each
// historical boundary. Colombia has a single transition for the general rate:
//
//	General: 16% (2006) → 19% (2017)
//	Reduced: 5% (2006, static)
//
// Since comparison is inclusive: an invoice issued on the Since date gets
// the new rate.
func TestVATRatesByDate(t *testing.T) {
	tests := []struct {
		name        string
		issueDate   cal.Date
		rate        cbc.Key
		expectedTax string
	}{
		// --- General rate ---

		// 2006 boundary: earliest rate (16%)
		{
			name:        "general: on 2006 Since date → 16%",
			issueDate:   cal.MakeDate(2006, 1, 1),
			rate:        "general",
			expectedTax: "16.00",
		},
		{
			name:        "general: well into 2006 era → 16%",
			issueDate:   cal.MakeDate(2015, 6, 15),
			rate:        "general",
			expectedTax: "16.00",
		},

		// 2017 boundary: 16% → 19%
		{
			name:        "general: day before 2017 Since → 16%",
			issueDate:   cal.MakeDate(2016, 12, 31),
			rate:        "general",
			expectedTax: "16.00",
		},
		{
			name:        "general: on 2017 Since date → 19%",
			issueDate:   cal.MakeDate(2017, 1, 1),
			rate:        "general",
			expectedTax: "19.00",
		},
		{
			name:        "general: day after 2017 Since → 19%",
			issueDate:   cal.MakeDate(2017, 1, 2),
			rate:        "general",
			expectedTax: "19.00",
		},
		{
			name:        "general: current era → 19%",
			issueDate:   cal.MakeDate(2025, 6, 15),
			rate:        "general",
			expectedTax: "19.00",
		},

		// --- Reduced rate (static, no transitions) ---
		{
			name:        "reduced: on 2006 Since date → 5%",
			issueDate:   cal.MakeDate(2006, 1, 1),
			rate:        "reduced",
			expectedTax: "5.00",
		},
		{
			name:        "reduced: current era → 5%",
			issueDate:   cal.MakeDate(2025, 6, 15),
			rate:        "reduced",
			expectedTax: "5.00",
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
