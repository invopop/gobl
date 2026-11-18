package hr_test

import (
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/hr"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsumptionTax(t *testing.T) {
	t.Run("added on top of VAT, with its own percent", func(t *testing.T) {
		inv := testInvoice(t)
		inv.Lines[0].Taxes = append(inv.Lines[0].Taxes, &tax.Combo{
			Category: hr.TaxCategoryPP,
			Percent:  num.NewPercentage(3, 2),
		})
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))

		pp := inv.Totals.Taxes.Category(hr.TaxCategoryPP)
		require.NotNil(t, pp)
		require.Len(t, pp.Rates, 1)
		assert.Equal(t, "100.00", pp.Rates[0].Base.String())
		assert.Equal(t, "3%", pp.Rates[0].Percent.String())
		assert.Equal(t, "3.00", pp.Amount.String())
		assert.Equal(t, "25.00", inv.Totals.Taxes.Category(tax.CategoryVAT).Amount.String())
		assert.Equal(t, "128.00", inv.Totals.Payable.String())
	})

	t.Run("not retained", func(t *testing.T) {
		cat := hr.New().CategoryDef(hr.TaxCategoryPP)
		require.NotNil(t, cat)
		assert.False(t, cat.Retained)
		assert.Empty(t, cat.Rates)
	})
}
