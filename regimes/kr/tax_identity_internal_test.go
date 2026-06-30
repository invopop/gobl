package kr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidBRNNonCode(t *testing.T) {
	t.Parallel()
	// isValidBRN guards against a non-cbc.Code value. The public validation
	// path always passes a cbc.Code, so exercise the guard directly.
	assert.False(t, isValidBRN(12345))
	assert.False(t, isValidBRN("1208147521"))
}
