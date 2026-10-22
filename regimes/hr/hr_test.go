package hr_test

import (
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/hr"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	regime := hr.New()
	require.NotNil(t, regime)
	require.NoError(t, rules.Validate(regime))
	assert.Equal(t, l10n.HR, regime.Country.Code())
	assert.Equal(t, currency.EUR, regime.Currency)
	assert.Equal(t, tax.CategoryVAT, regime.TaxScheme)
	assert.Equal(t, "Croatia", regime.Name.String())
	assert.Len(t, regime.Categories, 1)
	assert.Equal(t, "VAT", regime.Categories[0].Code.String())
	assert.Len(t, regime.Identities, 1)
	assert.Equal(t, "OIB", regime.Identities[0].Code.String())
}

func TestRateValues(t *testing.T) {
	cat := hr.New().Categories[0]
	tests := []struct {
		rate    cbc.Key
		date    cal.Date
		percent string
	}{
		{tax.RateGeneral, cal.MakeDate(2013, 7, 1), "25.0%"},
		{tax.RateGeneral, cal.MakeDate(2026, 1, 1), "25.0%"},
		{tax.RateReduced, cal.MakeDate(2013, 12, 31), "10.0%"},
		{tax.RateReduced, cal.MakeDate(2014, 1, 1), "13.0%"},
		{tax.RateSuperReduced, cal.MakeDate(2013, 7, 1), "5.0%"},
		{tax.RateSuperReduced, cal.MakeDate(2026, 1, 1), "5.0%"},
	}
	for _, tt := range tests {
		t.Run(tt.rate.String()+" "+tt.date.String(), func(t *testing.T) {
			rd := cat.RateDef(tax.KeyStandard, tt.rate)
			require.NotNil(t, rd)
			v := rd.Value(tt.date, tax.Extensions{})
			require.NotNil(t, v)
			assert.Equal(t, tt.percent, v.Percent.String())
		})
	}

	t.Run("no value before the current VAT Act", func(t *testing.T) {
		rd := cat.RateDef(tax.KeyStandard, tax.RateGeneral)
		require.NotNil(t, rd)
		assert.Nil(t, rd.Value(cal.MakeDate(2013, 6, 30), tax.Extensions{}))
	})
}
