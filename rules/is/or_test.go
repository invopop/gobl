package is_test

import (
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/stretchr/testify/assert"
)

func TestOr(t *testing.T) {
	alwaysPass := is.Func("pass", func(any) bool { return true })
	alwaysFail := is.Func("fail", func(any) bool { return false })

	t.Run("first test passes", func(t *testing.T) {
		assert.True(t, is.Or(alwaysPass, alwaysFail).Check("x"))
	})

	t.Run("only test fails", func(t *testing.T) {
		assert.False(t, is.Or(alwaysFail).Check("x"))
	})

	t.Run("second test passes", func(t *testing.T) {
		assert.True(t, is.Or(alwaysFail, alwaysPass).Check("x"))
	})

	t.Run("no tests (empty Or)", func(t *testing.T) {
		assert.False(t, is.Or().Check("x"))
	})

	t.Run("String output", func(t *testing.T) {
		result := is.Or(alwaysPass, alwaysFail).String()
		assert.Equal(t, "pass, or fail", result)
	})
}

// contextTest is a rules.Test + rules.ContextualTest for testing Or's context forwarding.
type contextTest struct {
	key    rules.ContextKey
	expect any
}

func (ct contextTest) Check(any) bool { return false }

func (ct contextTest) CheckWithContext(rc *rules.Context, _ any) bool {
	return rc.Value(ct.key) == ct.expect
}

func (ct contextTest) String() string { return "ctx-test" }

func TestOrCheckWithContext(t *testing.T) {
	alwaysFail := is.Func("fail", func(any) bool { return false })

	t.Run("context-aware inner test passes", func(t *testing.T) {
		inner := contextTest{key: "k", expect: "v"}
		or := is.Or(inner)
		rc := &rules.Context{}
		rc.Set("k", "v")
		ct := or.(rules.ContextualTest)
		assert.True(t, ct.CheckWithContext(rc, nil))
	})

	t.Run("all fail", func(t *testing.T) {
		inner := contextTest{key: "k", expect: "v"}
		or := is.Or(inner)
		rc := &rules.Context{}
		rc.Set("k", "other")
		ct := or.(rules.ContextualTest)
		assert.False(t, ct.CheckWithContext(rc, nil))
	})

	t.Run("mix of ContextualTest and plain Test", func(t *testing.T) {
		inner := contextTest{key: "k", expect: "v"}
		or := is.Or(alwaysFail, inner)
		rc := &rules.Context{}
		rc.Set("k", "v")
		ct := or.(rules.ContextualTest)
		assert.True(t, ct.CheckWithContext(rc, nil))
	})

	t.Run("plain test passes in CheckWithContext", func(t *testing.T) {
		alwaysPass := is.Func("pass", func(any) bool { return true })
		or := is.Or(alwaysPass)
		rc := &rules.Context{}
		ct := or.(rules.ContextualTest)
		assert.True(t, ct.CheckWithContext(rc, nil))
	})
}
