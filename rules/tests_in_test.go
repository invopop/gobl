package rules_test

import (
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

// taxCode is a named string type used to exercise In with cross-type comparison.
type taxCode string

func TestIn(t *testing.T) {
	t.Run("passes when value is in set", func(t *testing.T) {
		assert.True(t, rules.In("A", "B", "C").Check("B"))
	})

	t.Run("fails when value is not in set", func(t *testing.T) {
		assert.False(t, rules.In("A", "B", "C").Check("D"))
	})

	t.Run("passes when nil pointer", func(t *testing.T) {
		var s *string
		assert.True(t, rules.In("A", "B").Check(s))
	})

	t.Run("named string type matches string literal in set", func(t *testing.T) {
		assert.True(t, rules.In("standard", "reduced").Check(taxCode("standard")))
		assert.False(t, rules.In("standard", "reduced").Check(taxCode("exempt")))
	})

	t.Run("bool values", func(t *testing.T) {
		assert.True(t, rules.In(true).Check(true))
		assert.False(t, rules.In(true).Check(false))
	})

	t.Run("integer values", func(t *testing.T) {
		assert.True(t, rules.In(1, 2, 3).Check(2))
		assert.False(t, rules.In(1, 2, 3).Check(4))
	})

	t.Run("empty set always fails for non-nil value", func(t *testing.T) {
		assert.False(t, rules.In().Check("anything"))
	})

	t.Run("String lists the set values", func(t *testing.T) {
		assert.Equal(t, "one of [A, B, C]", rules.In("A", "B", "C").String())
	})

	t.Run("named types in set compiled via For", func(t *testing.T) {
		set := rules.For(taxCode(""),
			rules.Assert("001", "must be valid",
				rules.In(taxCode("standard"), taxCode("reduced")),
			),
		)
		assert.Nil(t, set.Validate(taxCode("standard")))
		assert.NotNil(t, set.Validate(taxCode("exempt")))
	})

	t.Run("integer set compiled via For", func(t *testing.T) {
		type level int
		set := rules.For(level(0),
			rules.Assert("001", "must be valid",
				rules.In(1, 2, 3),
			),
		)
		assert.Nil(t, set.Validate(level(1)))
		assert.NotNil(t, set.Validate(level(4)))
	})
}
