package nz_test

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/nz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	regime := nz.New()
	require.NotNil(t, regime)
	assert.Equal(t, l10n.NZ, regime.Country.Code())
	assert.Equal(t, "New Zealand", regime.Name.String())
	assert.Len(t, regime.Categories, 1)
	assert.Equal(t, "GST", regime.Categories[0].Code.String())
	assert.NotEmpty(t, regime.Description)
	assert.Len(t, regime.Corrections, 1)
}
