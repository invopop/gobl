package rules_test

import (
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

func TestRequired(t *testing.T) {
	t.Run("passes for non-empty string", func(t *testing.T) {
		assert.True(t, rules.Required.Check("hello"))
	})

	t.Run("fails for empty string", func(t *testing.T) {
		assert.False(t, rules.Required.Check(""))
	})

	t.Run("passes for non-zero int", func(t *testing.T) {
		assert.True(t, rules.Required.Check(42))
	})

	t.Run("fails for zero int", func(t *testing.T) {
		assert.False(t, rules.Required.Check(0))
	})

	t.Run("passes for true bool", func(t *testing.T) {
		assert.True(t, rules.Required.Check(true))
	})

	t.Run("fails for false bool", func(t *testing.T) {
		assert.False(t, rules.Required.Check(false))
	})

	t.Run("passes for non-empty slice", func(t *testing.T) {
		assert.True(t, rules.Required.Check([]string{"a"}))
	})

	t.Run("fails for empty slice", func(t *testing.T) {
		assert.False(t, rules.Required.Check([]string{}))
	})

	t.Run("fails for nil pointer", func(t *testing.T) {
		var s *string
		assert.False(t, rules.Required.Check(s))
	})

	t.Run("passes for non-nil pointer to non-empty value", func(t *testing.T) {
		s := "hello"
		assert.True(t, rules.Required.Check(&s))
	})

	t.Run("fails for non-nil pointer to empty value", func(t *testing.T) {
		s := ""
		assert.False(t, rules.Required.Check(&s))
	})

	t.Run("String returns required", func(t *testing.T) {
		assert.Equal(t, "required", rules.Required.String())
	})
}

func TestNilOrNotEmpty(t *testing.T) {
	t.Run("passes for nil pointer", func(t *testing.T) {
		var s *string
		assert.True(t, rules.NilOrNotEmpty.Check(s))
	})

	t.Run("passes for non-nil pointer to non-empty value", func(t *testing.T) {
		s := "hello"
		assert.True(t, rules.NilOrNotEmpty.Check(&s))
	})

	t.Run("fails for non-nil pointer to empty value", func(t *testing.T) {
		s := ""
		assert.False(t, rules.NilOrNotEmpty.Check(&s))
	})

	t.Run("passes for non-empty string", func(t *testing.T) {
		assert.True(t, rules.NilOrNotEmpty.Check("hello"))
	})

	t.Run("fails for empty string", func(t *testing.T) {
		assert.False(t, rules.NilOrNotEmpty.Check(""))
	})

	t.Run("String returns nil or not empty", func(t *testing.T) {
		assert.Equal(t, "nil or not empty", rules.NilOrNotEmpty.String())
	})
}
