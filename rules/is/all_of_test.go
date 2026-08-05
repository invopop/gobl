package is_test

import (
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/stretchr/testify/assert"
)

func TestAllOf(t *testing.T) {
	alwaysPass := is.Func("pass", func(any) bool { return true })
	alwaysFail := is.Func("fail", func(any) bool { return false })

	t.Run("all tests pass", func(t *testing.T) {
		assert.True(t, is.AllOf(alwaysPass, alwaysPass).Check("x"))
	})

	t.Run("only test passes", func(t *testing.T) {
		assert.True(t, is.AllOf(alwaysPass).Check("x"))
	})

	t.Run("first test fails", func(t *testing.T) {
		assert.False(t, is.AllOf(alwaysFail, alwaysPass).Check("x"))
	})

	t.Run("second test fails", func(t *testing.T) {
		assert.False(t, is.AllOf(alwaysPass, alwaysFail).Check("x"))
	})

	t.Run("no tests (empty AllOf)", func(t *testing.T) {
		assert.True(t, is.AllOf().Check("x"))
	})

	t.Run("String output", func(t *testing.T) {
		result := is.AllOf(alwaysPass, alwaysFail).String()
		assert.Equal(t, "pass, and fail", result)
	})
}

func TestAllOfCheckWithContext(t *testing.T) {
	alwaysPass := is.Func("pass", func(any) bool { return true })

	t.Run("context-aware inner test passes", func(t *testing.T) {
		inner := contextTest{key: "k", expect: "v"}
		aoTest := is.AllOf(inner, alwaysPass)
		rc := &rules.Context{}
		rc.Set("k", "v")
		ct := aoTest.(rules.ContextualTest)
		assert.True(t, ct.CheckWithContext(rc, nil))
	})

	t.Run("context-aware inner test fails", func(t *testing.T) {
		inner := contextTest{key: "k", expect: "v"}
		aoTest := is.AllOf(inner, alwaysPass)
		rc := &rules.Context{}
		rc.Set("k", "other")
		ct := aoTest.(rules.ContextualTest)
		assert.False(t, ct.CheckWithContext(rc, nil))
	})

	t.Run("plain inner test fails", func(t *testing.T) {
		alwaysFail := is.Func("fail", func(any) bool { return false })
		aoTest := is.AllOf(contextTest{key: "k", expect: "v"}, alwaysFail)
		rc := &rules.Context{}
		rc.Set("k", "v")
		ct := aoTest.(rules.ContextualTest)
		assert.False(t, ct.CheckWithContext(rc, nil))
	})
}

func TestAllOfCompile(t *testing.T) {
	type allOfThing struct {
		Code string `json:"code"`
	}

	t.Run("compiles wrapped tests", func(t *testing.T) {
		set := rules.For(new(allOfThing),
			rules.Field("code",
				rules.Assert("001", "code must be four digits",
					is.AllOf(is.Matches(`^\d+$`), is.Length(4, 4)),
				),
			),
		)
		assert.Nil(t, set.Validate(&allOfThing{Code: "1234"}))
		assert.NotNil(t, set.Validate(&allOfThing{Code: "123"}))
		assert.NotNil(t, set.Validate(&allOfThing{Code: "abcd"}))
	})

	t.Run("reports compilation errors", func(t *testing.T) {
		aoTest := is.AllOf(is.Matches(`^(\d+$`)).(interface{ Compile(any) error })
		assert.Error(t, aoTest.Compile(new(allOfThing)))
	})
}
