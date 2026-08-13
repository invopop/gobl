package is_test

import (
	"testing"

	"github.com/invopop/gobl/rules/is"
	"github.com/stretchr/testify/assert"
)

func TestLength(t *testing.T) {
	t.Run("passes when string is within range", func(t *testing.T) {
		assert.True(t, is.Length(2, 5).Check("abc"))
	})

	t.Run("passes when string length equals min", func(t *testing.T) {
		assert.True(t, is.Length(3, 5).Check("abc"))
	})

	t.Run("passes when string length equals max", func(t *testing.T) {
		assert.True(t, is.Length(2, 5).Check("abcde"))
	})

	t.Run("fails when string is too short", func(t *testing.T) {
		assert.False(t, is.Length(5, 10).Check("abc"))
	})

	t.Run("fails when string is too long", func(t *testing.T) {
		assert.False(t, is.Length(1, 3).Check("abcde"))
	})

	t.Run("passes for empty string (skipped)", func(t *testing.T) {
		assert.True(t, is.Length(2, 5).Check(""))
	})

	t.Run("passes for nil pointer (skipped)", func(t *testing.T) {
		var s *string
		assert.True(t, is.Length(2, 5).Check(s))
	})

	t.Run("zero min means no lower bound", func(t *testing.T) {
		assert.True(t, is.Length(0, 5).Check("a"))
		assert.False(t, is.Length(0, 5).Check("abcdef"))
	})

	t.Run("zero max means no upper bound", func(t *testing.T) {
		assert.True(t, is.Length(2, 0).Check("abcdefghij"))
		assert.False(t, is.Length(2, 0).Check("a"))
	})

	t.Run("both zero requires empty value", func(t *testing.T) {
		assert.True(t, is.Length(0, 0).Check(""))
		assert.False(t, is.Length(0, 0).Check("a"))
	})

	t.Run("works on slices", func(t *testing.T) {
		assert.True(t, is.Length(2, 4).Check([]string{"a", "b", "c"}))
		assert.False(t, is.Length(2, 4).Check([]string{"a"}))
	})

	t.Run("works on maps", func(t *testing.T) {
		assert.True(t, is.Length(1, 3).Check(map[string]int{"a": 1, "b": 2}))
		assert.False(t, is.Length(3, 5).Check(map[string]int{"a": 1}))
	})

	t.Run("String returns length description", func(t *testing.T) {
		assert.Equal(t, "length between 2 and 5", is.Length(2, 5).String())
	})
}

func TestRuneLength(t *testing.T) {
	t.Run("passes for ASCII string within range", func(t *testing.T) {
		assert.True(t, is.RuneLength(2, 5).Check("abc"))
	})

	t.Run("passes for unicode string within rune range", func(t *testing.T) {
		// "héllo" is 5 runes but 6 bytes
		assert.True(t, is.RuneLength(2, 5).Check("héllo"))
	})

	t.Run("fails for unicode string exceeding rune limit", func(t *testing.T) {
		// "héllo!" is 6 runes
		assert.False(t, is.RuneLength(2, 5).Check("héllo!"))
	})

	t.Run("counts runes not bytes", func(t *testing.T) {
		// "日本語" is 3 runes but 9 bytes
		assert.True(t, is.RuneLength(1, 3).Check("日本語"))
		assert.False(t, is.RuneLength(1, 2).Check("日本語"))
	})

	t.Run("passes for empty string (skipped)", func(t *testing.T) {
		assert.True(t, is.RuneLength(2, 5).Check(""))
	})

	t.Run("passes for nil pointer (skipped)", func(t *testing.T) {
		var s *string
		assert.True(t, is.RuneLength(2, 5).Check(s))
	})

	t.Run("falls back to element count for non-string types", func(t *testing.T) {
		assert.True(t, is.RuneLength(2, 4).Check([]string{"a", "b", "c"}))
		assert.False(t, is.RuneLength(2, 4).Check([]string{"a"}))
	})

	t.Run("String returns rune length description", func(t *testing.T) {
		assert.Equal(t, "rune length between 2 and 5", is.RuneLength(2, 5).String())
	})
}
