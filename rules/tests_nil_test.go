package rules_test

import (
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

func TestNil(t *testing.T) {
	t.Run("passes for nil pointer", func(t *testing.T) {
		var s *string
		assert.True(t, rules.Nil.Check(s))
	})

	t.Run("fails for empty string", func(t *testing.T) {
		assert.False(t, rules.Nil.Check(""))
	})

	t.Run("fails for non-empty string", func(t *testing.T) {
		assert.False(t, rules.Nil.Check("hello"))
	})

	t.Run("fails for zero int", func(t *testing.T) {
		assert.False(t, rules.Nil.Check(0))
	})

	t.Run("fails for false bool", func(t *testing.T) {
		assert.False(t, rules.Nil.Check(false))
	})

	t.Run("fails for non-nil pointer to empty value", func(t *testing.T) {
		s := ""
		assert.False(t, rules.Nil.Check(&s))
	})

	t.Run("fails for non-nil pointer to non-empty value", func(t *testing.T) {
		s := "hello"
		assert.False(t, rules.Nil.Check(&s))
	})

	t.Run("String returns nil", func(t *testing.T) {
		assert.Equal(t, "nil", rules.Nil.String())
	})
}
