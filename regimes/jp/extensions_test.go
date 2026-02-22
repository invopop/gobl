package jp_test

import (
	"testing"

	_ "github.com/invopop/gobl" // ensure all regimes loaded
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/jp"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ptrAmt(val int64) *num.Amount {
	a := num.MakeAmount(val, 0)
	return &a
}

func TestExtensionsPresent(t *testing.T) {
	regime := jp.New()
	require.NotNil(t, regime)
	require.NotNil(t, regime.Extensions)
	assert.Len(t, regime.Extensions, 1, "should have one extension definition")

	ext := regime.Extensions[0]
	assert.Equal(t, cbc.Key(jp.ExtKeyReducedRate), ext.Key)
	assert.NotEmpty(t, ext.Name)
	assert.NotEmpty(t, ext.Desc)
	assert.Len(t, ext.Values, 2, "should have two extension codes")

	// Check food-beverage extension
	foodBev := ext.Values[0]
	assert.Equal(t, cbc.Code(jp.ExtCodeFoodBeverage), foodBev.Code)
	assert.NotEmpty(t, foodBev.Name)
	assert.NotEmpty(t, foodBev.Desc)

	// Check newspaper extension
	newspaper := ext.Values[1]
	assert.Equal(t, cbc.Code(jp.ExtCodeNewspaper), newspaper.Code)
	assert.NotEmpty(t, newspaper.Name)
	assert.NotEmpty(t, newspaper.Desc)
}

func TestInvoiceWithReducedRateExtensions(t *testing.T) {
	t.Run("valid food-beverage extension", func(t *testing.T) {
		inv := &bill.Invoice{
			Regime:    tax.WithRegime(l10n.TaxCountryCode(l10n.JP)),
			Code:      "TEST-001",
			Currency:  currency.JPY,
			IssueDate: *cal.NewDate(2024, 11, 1),
			Supplier: &org.Party{
				Name: "Test Supplier",
				TaxID: &tax.Identity{
					Country: l10n.TaxCountryCode(l10n.JP),
					Code:    "T7000012050002",
				},
			},
			Customer: &org.Party{
				Name: "Test Customer",
			},
			Lines: []*bill.Line{
				{
					Index:    1,
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Bento Box",
						Price: ptrAmt(500),
					},
					Taxes: []*tax.Combo{
						{
							Category: tax.CategoryVAT,
							Key:      tax.KeyStandard,
							Rate:     tax.RateReduced,
							Ext: tax.Extensions{
								jp.ExtKeyReducedRate: jp.ExtCodeFoodBeverage,
							},
						},
					},
				},
			},
		}

		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())

		// Verify the reduced rate (8%) was applied
		assert.Equal(t, num.MakePercentage(8, 2), *inv.Lines[0].Taxes[0].Percent)
	})

	t.Run("valid newspaper extension", func(t *testing.T) {
		inv := &bill.Invoice{
			Regime:    tax.WithRegime(l10n.TaxCountryCode(l10n.JP)),
			Code:      "TEST-002",
			Currency:  currency.JPY,
			IssueDate: *cal.NewDate(2024, 11, 1),
			Supplier: &org.Party{
				Name: "Test Supplier",
				TaxID: &tax.Identity{
					Country: l10n.TaxCountryCode(l10n.JP),
					Code:    "T7000012050002",
				},
			},
			Customer: &org.Party{
				Name: "Test Customer",
			},
			Lines: []*bill.Line{
				{
					Index:    1,
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Newspaper Subscription",
						Price: ptrAmt(4530),
					},
					Taxes: []*tax.Combo{
						{
							Category: tax.CategoryVAT,
							Key:      tax.KeyStandard,
							Rate:     tax.RateReduced,
							Ext: tax.Extensions{
								jp.ExtKeyReducedRate: jp.ExtCodeNewspaper,
							},
						},
					},
				},
			},
		}

		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())

		// Verify the reduced rate (8%) was applied
		assert.Equal(t, num.MakePercentage(8, 2), *inv.Lines[0].Taxes[0].Percent)
	})

	t.Run("reduced rate without extension is valid", func(t *testing.T) {
		inv := &bill.Invoice{
			Regime:    tax.WithRegime(l10n.TaxCountryCode(l10n.JP)),
			Code:      "TEST-003",
			Currency:  currency.JPY,
			IssueDate: *cal.NewDate(2024, 11, 1),
			Supplier: &org.Party{
				Name: "Test Supplier",
				TaxID: &tax.Identity{
					Country: l10n.TaxCountryCode(l10n.JP),
					Code:    "T7000012050002",
				},
			},
			Customer: &org.Party{
				Name: "Test Customer",
			},
			Lines: []*bill.Line{
				{
					Index:    1,
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Food Item",
						Price: ptrAmt(100),
					},
					Taxes: []*tax.Combo{
						{
							Category: tax.CategoryVAT,
							Key:      tax.KeyStandard,
							Rate:     tax.RateReduced,
							// No extension - should still be valid
						},
					},
				},
			},
		}

		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())

		// Verify the reduced rate (8%) was still applied
		assert.Equal(t, num.MakePercentage(8, 2), *inv.Lines[0].Taxes[0].Percent)
	})

	t.Run("standard rate with reduced extension should still use standard rate", func(t *testing.T) {
		inv := &bill.Invoice{
			Regime:    tax.WithRegime(l10n.TaxCountryCode(l10n.JP)),
			Code:      "TEST-004",
			Currency:  currency.JPY,
			IssueDate: *cal.NewDate(2024, 11, 1),
			Supplier: &org.Party{
				Name: "Test Supplier",
				TaxID: &tax.Identity{
					Country: l10n.TaxCountryCode(l10n.JP),
					Code:    "T7000012050002",
				},
			},
			Customer: &org.Party{
				Name: "Test Customer",
			},
			Lines: []*bill.Line{
				{
					Index:    1,
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Regular Item",
						Price: ptrAmt(1000),
					},
					Taxes: []*tax.Combo{
						{
							Category: tax.CategoryVAT,
							Key:      tax.KeyStandard,
							Rate:     tax.RateGeneral,
							Ext: tax.Extensions{
								jp.ExtKeyReducedRate: jp.ExtCodeFoodBeverage,
							},
						},
					},
				},
			},
		}

		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())

		// Extension doesn't change the rate - standard rate (10%) should apply
		assert.Equal(t, num.MakePercentage(10, 2), *inv.Lines[0].Taxes[0].Percent)
	})

	t.Run("invalid extension code", func(t *testing.T) {
		inv := &bill.Invoice{
			Regime:    tax.WithRegime(l10n.TaxCountryCode(l10n.JP)),
			Code:      "TEST-005",
			Currency:  currency.JPY,
			IssueDate: *cal.NewDate(2024, 11, 1),
			Supplier: &org.Party{
				Name: "Test Supplier",
				TaxID: &tax.Identity{
					Country: l10n.TaxCountryCode(l10n.JP),
					Code:    "T7000012050002",
				},
			},
			Customer: &org.Party{
				Name: "Test Customer",
			},
			Lines: []*bill.Line{
				{
					Index:    1,
					Quantity: num.MakeAmount(1, 0),
					Item: &org.Item{
						Name:  "Food Item",
						Price: ptrAmt(100),
					},
					Taxes: []*tax.Combo{
						{
							Category: tax.CategoryVAT,
							Key:      tax.KeyStandard,
							Rate:     tax.RateReduced,
							Ext: tax.Extensions{
								jp.ExtKeyReducedRate: "invalid-code",
							},
						},
					},
				},
			},
		}

		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "jp-reduced-rate")
	})
}

