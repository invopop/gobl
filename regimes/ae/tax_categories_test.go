package ae_test

import (
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/ae"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarginSchemeKey(t *testing.T) {
	t.Run("defined on VAT without a percent", func(t *testing.T) {
		kd := tax.RegimeDefFor("AE").CategoryDef(tax.CategoryVAT).KeyDef(ae.TaxKeyMarginScheme)
		require.NotNil(t, kd)
		assert.True(t, kd.NoPercent)
	})

	t.Run("the global VAT keys are unchanged", func(t *testing.T) {
		assert.Nil(t, tax.GlobalVAT().KeyDef(ae.TaxKeyMarginScheme))
	})

	t.Run("margin scheme lines have no VAT", func(t *testing.T) {
		inv := testInvoiceReverseCharge(t)
		inv.Tags = tax.Tags{}
		inv.Lines[0].Taxes = tax.Set{{Category: tax.CategoryVAT, Key: ae.TaxKeyMarginScheme}}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
		assert.Nil(t, inv.Lines[0].Taxes[0].Percent)
		assert.Equal(t, num.MakeAmount(0, 2).String(), inv.Totals.Tax.String())
	})
}
