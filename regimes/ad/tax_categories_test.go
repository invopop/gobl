package ad_test

import (
	"testing"

	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategories(t *testing.T) {
	r := tax.RegimeDefFor("AD")
	require.NotNil(t, r)

	cat := r.CategoryDef(tax.CategoryVAT)
	require.NotNil(t, cat)

	t.Run("standard rate is 4.5%", func(t *testing.T) {
		rate := cat.RateDef(tax.KeyStandard, tax.RateGeneral)
		require.NotNil(t, rate)
		assert.Equal(t, "4.5%", rate.Values[0].Percent.String())
	})

	t.Run("reduced rate is 1%", func(t *testing.T) {
		rate := cat.RateDef("", tax.RateReduced)
		require.NotNil(t, rate)
		assert.Equal(t, "1%", rate.Values[0].Percent.String())
	})

	t.Run("intermediate rate is 2.5%", func(t *testing.T) {
		rate := cat.RateDef("", tax.RateIntermediate)
		require.NotNil(t, rate)
		assert.Equal(t, "2.5%", rate.Values[0].Percent.String())
	})

	t.Run("zero rate is 0%", func(t *testing.T) {
		rate := cat.RateDef(tax.KeyZero, tax.RateZero)
		require.NotNil(t, rate)
		assert.Equal(t, "0%", rate.Values[0].Percent.String())
	})

	t.Run("increased rate is 9.5%", func(t *testing.T) {
		rate := cat.RateDef("", "increased")
		require.NotNil(t, rate)
		assert.Equal(t, "9.5%", rate.Values[0].Percent.String())
	})
}
