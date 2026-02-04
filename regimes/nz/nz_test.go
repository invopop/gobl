package nz_test

import (
	"testing"

	"github.com/invopop/gobl/regimes/nz"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRegime(t *testing.T) {
	r := nz.New()
	require.NotNil(t, r)
	assert.Equal(t, "NZ", r.Country.String())
	assert.Equal(t, "NZD", r.Currency.String())
	assert.Equal(t, "Pacific/Auckland", r.TimeZone)
	assert.Equal(t, tax.CategoryGST, r.TaxScheme)
	assert.NotEmpty(t, r.Scenarios, "Should have scenarios")
	assert.NotEmpty(t, r.Corrections, "Should have corrections")
}
