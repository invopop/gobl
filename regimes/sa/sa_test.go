package sa_test

import (
	"testing"

	"github.com/invopop/gobl/regimes/sa"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	r := sa.New()
	require.NotNil(t, r)

	assert.Equal(t, "SA", r.Country.String())
	assert.Equal(t, "SAR", string(r.Currency))
	assert.Equal(t, "Asia/Riyadh", r.TimeZone)
	assert.Equal(t, tax.CategoryVAT, r.TaxScheme)

	assert.NotEmpty(t, r.Categories)
	assert.NotEmpty(t, r.Identities)
	assert.NotEmpty(t, r.Scenarios)
	assert.NotEmpty(t, r.Corrections)
}
