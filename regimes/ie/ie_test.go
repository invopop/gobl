package ie_test

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/ie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	regime := ie.New()
	require.NotNil(t, regime)
	assert.Equal(t, l10n.IE, regime.Country.Code())
	assert.Equal(t, "Ireland", regime.Name.String())
	assert.NotNil(t, regime.Categories)
	assert.NotNil(t, regime.Validator)
	assert.NotNil(t, regime.Normalizer)
	assert.Len(t, regime.Categories, 1)
	assert.Equal(t, "VAT", regime.Categories[0].Code.String())
}
