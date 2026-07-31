package is_test

import (
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNot(t *testing.T) {
	alwaysPass := is.Func("pass", func(any) bool { return true })
	alwaysFail := is.Func("fail", func(any) bool { return false })

	t.Run("wrapped test passes", func(t *testing.T) {
		assert.False(t, is.Not(alwaysPass).Check("x"))
	})

	t.Run("wrapped test fails", func(t *testing.T) {
		assert.True(t, is.Not(alwaysFail).Check("x"))
	})

	t.Run("String output", func(t *testing.T) {
		assert.Equal(t, "not (pass)", is.Not(alwaysPass).String())
	})

	t.Run("context-aware wrapped test", func(t *testing.T) {
		nt := is.Not(contextTest{key: "k", expect: "v"})
		ct, ok := nt.(rules.ContextualTest)
		require.True(t, ok)

		rc := &rules.Context{}
		rc.Set("k", "v")
		assert.False(t, ct.CheckWithContext(rc, nil))

		rc = &rules.Context{}
		rc.Set("k", "other")
		assert.True(t, ct.CheckWithContext(rc, nil))
	})

	t.Run("compiles wrapped test", func(t *testing.T) {
		type notThing struct {
			Code string `json:"code"`
		}
		set := rules.For(new(notThing),
			rules.Field("code",
				rules.Assert("001", "code must not be numeric",
					is.Not(is.Matches(`^\d+$`)),
				),
			),
		)
		assert.Nil(t, set.Validate(&notThing{Code: "foo"}))
		assert.NotNil(t, set.Validate(&notThing{Code: "1234"}))
	})

	t.Run("wrapped test needs no compilation", func(t *testing.T) {
		type notPlain struct {
			Code string `json:"code"`
		}
		set := rules.For(new(notPlain),
			rules.Field("code",
				rules.Assert("001", "code is required",
					is.Not(is.Empty),
				),
			),
		)
		assert.Nil(t, set.Validate(&notPlain{Code: "foo"}))
		assert.NotNil(t, set.Validate(&notPlain{}))
	})
}
