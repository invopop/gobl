package cbc_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValueDefintionValidation(t *testing.T) {
	vd := new(cbc.ValueDefinition)

	err := vd.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "value: cannot be blank")
	assert.Contains(t, err.Error(), "name: cannot be blank")
}
