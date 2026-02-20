package ro_test

import (
	"testing"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/ro"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	regime := ro.New()
	require.NotNil(t, regime)
	assert.Equal(t, l10n.RO, regime.Country.Code())
	assert.Equal(t, currency.RON, regime.Currency)
	assert.Equal(t, "Romania", regime.Name.String())
	require.Len(t, regime.Categories, 1)
	assert.Equal(t, tax.CategoryVAT, regime.Categories[0].Code)
	assert.NotNil(t, regime.Validator)
	assert.NotNil(t, regime.Normalizer)
}
