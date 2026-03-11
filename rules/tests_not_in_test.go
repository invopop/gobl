package rules_test

import (
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

func TestNotIn(t *testing.T) {
	t.Run("passes when value is not in set", func(t *testing.T) {
		assert.True(t, rules.NotIn("A", "B", "C").Check("D"))
	})

	t.Run("fails when value is in set", func(t *testing.T) {
		assert.False(t, rules.NotIn("A", "B", "C").Check("B"))
	})

	t.Run("passes when nil pointer", func(t *testing.T) {
		var s *string
		assert.True(t, rules.NotIn("A", "B").Check(s))
	})

	t.Run("named string type matches string literal in set", func(t *testing.T) {
		assert.False(t, rules.NotIn("standard", "reduced").Check(taxCode("standard")))
		assert.True(t, rules.NotIn("standard", "reduced").Check(taxCode("exempt")))
	})

	t.Run("bool values", func(t *testing.T) {
		assert.False(t, rules.NotIn(true).Check(true))
		assert.True(t, rules.NotIn(true).Check(false))
	})

	t.Run("empty set always passes", func(t *testing.T) {
		assert.True(t, rules.NotIn().Check("anything"))
	})

	t.Run("integer values", func(t *testing.T) {
		assert.False(t, rules.NotIn(1, 2, 3).Check(2))
		assert.True(t, rules.NotIn(1, 2, 3).Check(4))
	})

	t.Run("String lists the excluded values", func(t *testing.T) {
		assert.Equal(t, "not one of [A, B, C]", rules.NotIn("A", "B", "C").String())
	})

	t.Run("named types in set compiled via For", func(t *testing.T) {
		set := rules.For(taxCode(""),
			rules.Assert("001", "must not be reserved",
				rules.NotIn(taxCode("reserved"), taxCode("internal")),
			),
		)
		assert.Nil(t, set.Validate(taxCode("standard")))
		assert.NotNil(t, set.Validate(taxCode("reserved")))
	})
}
