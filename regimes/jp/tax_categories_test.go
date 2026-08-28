package jp_test

import (
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/regimes/jp"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsumptionTaxRates(t *testing.T) {
	cat := jp.New().CategoryDef(tax.CategoryVAT)
	require.NotNil(t, cat, "VAT category should be defined in the JP regime")

	t.Run("standard rate history", func(t *testing.T) {
		rate := cat.RateDef(tax.KeyStandard, tax.RateGeneral)
		require.NotNil(t, rate)

		tests := []struct {
			date    cal.Date
			percent string
		}{
			{cal.MakeDate(2024, 1, 1), "10%"},
			{cal.MakeDate(2019, 10, 1), "10%"}, // effective-from boundary
			{cal.MakeDate(2019, 9, 30), "8%"},  // day before the 10% rate
			{cal.MakeDate(2015, 6, 1), "8%"},
			{cal.MakeDate(2000, 1, 1), "5%"},
			{cal.MakeDate(1990, 1, 1), "3%"},
		}
		for _, tt := range tests {
			t.Run(tt.date.String(), func(t *testing.T) {
				v := rate.Value(tt.date, tax.Extensions{})
				require.NotNil(t, v, "expected a rate value for %s", tt.date)
				assert.Equal(t, tt.percent, v.Percent.String())
			})
		}

		t.Run("before the tax existed", func(t *testing.T) {
			assert.Nil(t, rate.Value(cal.MakeDate(1989, 3, 31), tax.Extensions{}))
		})
	})

	t.Run("reduced rate", func(t *testing.T) {
		rate := cat.RateDef(tax.KeyStandard, tax.RateReduced)
		require.NotNil(t, rate)

		v := rate.Value(cal.MakeDate(2024, 1, 1), tax.Extensions{})
		require.NotNil(t, v)
		assert.Equal(t, "8%", v.Percent.String())

		// The reduced rate was introduced with the 10% standard rate on
		// 2019-10-01 and did not exist before then.
		assert.Nil(t, rate.Value(cal.MakeDate(2019, 9, 30), tax.Extensions{}))
	})
}
