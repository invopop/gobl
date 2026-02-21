package hr_test

import (
	"testing"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/hr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	regime := hr.New()
	require.NotNil(t, regime)
	assert.Equal(t, l10n.HR.Tax(), regime.Country)
	assert.Equal(t, currency.EUR, regime.Currency)
	assert.Equal(t, tax.CategoryVAT, regime.TaxScheme)
	assert.Equal(t, "Croatia", regime.Name.String())
	assert.NotNil(t, regime.Categories)
	assert.NotNil(t, regime.Identities)
	assert.NotNil(t, regime.Validator)
	assert.NotNil(t, regime.Normalizer)
	assert.Len(t, regime.Categories, 1)
	assert.Equal(t, "VAT", regime.Categories[0].Code.String())
	assert.Len(t, regime.Identities, 1)
	assert.Equal(t, "OIB", regime.Identities[0].Code.String())
}

func TestRegimeValidation(t *testing.T) {
	regime := hr.New()
	require.NotNil(t, regime)
	err := regime.Validate()
	assert.NoError(t, err, "regime definition should be valid")
}
