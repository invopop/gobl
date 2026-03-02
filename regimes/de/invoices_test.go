package de_test

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
		Regime: tax.WithRegime("DE"),
		Series: "TEST",
		Code:   "0002",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "111111125",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "282741168",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "bogus",
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

func TestInvoiceValidation(t *testing.T) {
	t.Run("normal invoice", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())

		inv = validInvoice()
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, inv.Validate(), "supplier: (identities: missing key 'de-tax-number'; tax_id: cannot be blank.).")
	})

	t.Run("simplified invoice - no tax details", func(t *testing.T) {
		inv := validInvoice()
		inv.SetTags("simplified")
		inv.Supplier.TaxID.Code = ""
		inv.Customer = nil

		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("regular invoice - only tax number", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID.Code = ""
		inv.Supplier.Identities = []*org.Identity{
			{
				Key:  "de-tax-number",
				Code: "92/345/67894",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})

	t.Run("regular invoice - only tax number nil tax ID", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = nil
		inv.Supplier.Identities = []*org.Identity{
			{
				Key:  "de-tax-number",
				Code: "92/345/67894",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})
}

// TestVATRatesByDate verifies that the correct VAT rate is applied at each
// historical boundary. Germany's rate transitions:
//
//	Standard: 15% (1993) → 16% (1998) → 19% (2007) → 16% (2020 COVID) → 19% (2021)
//	Reduced:  7% (1993) → 5% (2020 COVID) → 7% (2021)
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
		// --- Standard (general) rate ---

		// 1993 boundary: earliest rate (15%)
		{
			name:        "general: on 1993 Since date → 15%",
			issueDate:   cal.MakeDate(1993, 1, 1),
			rate:        "general",
			expectedTax: "15.00",
		},

		// 1998 boundary: 15% → 16%
		{
			name:        "general: day before 1998 Since → 15%",
			issueDate:   cal.MakeDate(1998, 3, 31),
			rate:        "general",
			expectedTax: "15.00",
		},
		{
			name:        "general: on 1998 Since date → 16%",
			issueDate:   cal.MakeDate(1998, 4, 1),
			rate:        "general",
			expectedTax: "16.00",
		},

		// 2007 boundary: 16% → 19%
		{
			name:        "general: day before 2007 Since → 16%",
			issueDate:   cal.MakeDate(2006, 12, 31),
			rate:        "general",
			expectedTax: "16.00",
		},
		{
			name:        "general: on 2007 Since date → 19%",
			issueDate:   cal.MakeDate(2007, 1, 1),
			rate:        "general",
			expectedTax: "19.00",
		},

		// 2020 COVID boundary: 19% → 16%
		{
			name:        "general: day before COVID Since → 19%",
			issueDate:   cal.MakeDate(2020, 6, 30),
			rate:        "general",
			expectedTax: "19.00",
		},
		{
			name:        "general: on COVID Since date → 16%",
			issueDate:   cal.MakeDate(2020, 7, 1),
			rate:        "general",
			expectedTax: "16.00",
		},

		// 2021 restoration boundary: 16% → 19%
		{
			name:        "general: last day of COVID → 16%",
			issueDate:   cal.MakeDate(2020, 12, 31),
			rate:        "general",
			expectedTax: "16.00",
		},
		{
			name:        "general: on restoration Since date → 19%",
			issueDate:   cal.MakeDate(2021, 1, 1),
			rate:        "general",
			expectedTax: "19.00",
		},
		{
			name:        "general: current era → 19%",
			issueDate:   cal.MakeDate(2025, 6, 15),
			rate:        "general",
			expectedTax: "19.00",
		},

		// --- Reduced rate ---

		// No 2007 transition — reduced stayed at 7% from 1993 until COVID
		{
			name:        "reduced: on 1993 Since date → 7%",
			issueDate:   cal.MakeDate(1993, 1, 1),
			rate:        "reduced",
			expectedTax: "7.00",
		},
		{
			name:        "reduced: pre-COVID → 7%",
			issueDate:   cal.MakeDate(2019, 6, 15),
			rate:        "reduced",
			expectedTax: "7.00",
		},

		// 2020 COVID boundary: 7% → 5%
		{
			name:        "reduced: day before COVID Since → 7%",
			issueDate:   cal.MakeDate(2020, 6, 30),
			rate:        "reduced",
			expectedTax: "7.00",
		},
		{
			name:        "reduced: on COVID Since date → 5%",
			issueDate:   cal.MakeDate(2020, 7, 1),
			rate:        "reduced",
			expectedTax: "5.00",
		},

		// 2021 restoration boundary: 5% → 7%
		{
			name:        "reduced: last day of COVID → 5%",
			issueDate:   cal.MakeDate(2020, 12, 31),
			rate:        "reduced",
			expectedTax: "5.00",
		},
		{
			name:        "reduced: on restoration Since date → 7%",
			issueDate:   cal.MakeDate(2021, 1, 1),
			rate:        "reduced",
			expectedTax: "7.00",
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
