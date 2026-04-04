package au_test

import (
	"testing"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/au"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	regime := au.New()

	assert.Equal(t, l10n.AU.Tax(), regime.Country)
	assert.Equal(t, currency.AUD, regime.Currency)
	assert.Equal(t, tax.CategoryGST, regime.TaxScheme)
	assert.Equal(t, "Australia", regime.Name.String())
	assert.NotEmpty(t, regime.Categories)
	assert.NotNil(t, regime.Validator)
	assert.NotNil(t, regime.Normalizer)
}

func TestRegimeValidation(t *testing.T) {
	t.Parallel()

	regime := au.New()
	require.NoError(t, regime.Validate())
}
