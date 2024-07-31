package l10n_test

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/stretchr/testify/assert"
)

func TestISOCountryCodeValidation(t *testing.T) {
	// Test that the country code is valid
	c := l10n.ISOCountryCode("US")
	err := c.Validate()
	assert.NoError(t, err)

	// Test that the country code is invalid
	c = l10n.ISOCountryCode("XX")
	err = c.Validate()
	assert.ErrorContains(t, err, "must be a valid ISO country code")

	c = l10n.ISOCountryCode("XI") // Northern Ireland
	err = c.Validate()
	assert.ErrorContains(t, err, "must be a valid ISO country code")
}

func TestTaxCountryCodeValidation(t *testing.T) {
	// Test that the country code is valid
	c := l10n.TaxCountryCode("US")
	err := c.Validate()
	assert.NoError(t, err)

	// Test that the country code is invalid
	c = l10n.TaxCountryCode("XX")
	err = c.Validate()
	assert.ErrorContains(t, err, "must be a valid tax country code")

	c = l10n.TaxCountryCode("XI") // Northern Ireland
	err = c.Validate()
	assert.NoError(t, err)
}

func TestISOCountryCodeIn(t *testing.T) {
	// Test that the country code is valid
	c := l10n.ISOCountryCode("US")
	assert.True(t, c.In("US", "CA"))

	// Test that the country code is invalid
	c = l10n.ISOCountryCode("XI")
	assert.False(t, c.In("US", "CA"))
}

func TestISOCountryCodeOutput(t *testing.T) {
	c := l10n.ISOCountryCode("US")

	assert.Equal(t, "US", c.String())
	assert.Equal(t, "United States of America", c.Name())
	assert.Equal(t, "USA", c.Alpha3())
	assert.Equal(t, l10n.US, c.Code())

	c = l10n.ISOCountryCode("XX")
	assert.Empty(t, c.Name())
	assert.Empty(t, c.Alpha3())

}

func TestISOCountryCodeSchema(t *testing.T) {
	s := l10n.ISOCountryCode("").JSONSchema()
	assert.Equal(t, "ISO Country Code", s.Title)
	assert.Equal(t, "string", s.Type)
	assert.Equal(t, 249, len(s.OneOf))
	assert.Equal(t, l10n.AF, s.OneOf[0].Const)
	assert.Equal(t, "Afghanistan", s.OneOf[0].Title)
}

func TestTaxCountryCodeIn(t *testing.T) {
	// Test that the country code is valid
	c := l10n.TaxCountryCode("US")
	assert.True(t, c.In("US", "CA"))

	// Test that the country code is invalid
	c = l10n.TaxCountryCode("XI")
	assert.False(t, c.In("US", "CA"))
}

func TestTaxCountryCodeOutput(t *testing.T) {
	c := l10n.TaxCountryCode("US")
	assert.Equal(t, "US", c.String())
	assert.Equal(t, "United States of America", c.Name())
	assert.Equal(t, l10n.US, c.Code())
	c = l10n.TaxCountryCode("XX")
	assert.Empty(t, c.Name())
}

func TestTaxCountryCodeSchema(t *testing.T) {
	s := l10n.TaxCountryCode("").JSONSchema()
	assert.Equal(t, "Tax Country Code", s.Title)
	assert.Equal(t, "string", s.Type)
	assert.Equal(t, 252, len(s.OneOf))
	assert.Equal(t, l10n.AF, s.OneOf[0].Const)
	assert.Equal(t, "Afghanistan", s.OneOf[0].Title)
}
