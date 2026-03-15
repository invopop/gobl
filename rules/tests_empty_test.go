package rules_test

import (
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	t.Run("passes for empty string", func(t *testing.T) {
		assert.True(t, rules.Empty.Check(""))
	})

	t.Run("fails for non-empty string", func(t *testing.T) {
		assert.False(t, rules.Empty.Check("hello"))
	})

	t.Run("passes for zero int", func(t *testing.T) {
		assert.True(t, rules.Empty.Check(0))
	})

	t.Run("fails for non-zero int", func(t *testing.T) {
		assert.False(t, rules.Empty.Check(42))
	})

	t.Run("passes for false bool", func(t *testing.T) {
		assert.True(t, rules.Empty.Check(false))
	})

	t.Run("fails for true bool", func(t *testing.T) {
		assert.False(t, rules.Empty.Check(true))
	})

	t.Run("passes for empty slice", func(t *testing.T) {
		assert.True(t, rules.Empty.Check([]string{}))
	})

	t.Run("fails for non-empty slice", func(t *testing.T) {
		assert.False(t, rules.Empty.Check([]string{"a"}))
	})

	t.Run("passes for nil pointer", func(t *testing.T) {
		var s *string
		assert.True(t, rules.Empty.Check(s))
	})

	t.Run("passes for non-nil pointer to empty value", func(t *testing.T) {
		s := ""
		assert.True(t, rules.Empty.Check(&s))
	})

	t.Run("fails for non-nil pointer to non-empty value", func(t *testing.T) {
		s := "hello"
		assert.False(t, rules.Empty.Check(&s))
	})

	t.Run("String returns empty", func(t *testing.T) {
		assert.Equal(t, "empty", rules.Empty.String())
	})
}
