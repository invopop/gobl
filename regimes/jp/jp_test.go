package jp_test

import (
	"testing"

	"github.com/invopop/gobl/regimes/jp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegimeValidation(t *testing.T) {
	regime := jp.New()
	require.NotNil(t, regime)

	err := regime.Validate()
	assert.NoError(t, err, "RegimeDef should be valid")
}