func TestMixedRateInvoice(t *testing.T) {
	inv := &bill.Invoice{
		Regime:    tax.WithRegime(l10n.TaxCountryCode(l10n.JP)),
		Code:      "TEST-006",
		Currency:  currency.JPY,
		IssueDate: *cal.NewDate(2024, 11, 1),
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: l10n.TaxCountryCode(l10n.JP),
				Code:    "T7000012050002",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
		},
		Lines: []*bill.Line{
			{
				Index:    1,
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Office Supplies",
					Price: ptrAmt(1000),
				},
				Taxes: []*tax.Combo{
					{
						Category: tax.CategoryVAT,
						Key:      tax.KeyStandard,
						Rate:     tax.RateGeneral,
					},
				},
			},
			{
				Index:    2,
				Quantity: num.MakeAmount(3, 0),
				Item: &org.Item{
					Name:  "Tea Box",
					Price: ptrAmt(500),
				},
				Taxes: []*tax.Combo{
					{
						Category: tax.CategoryVAT,
						Key:      tax.KeyStandard,
						Rate:     tax.RateReduced,
						Ext: tax.Extensions{
							jp.ExtKeyReducedRate: jp.ExtCodeFoodBeverage,
						},
					},
				},
			},
			{
				Index:    3,
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Newspaper Subscription",
					Price: ptrAmt(4530),
				},
				Taxes: []*tax.Combo{
					{
						Category: tax.CategoryVAT,
						Key:      tax.KeyStandard,
						Rate:     tax.RateReduced,
						Ext: tax.Extensions{
							jp.ExtKeyReducedRate: jp.ExtCodeNewspaper,
						},
					},
				},
			},
		},
	}

	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())

	// Verify standard rate (10%) on line 1
	assert.Equal(t, num.MakePercentage(10, 2), *inv.Lines[0].Taxes[0].Percent)

	// Verify reduced rate (8%) on line 2
	assert.Equal(t, num.MakePercentage(8, 2), *inv.Lines[1].Taxes[0].Percent)

	// Verify reduced rate (8%) on line 3
	assert.Equal(t, num.MakePercentage(8, 2), *inv.Lines[2].Taxes[0].Percent)

	// Verify totals exist
	assert.NotNil(t, inv.Totals)
	assert.NotNil(t, inv.Totals.Sum)
	assert.NotNil(t, inv.Totals.Tax)
	assert.NotNil(t, inv.Totals.TotalWithTax)
}
