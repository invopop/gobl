package au_test

import (
	"testing"

	"github.com/invopop/gobl/tax"
	_ "github.com/invopop/gobl/regimes/au"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegimeRegistration(t *testing.T) {
	regime := tax.RegimeDefFor("AU")
	require.NotNil(t, regime, "AU regime should be registered")
	assert.Equal(t, "AU", regime.Country.String())
	assert.Equal(t, "Australia/Sydney", regime.TimeZone)
	assert.Equal(t, "AUD", regime.Currency.String())
}
