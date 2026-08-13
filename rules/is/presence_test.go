package is_test

import (
	"testing"

	"github.com/invopop/gobl/rules/is"
	"github.com/stretchr/testify/assert"
)

func TestPresent(t *testing.T) {
	t.Run("passes for non-empty string", func(t *testing.T) {
		assert.True(t, is.Present.Check("hello"))
	})

	t.Run("fails for empty string", func(t *testing.T) {
		assert.False(t, is.Present.Check(""))
	})

	t.Run("passes for non-zero int", func(t *testing.T) {
		assert.True(t, is.Present.Check(42))
	})

	t.Run("fails for zero int", func(t *testing.T) {
		assert.False(t, is.Present.Check(0))
	})

	t.Run("passes for true bool", func(t *testing.T) {
		assert.True(t, is.Present.Check(true))
	})

	t.Run("fails for false bool", func(t *testing.T) {
		assert.False(t, is.Present.Check(false))
	})

	t.Run("passes for non-empty slice", func(t *testing.T) {
		assert.True(t, is.Present.Check([]string{"a"}))
	})

	t.Run("fails for empty slice", func(t *testing.T) {
		assert.False(t, is.Present.Check([]string{}))
	})

	t.Run("fails for nil pointer", func(t *testing.T) {
		var s *string
		assert.False(t, is.Present.Check(s))
	})

	t.Run("passes for non-nil pointer to non-empty value", func(t *testing.T) {
		s := "hello"
		assert.True(t, is.Present.Check(&s))
	})

	t.Run("fails for non-nil pointer to empty value", func(t *testing.T) {
		s := ""
		assert.False(t, is.Present.Check(&s))
	})

	t.Run("String returns present", func(t *testing.T) {
		assert.Equal(t, "present", is.Present.String())
	})
}

func TestNilOrNotEmpty(t *testing.T) {
	t.Run("passes for nil pointer", func(t *testing.T) {
		var s *string
		assert.True(t, is.NilOrNotEmpty.Check(s))
	})

	t.Run("passes for non-nil pointer to non-empty value", func(t *testing.T) {
		s := "hello"
		assert.True(t, is.NilOrNotEmpty.Check(&s))
	})

	t.Run("fails for non-nil pointer to empty value", func(t *testing.T) {
		s := ""
		assert.False(t, is.NilOrNotEmpty.Check(&s))
	})

	t.Run("passes for non-empty string", func(t *testing.T) {
		assert.True(t, is.NilOrNotEmpty.Check("hello"))
	})

	t.Run("fails for empty string", func(t *testing.T) {
		assert.False(t, is.NilOrNotEmpty.Check(""))
	})

	t.Run("String returns nil or not empty", func(t *testing.T) {
		assert.Equal(t, "nil or not empty", is.NilOrNotEmpty.String())
	})
}
