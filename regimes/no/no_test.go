package no_test

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/no"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()
	regime := no.New()
	require.NotNil(t, regime)
	assert.Equal(t, l10n.NO, regime.Country.Code())
	assert.Equal(t, "Norway", regime.Name.String())
	assert.Len(t, regime.Categories, 1)
	assert.Equal(t, "VAT", regime.Categories[0].Code.String())
	assert.Len(t, regime.Identities, 1)
	assert.Len(t, regime.Tags, 1)
	assert.Len(t, regime.Corrections, 1)
	assert.NotNil(t, regime.Scenarios)
	assert.NotNil(t, regime.Validator)
	assert.NotNil(t, regime.Normalizer)
}
