package l10n_test

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/stretchr/testify/assert"
)

func TestCountries(t *testing.T) {
	t.Parallel()
	c := l10n.Countries().Code(l10n.US)
	assert.NotNil(t, c)
}

func TestCountriesMap(t *testing.T) {
	t.Parallel()
	t.Run("contains all countries in the expected format", func(t *testing.T) {
		t.Parallel()
		countriesMap := l10n.CountriesMap()
		countriesList := l10n.Countries()

		assert.Equal(t, len(countriesList), len(countriesMap),
			"Map and slice should contain the same number of countries")

		// Verify every country from the slice exists in the map
		for _, country := range countriesList {
			mapEntry, exists := countriesMap[country.Code]
			assert.True(t, exists, "Country %s should exist in map", country.Code)
			assert.Equal(t, country, mapEntry, "Country data should match between map and slice")
		}
	})

	t.Run("all countries have the required fields", func(t *testing.T) {
		t.Parallel()
		countriesMap := l10n.CountriesMap()

		for code, country := range countriesMap {
			assert.NotEmpty(t, country.Code, "Country %s should have a code", code)
			assert.NotEmpty(t, country.Alpha3, "Country %s should have Alpha3 code", code)
			assert.NotEmpty(t, country.Name, "Country %s should have a name", code)
			// TLD can be empty for some countries (like Svalbard and Jan Mayen, US Minor Outlying Islands)
			// ISO and Tax are booleans, so they always have values
		}
	})

	t.Run("returns a new map instance each time", func(t *testing.T) {
		t.Parallel()
		map1 := l10n.CountriesMap()
		map2 := l10n.CountriesMap()

		assert.Equal(t, len(map1), len(map2), "Maps should have same length")

		delete(map1, l10n.US)

		assert.NotEqual(t, len(map1), len(map2), "First map should be modified without affecting the second")
	})
}
