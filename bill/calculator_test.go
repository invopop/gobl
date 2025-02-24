package bill_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: many calculation tests are distributed throughout this package.

func TestCalculate(t *testing.T) {
	t.Run("with round-then-sum rounding rule", func(t *testing.T) {
		inv := baseInvoice(t, &bill.Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "test item 1",
				Price: num.NewAmount(942, 2),
			},
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Percent:  num.NewPercentage(24, 2),
				},
			},
		}, &bill.Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "test item 2",
				Price: num.NewAmount(942, 2),
			},
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Percent:  num.NewPercentage(13, 2),
				},
			},
		})
		inv.Tax.PricesInclude = ""
		inv.Tax.Rounding = tax.RoundingRuleRoundThenSum
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "3.48", inv.Totals.Tax.String())

		inv.Tax.Rounding = tax.RoundingRuleSumThenRound
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "3.49", inv.Totals.Tax.String())
	})
}

func TestRemoveIncludedTaxes(t *testing.T) {
	t.Run("no included tax", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Tax = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.RemoveIncludedTaxes())
	})

	t.Run("from discounts", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Discounts = []*bill.Discount{
			{
				Amount: num.MakeAmount(1000, 2),
				Reason: "testing",
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     "standard",
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.RemoveIncludedTaxes())
		assert.Equal(t, "8.26", inv.Totals.Discount.String())
	})

	t.Run("from charges", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Charges = []*bill.Charge{
			{
				Amount: num.MakeAmount(1000, 2),
				Reason: "testing",
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     "standard",
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.RemoveIncludedTaxes())
		assert.Equal(t, "8.26", inv.Totals.Charge.String())
	})
}
