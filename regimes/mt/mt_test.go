package mt_test

import (
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/mt"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	regime := mt.New()
	require.NotNil(t, regime)
	require.NoError(t, rules.Validate(regime))

	assert.Equal(t, l10n.MT, regime.Country.Code())
	assert.Equal(t, "Malta", regime.Name.String())
	assert.Equal(t, "EUR", regime.Currency.String())
	assert.Equal(t, "Europe/Malta", regime.TimeZone)

	require.Len(t, regime.Categories, 1)
	assert.Equal(t, "VAT", regime.Categories[0].Code.String())

	// Corrections and Identities are deliberately absent; see the notes in mt.go.
	assert.Empty(t, regime.Corrections)
	assert.Empty(t, regime.Identities)
}

func TestTaxCategoryKeys(t *testing.T) {
	// Every global VAT key is offered, matching every other VAT regime in the
	// codebase. Malta needs export and intra-community in particular: items 1 and 3
	// of Part One of the Fifth Schedule are exempt *with* credit, so flattening
	// them into zero would lose the distinction the VAT Act draws.
	cat := mt.New().Categories[0]
	keys := make([]string, 0, len(cat.Keys))
	for _, k := range cat.Keys {
		keys = append(keys, k.Key.String())
	}
	for _, want := range []string{
		"standard", "zero", "exempt", "export", "intra-community",
	} {
		assert.Contains(t, keys, want)
	}
	assert.Equal(t, len(tax.GlobalVATKeys()), len(cat.Keys),
		"the regime should offer the global VAT key set unchanged")
}

func TestTaxCategoryRates(t *testing.T) {
	cat := mt.New().Categories[0]
	require.Len(t, cat.Rates, 4)

	tests := []struct {
		rate    cbc.Key
		percent string
	}{
		{tax.RateGeneral, "18.0%"},
		{tax.RateIntermediate, "12.0%"},
		{tax.RateReduced, "7.0%"},
		{tax.RateSuperReduced, "5.0%"},
	}

	for _, tt := range tests {
		t.Run(tt.rate.String(), func(t *testing.T) {
			rate := cat.RateDef(tax.KeyStandard, tt.rate)
			require.NotNil(t, rate, "rate %s not defined", tt.rate)
			val := rate.Value(cal.MakeDate(2026, 1, 1), tax.Extensions{})
			require.NotNil(t, val)
			assert.Equal(t, tt.percent, val.Percent.String())
		})
	}
}

func TestIntermediateRateCommencement(t *testing.T) {
	// LN 231/2023 introduced the 12% rate with effect from 1 January 2024.
	cat := mt.New().Categories[0]

	intermediate := cat.RateDef(tax.KeyStandard, tax.RateIntermediate)
	require.NotNil(t, intermediate)
	assert.Nil(t, intermediate.Value(cal.MakeDate(2023, 12, 31), tax.Extensions{}),
		"the 12%% rate must not apply before it commenced")
	val := intermediate.Value(cal.MakeDate(2024, 1, 1), tax.Extensions{})
	require.NotNil(t, val)
	assert.Equal(t, "12.0%", val.Percent.String())

	general := cat.RateDef(tax.KeyStandard, tax.RateGeneral)
	require.NotNil(t, general)
	assert.NotNil(t, general.Value(cal.MakeDate(2023, 12, 31), tax.Extensions{}),
		"the standard rate is unaffected by the 12%% rate's commencement")
}

// TestRateHistory pins the dates from the European Commission's rate table, cited in the
// category sources.
func TestRateHistory(t *testing.T) {
	cat := mt.New().Categories[0]

	tests := []struct {
		name    string
		rate    cbc.Key
		date    cal.Date
		percent string // empty means the rate must not resolve
	}{
		{"general before the Act commenced", tax.RateGeneral, cal.MakeDate(1998, 12, 31), ""},
		{"general on commencement", tax.RateGeneral, cal.MakeDate(1999, 1, 1), "15.0%"},
		{"general day before the rise", tax.RateGeneral, cal.MakeDate(2003, 12, 31), "15.0%"},
		{"general on the rise to 18%", tax.RateGeneral, cal.MakeDate(2004, 1, 1), "18.0%"},
		{"general today", tax.RateGeneral, cal.MakeDate(2026, 1, 1), "18.0%"},
		{"reduced before 7% existed", tax.RateReduced, cal.MakeDate(2010, 12, 31), ""},
		{"reduced on commencement", tax.RateReduced, cal.MakeDate(2011, 1, 1), "7.0%"},
		{"super-reduced before the Act", tax.RateSuperReduced, cal.MakeDate(1998, 12, 31), ""},
		{"super-reduced on commencement", tax.RateSuperReduced, cal.MakeDate(1999, 1, 1), "5.0%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rate := cat.RateDef(tax.KeyStandard, tt.rate)
			require.NotNil(t, rate)
			val := rate.Value(tt.date, tax.Extensions{})
			if tt.percent == "" {
				assert.Nil(t, val, "%s must not resolve on %s", tt.rate, tt.date)
				return
			}
			require.NotNil(t, val)
			assert.Equal(t, tt.percent, val.Percent.String())
		})
	}
}
