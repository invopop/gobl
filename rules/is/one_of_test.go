package is_test

import (
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/stretchr/testify/assert"
)

func TestOneOf(t *testing.T) {
	alwaysPass := is.Func("pass", func(any) bool { return true })
	alwaysFail := is.Func("fail", func(any) bool { return false })

	t.Run("exactly one passes", func(t *testing.T) {
		assert.True(t, is.OneOf(alwaysPass, alwaysFail).Check("x"))
	})

	t.Run("only test fails", func(t *testing.T) {
		assert.False(t, is.OneOf(alwaysFail).Check("x"))
	})

	t.Run("only test passes", func(t *testing.T) {
		assert.True(t, is.OneOf(alwaysPass).Check("x"))
	})

	t.Run("two pass", func(t *testing.T) {
		assert.False(t, is.OneOf(alwaysPass, alwaysPass).Check("x"))
	})

	t.Run("three pass", func(t *testing.T) {
		assert.False(t, is.OneOf(alwaysPass, alwaysPass, alwaysPass).Check("x"))
	})

	t.Run("one of many passes", func(t *testing.T) {
		assert.True(t, is.OneOf(alwaysFail, alwaysPass, alwaysFail).Check("x"))
	})

	t.Run("all fail", func(t *testing.T) {
		assert.False(t, is.OneOf(alwaysFail, alwaysFail, alwaysFail).Check("x"))
	})

	t.Run("no tests (empty OneOf)", func(t *testing.T) {
		assert.False(t, is.OneOf().Check("x"))
	})

	t.Run("String output", func(t *testing.T) {
		result := is.OneOf(alwaysPass, alwaysFail).String()
		assert.Equal(t, "exactly one of: pass, fail", result)
	})
}

func TestOneOfCheckWithContext(t *testing.T) {
	alwaysPass := is.Func("pass", func(any) bool { return true })
	alwaysFail := is.Func("fail", func(any) bool { return false })

	t.Run("context-aware inner test passes alone", func(t *testing.T) {
		inner := contextTest{key: "k", expect: "v"}
		one := is.OneOf(inner)
		rc := &rules.Context{}
		rc.Set("k", "v")
		ct := one.(rules.ContextualTest)
		assert.True(t, ct.CheckWithContext(rc, nil))
	})

	t.Run("all fail", func(t *testing.T) {
		inner := contextTest{key: "k", expect: "v"}
		one := is.OneOf(inner)
		rc := &rules.Context{}
		rc.Set("k", "other")
		ct := one.(rules.ContextualTest)
		assert.False(t, ct.CheckWithContext(rc, nil))
	})

	t.Run("plain test and contextual test both pass fails", func(t *testing.T) {
		inner := contextTest{key: "k", expect: "v"}
		one := is.OneOf(alwaysPass, inner)
		rc := &rules.Context{}
		rc.Set("k", "v")
		ct := one.(rules.ContextualTest)
		assert.False(t, ct.CheckWithContext(rc, nil))
	})

	t.Run("only contextual test passes among mixed", func(t *testing.T) {
		inner := contextTest{key: "k", expect: "v"}
		one := is.OneOf(alwaysFail, inner)
		rc := &rules.Context{}
		rc.Set("k", "v")
		ct := one.(rules.ContextualTest)
		assert.True(t, ct.CheckWithContext(rc, nil))
	})

	t.Run("plain test passes alone in CheckWithContext", func(t *testing.T) {
		one := is.OneOf(alwaysPass)
		rc := &rules.Context{}
		ct := one.(rules.ContextualTest)
		assert.True(t, ct.CheckWithContext(rc, nil))
	})
}
