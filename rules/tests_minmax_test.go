package rules_test

import (
	"testing"
	"time"

	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

// amount is a named int type used to exercise Min/Max with named types.
type amount int

func TestMin(t *testing.T) {
	t.Run("passes when int value equals threshold", func(t *testing.T) {
		assert.True(t, rules.Min(5).Check(5))
	})

	t.Run("passes when int value exceeds threshold", func(t *testing.T) {
		assert.True(t, rules.Min(5).Check(10))
	})

	t.Run("fails when int value is below threshold", func(t *testing.T) {
		assert.False(t, rules.Min(5).Check(3))
	})

	t.Run("passes for zero (empty) value", func(t *testing.T) {
		assert.True(t, rules.Min(5).Check(0))
	})

	t.Run("passes for nil pointer", func(t *testing.T) {
		var p *int
		assert.True(t, rules.Min(5).Check(p))
	})

	t.Run("named int type works", func(t *testing.T) {
		assert.True(t, rules.Min(5).Check(amount(10)))
		assert.False(t, rules.Min(5).Check(amount(3)))
	})

	t.Run("uint threshold", func(t *testing.T) {
		assert.True(t, rules.Min(uint(5)).Check(uint(10)))
		assert.False(t, rules.Min(uint(5)).Check(uint(3)))
		assert.True(t, rules.Min(uint(5)).Check(uint(5)))
	})

	t.Run("float threshold", func(t *testing.T) {
		assert.True(t, rules.Min(1.5).Check(2.0))
		assert.False(t, rules.Min(1.5).Check(1.0))
		assert.True(t, rules.Min(1.5).Check(1.5))
	})

	t.Run("time.Time threshold", func(t *testing.T) {
		base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		assert.True(t, rules.Min(base).Check(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)))
		assert.True(t, rules.Min(base).Check(base))
		assert.False(t, rules.Min(base).Check(time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)))
	})

	t.Run("zero time is treated as empty", func(t *testing.T) {
		base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		assert.True(t, rules.Min(base).Check(time.Time{}))
	})

	t.Run("String returns at least description", func(t *testing.T) {
		assert.Equal(t, "at least 5", rules.Min(5).String())
	})
}

func TestMax(t *testing.T) {
	t.Run("passes when int value equals threshold", func(t *testing.T) {
		assert.True(t, rules.Max(10).Check(10))
	})

	t.Run("passes when int value is below threshold", func(t *testing.T) {
		assert.True(t, rules.Max(10).Check(5))
	})

	t.Run("fails when int value exceeds threshold", func(t *testing.T) {
		assert.False(t, rules.Max(10).Check(15))
	})

	t.Run("passes for zero (empty) value", func(t *testing.T) {
		assert.True(t, rules.Max(10).Check(0))
	})

	t.Run("passes for nil pointer", func(t *testing.T) {
		var p *int
		assert.True(t, rules.Max(10).Check(p))
	})

	t.Run("float threshold", func(t *testing.T) {
		assert.True(t, rules.Max(1.5).Check(1.0))
		assert.False(t, rules.Max(1.5).Check(2.0))
		assert.True(t, rules.Max(1.5).Check(1.5))
	})

	t.Run("time.Time threshold", func(t *testing.T) {
		base := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
		assert.True(t, rules.Max(base).Check(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)))
		assert.True(t, rules.Max(base).Check(base))
		assert.False(t, rules.Max(base).Check(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)))
	})

	t.Run("String returns at most description", func(t *testing.T) {
		assert.Equal(t, "at most 10", rules.Max(10).String())
	})
}

func TestMinMaxExclusive(t *testing.T) {
	t.Run("Min Exclusive passes when strictly greater", func(t *testing.T) {
		assert.True(t, rules.Min(5).Exclusive().Check(6))
	})

	t.Run("Min Exclusive fails when equal to threshold", func(t *testing.T) {
		assert.False(t, rules.Min(5).Exclusive().Check(5))
	})

	t.Run("Min Exclusive fails when below threshold", func(t *testing.T) {
		assert.False(t, rules.Min(5).Exclusive().Check(3))
	})

	t.Run("Max Exclusive passes when strictly less", func(t *testing.T) {
		assert.True(t, rules.Max(10).Exclusive().Check(9))
	})

	t.Run("Max Exclusive fails when equal to threshold", func(t *testing.T) {
		assert.False(t, rules.Max(10).Exclusive().Check(10))
	})

	t.Run("Max Exclusive fails when above threshold", func(t *testing.T) {
		assert.False(t, rules.Max(10).Exclusive().Check(15))
	})

	t.Run("Min Exclusive String returns greater than description", func(t *testing.T) {
		assert.Equal(t, "greater than 5", rules.Min(5).Exclusive().String())
	})

	t.Run("Max Exclusive String returns less than description", func(t *testing.T) {
		assert.Equal(t, "less than 10", rules.Max(10).Exclusive().String())
	})

	t.Run("time.Time exclusive", func(t *testing.T) {
		base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		assert.True(t, rules.Min(base).Exclusive().Check(time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)))
		assert.False(t, rules.Min(base).Exclusive().Check(base))
		assert.True(t, rules.Max(base).Exclusive().Check(time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)))
		assert.False(t, rules.Max(base).Exclusive().Check(base))
	})
}
