package rules_test

import (
	"errors"
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

func TestBy(t *testing.T) {
	isPositive := func(v any) bool {
		n, ok := v.(int)
		return ok && n > 0
	}

	t.Run("passes when function returns true", func(t *testing.T) {
		assert.True(t, rules.By("positive", isPositive).Check(5))
	})

	t.Run("fails when function returns false", func(t *testing.T) {
		assert.False(t, rules.By("positive", isPositive).Check(-1))
	})

	t.Run("fails for unexpected type", func(t *testing.T) {
		assert.False(t, rules.By("positive", isPositive).Check("hello"))
	})

	t.Run("String returns description", func(t *testing.T) {
		assert.Equal(t, "positive", rules.By("positive", isPositive).String())
	})
}

func TestByError(t *testing.T) {
	validate := func(v any) error {
		s, ok := v.(string)
		if !ok || s == "" {
			return errors.New("must be a non-empty string")
		}
		return nil
	}

	t.Run("passes when function returns nil", func(t *testing.T) {
		assert.True(t, rules.ByError("non-empty string", validate).Check("hello"))
	})

	t.Run("fails when function returns error", func(t *testing.T) {
		assert.False(t, rules.ByError("non-empty string", validate).Check(""))
	})

	t.Run("fails for wrong type", func(t *testing.T) {
		assert.False(t, rules.ByError("non-empty string", validate).Check(42))
	})

	t.Run("String returns description", func(t *testing.T) {
		assert.Equal(t, "non-empty string", rules.ByError("non-empty string", validate).String())
	})
}
