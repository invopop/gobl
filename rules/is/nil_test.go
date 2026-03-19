package is_test

import (
	"testing"

	"github.com/invopop/gobl/rules/is"
	"github.com/stretchr/testify/assert"
)

func TestNil(t *testing.T) {
	t.Run("passes for nil pointer", func(t *testing.T) {
		var s *string
		assert.True(t, is.Nil.Check(s))
	})

	t.Run("fails for empty string", func(t *testing.T) {
		assert.False(t, is.Nil.Check(""))
	})

	t.Run("fails for non-empty string", func(t *testing.T) {
		assert.False(t, is.Nil.Check("hello"))
	})

	t.Run("fails for zero int", func(t *testing.T) {
		assert.False(t, is.Nil.Check(0))
	})

	t.Run("fails for false bool", func(t *testing.T) {
		assert.False(t, is.Nil.Check(false))
	})

	t.Run("fails for non-nil pointer to empty value", func(t *testing.T) {
		s := ""
		assert.False(t, is.Nil.Check(&s))
	})

	t.Run("fails for non-nil pointer to non-empty value", func(t *testing.T) {
		s := "hello"
		assert.False(t, is.Nil.Check(&s))
	})

	t.Run("String returns nil", func(t *testing.T) {
		assert.Equal(t, "nil", is.Nil.String())
	})
}
