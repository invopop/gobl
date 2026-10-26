package cl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidRUT(t *testing.T) {
	t.Run("wrong type", func(t *testing.T) {
		got := isValidRUT("not a code")
		assert.False(t, got)
	})
}
