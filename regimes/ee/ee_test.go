package ee_test

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/ee"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	regime := ee.New()
	require.NotNil(t, regime)
	assert.Equal(t, l10n.EE, regime.Country.Code())
	assert.Equal(t, "Estonia", regime.Name.String())
	assert.NotNil(t, regime.Categories)
	assert.Len(t, regime.Categories, 1)
	assert.Equal(t, "VAT", regime.Categories[0].Code.String())
}
