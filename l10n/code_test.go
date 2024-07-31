package l10n_test

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/stretchr/testify/assert"
)

func TestCodeValidate(t *testing.T) {
	// Test that the code is valid
	c := l10n.Code("US")
	err := c.Validate()
	assert.NoError(t, err)

	// Test that the code is invalid
	c = l10n.Code("X-X")
	err = c.Validate()
	assert.ErrorContains(t, err, "must be in a valid format")
}

func TestCodeIn(t *testing.T) {
	c := l10n.Code("MAD")

	assert.True(t, c.In("A", "MAD"))
	assert.False(t, c.In("A", "V"))
}

func TestCodeOutput(t *testing.T) {
	c := l10n.US

	assert.Equal(t, "US", c.String())
	assert.Equal(t, l10n.ISOCountryCode("US"), c.ISO())
	assert.Equal(t, l10n.TaxCountryCode("US"), c.Tax())
}
