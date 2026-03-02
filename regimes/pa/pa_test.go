package pa_test

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/pa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	regime := pa.New()
	require.NotNil(t, regime)
	assert.Equal(t, l10n.PA, regime.Country.Code())
	assert.Equal(t, "Panama", regime.Name.String())
	assert.NotNil(t, regime.Categories)
	assert.NotNil(t, regime.Validator)
	assert.NotNil(t, regime.Normalizer)
	assert.Len(t, regime.Categories, 2)
	assert.Equal(t, "VAT", regime.Categories[0].Code.String())
	assert.Equal(t, "ISC", regime.Categories[1].Code.String())
}
