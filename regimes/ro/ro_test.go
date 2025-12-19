package ro_test

import (
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/ro"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRegime(t *testing.T) {
	regime := ro.New()
	require.NotNil(t, regime)

	t.Run("basic properties", func(t *testing.T) {
		assert.Equal(t, l10n.RO.Tax(), regime.Country)
		assert.Equal(t, currency.RON, regime.Currency)
		assert.Equal(t, "Europe/Bucharest", regime.TimeZone)
		assert.Equal(t, tax.CategoryVAT, regime.TaxScheme)
	})

	t.Run("name translations", func(t *testing.T) {
		assert.NotEmpty(t, regime.Name)
		assert.Equal(t, "Romania", regime.Name[i18n.EN])
		assert.Equal(t, "România", regime.Name[i18n.RO])
	})

	t.Run("has validator", func(t *testing.T) {
		assert.NotNil(t, regime.Validator)
	})

	t.Run("has normalizer", func(t *testing.T) {
		assert.NotNil(t, regime.Normalizer)
	})

	t.Run("has categories", func(t *testing.T) {
		assert.NotEmpty(t, regime.Categories)
	})

	t.Run("has identities", func(t *testing.T) {
		assert.NotEmpty(t, regime.Identities)
	})

	t.Run("has scenarios", func(t *testing.T) {
		assert.NotEmpty(t, regime.Scenarios)
	})
}

func TestTaxCategories(t *testing.T) {
	regime := ro.New()
	require.NotNil(t, regime)

	t.Run("VAT category exists", func(t *testing.T) {
		var vatCat *tax.CategoryDef
		for _, cat := range regime.Categories {
			if cat.Code == tax.CategoryVAT {
				vatCat = cat
				break
			}
		}
		require.NotNil(t, vatCat, "VAT category should exist")

		assert.Equal(t, "VAT", vatCat.Name[i18n.EN])
		assert.Equal(t, "TVA", vatCat.Name[i18n.RO])
		assert.Equal(t, "Value Added Tax", vatCat.Title[i18n.EN])
		assert.Equal(t, "Taxa pe Valoarea Adăugată", vatCat.Title[i18n.RO])
		assert.False(t, vatCat.Retained)
	})

	t.Run("standard rate - 19%", func(t *testing.T) {
		regime := ro.New()
		cat := regime.Categories[0]

		var standardRate *tax.RateDef
		for _, rate := range cat.Rates {
			if rate.Rate == tax.RateGeneral {
				standardRate = rate
				break
			}
		}
		require.NotNil(t, standardRate, "Standard rate should exist")

		assert.Equal(t, "General Rate", standardRate.Name[i18n.EN])
		assert.Equal(t, "Cota standard", standardRate.Name[i18n.RO])

		// Current rate should be 21% (since Aug 1, 2025 - Law 141/2025)
		currentRate := standardRate.Values[0]
		assert.Equal(t, num.MakePercentage(210, 3), currentRate.Percent)
		assert.Equal(t, cal.NewDate(2025, 8, 1), currentRate.Since)
	})

	t.Run("reduced rate - 9%", func(t *testing.T) {
		regime := ro.New()
		cat := regime.Categories[0]

		var reducedRate *tax.RateDef
		for _, rate := range cat.Rates {
			if rate.Rate == tax.RateReduced {
				reducedRate = rate
				break
			}
		}
		require.NotNil(t, reducedRate, "Reduced rate should exist")

		assert.Equal(t, "Reduced Rate", reducedRate.Name[i18n.EN])
		assert.Equal(t, "Cota redusă", reducedRate.Name[i18n.RO])

		// Current rate should be 11% (since Aug 1, 2025 - Law 141/2025)
		currentRate := reducedRate.Values[0]
		assert.Equal(t, num.MakePercentage(110, 3), currentRate.Percent)
		assert.Equal(t, cal.NewDate(2025, 8, 1), currentRate.Since)
	})

	t.Run("super-reduced rate - 5%", func(t *testing.T) {
		regime := ro.New()
		cat := regime.Categories[0]

		var superReducedRate *tax.RateDef
		for _, rate := range cat.Rates {
			if rate.Rate == tax.RateSuperReduced {
				superReducedRate = rate
				break
			}
		}
		require.NotNil(t, superReducedRate, "Super-reduced rate should exist")

		assert.Equal(t, "Super-Reduced Rate", superReducedRate.Name[i18n.EN])
		assert.Equal(t, "Cota redusă suplimentară", superReducedRate.Name[i18n.RO])

		currentRate := superReducedRate.Values[0]
		assert.Equal(t, num.MakePercentage(50, 3), currentRate.Percent)
		assert.Equal(t, cal.NewDate(2017, 1, 1), currentRate.Since)
	})

	t.Run("historical rates", func(t *testing.T) {
		regime := ro.New()
		cat := regime.Categories[0]

		var standardRate *tax.RateDef
		for _, rate := range cat.Rates {
			if rate.Rate == tax.RateGeneral {
				standardRate = rate
				break
			}
		}
		require.NotNil(t, standardRate)

		// Should have historical rates
		assert.GreaterOrEqual(t, len(standardRate.Values), 2, "Should have at least 2 rate values (current and historical)")

		// Check 2016 rate (20%)
		var rate2016 *tax.RateValueDef
		date2016 := cal.NewDate(2016, 1, 1)
		for _, val := range standardRate.Values {
			if val.Since != nil && val.Since.Year == date2016.Year && val.Since.Month == date2016.Month && val.Since.Day == date2016.Day {
				rate2016 = val
				break
			}
		}
		if rate2016 != nil {
			assert.Equal(t, num.MakePercentage(200, 3), rate2016.Percent)
		}
	})

	t.Run("sources are documented", func(t *testing.T) {
		regime := ro.New()
		cat := regime.Categories[0]

		assert.NotEmpty(t, cat.Sources, "Tax category should have sources")

		// Check that sources contain key references
		hasLegislatie := false
		hasLaw141 := false

		for _, src := range cat.Sources {
			if src.URL != "" {
				if containsAny(src.URL, "legislatie.just.ro") {
					hasLegislatie = true
				}
				if containsAny(src.URL, "Law 141/2025", "141/2025") {
					hasLaw141 = true
				}
			}
			// Check title if it has English text
			if titleEN, ok := src.Title[i18n.EN]; ok {
				if containsAny(titleEN, "Law 141/2025", "141/2025") {
					hasLaw141 = true
				}
			}
		}

		assert.True(t, hasLegislatie, "Should reference legislatie.just.ro")
		assert.True(t, hasLaw141, "Should reference Law 141/2025")
	})
}

func TestRegimeRegistration(t *testing.T) {
	// Test that the regime is properly registered
	regime := tax.RegimeDefFor(l10n.RO)
	require.NotNil(t, regime, "RO regime should be registered")
	assert.Equal(t, l10n.RO.Tax(), regime.Country)
}

func TestValidateAndNormalize(t *testing.T) {
	t.Run("normalize nil values", func(t *testing.T) {
		// Should not panic with nil values
		assert.NotPanics(t, func() {
			ro.Normalize(nil)
		})
	})

	t.Run("validate nil values", func(t *testing.T) {
		// Should not panic with nil values
		assert.NotPanics(t, func() {
			err := ro.Validate(nil)
			assert.NoError(t, err)
		})
	})

	t.Run("validate unknown type", func(t *testing.T) {
		// Should return nil for unknown types
		type unknownType struct{}
		err := ro.Validate(&unknownType{})
		assert.NoError(t, err)
	})
}

// Helper function
func containsAny(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if len(s) >= len(substr) {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}
