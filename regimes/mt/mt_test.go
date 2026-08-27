package mt_test

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/mt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()
	regime := mt.New()
	require.NotNil(t, regime)
	assert.Equal(t, l10n.MT, regime.Country.Code())
	assert.Equal(t, "Malta", regime.Name.String())
	assert.Len(t, regime.Categories, 1)
	assert.Equal(t, "VAT", regime.Categories[0].Code.String())
	assert.NotEmpty(t, regime.Description)
	assert.Empty(t, regime.Tags)
	assert.Len(t, regime.Corrections, 1)
}
