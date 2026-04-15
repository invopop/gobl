package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextSetAndValue(t *testing.T) {
	t.Run("set and get", func(t *testing.T) {
		var ctx Context
		ctx.Set("foo", "bar")
		assert.Equal(t, "bar", ctx.Value("foo"))
	})

	t.Run("missing key returns nil", func(t *testing.T) {
		var ctx Context
		assert.Nil(t, ctx.Value("missing"))
	})

	t.Run("multiple keys coexist", func(t *testing.T) {
		var ctx Context
		ctx.Set("a", 1)
		ctx.Set("b", 2)
		assert.Equal(t, 1, ctx.Value("a"))
		assert.Equal(t, 2, ctx.Value("b"))
	})
}

func TestContextEach(t *testing.T) {
	t.Run("short-circuit on true", func(t *testing.T) {
		var ctx Context
		ctx.Set("a", 1)
		ctx.Set("b", 2)
		count := 0
		result := ctx.Each(func(v any) bool {
			count++
			return v.(int) == 1
		})
		assert.True(t, result)
		assert.Equal(t, 1, count)
	})

	t.Run("iterate all returning false", func(t *testing.T) {
		var ctx Context
		ctx.Set("a", 1)
		ctx.Set("b", 2)
		count := 0
		result := ctx.Each(func(_ any) bool {
			count++
			return false
		})
		assert.False(t, result)
		assert.Equal(t, 2, count)
	})

	t.Run("empty context", func(t *testing.T) {
		var ctx Context
		result := ctx.Each(func(_ any) bool {
			return true
		})
		assert.False(t, result)
	})
}
