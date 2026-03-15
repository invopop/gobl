package rules_test

import (
	"strings"
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	isUpperCase := func(s string) bool {
		return s == strings.ToUpper(s)
	}

	t.Run("passes when function returns true", func(t *testing.T) {
		assert.True(t, rules.String(isUpperCase).Check("HELLO"))
	})

	t.Run("fails when function returns false", func(t *testing.T) {
		assert.False(t, rules.String(isUpperCase).Check("hello"))
	})

	t.Run("fails for non-string value", func(t *testing.T) {
		assert.False(t, rules.String(isUpperCase).Check(42))
	})

	t.Run("String returns default description", func(t *testing.T) {
		assert.Equal(t, "custom string test", rules.String(isUpperCase).String())
	})
}

func TestByString(t *testing.T) {
	noSpaces := func(s string) bool {
		return !strings.Contains(s, " ")
	}

	t.Run("passes when function returns true", func(t *testing.T) {
		assert.True(t, rules.ByString("no spaces", noSpaces).Check("hello"))
	})

	t.Run("fails when function returns false", func(t *testing.T) {
		assert.False(t, rules.ByString("no spaces", noSpaces).Check("hello world"))
	})

	t.Run("fails for non-string value", func(t *testing.T) {
		assert.False(t, rules.ByString("no spaces", noSpaces).Check(42))
	})

	t.Run("String returns custom description", func(t *testing.T) {
		assert.Equal(t, "no spaces", rules.ByString("no spaces", noSpaces).String())
	})
}
