package us_test

import (
	"testing"

	"github.com/invopop/gobl/regimes/us"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("should create a new US regime", func(t *testing.T) {
		regime := us.New()
		require.NotNil(t, regime)
		assert.Equal(t, "US", regime.Country.String())
		assert.Equal(t, "United States of America", regime.Name["en"])
	})
}
