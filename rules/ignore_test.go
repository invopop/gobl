package rules_test

import (
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ignoreItem exercises rules.Ignore via the single-set Set.Validate entry
// point. Builders below are package-level so For omits the package segment
// from the ID (same convention as emailRules), giving a predictable code.
type ignoreItem struct {
	Name string `json:"name"`
}

// ignoreNameRequired fails assertion "01" when Name is empty. Extra defs
// (e.g. rules.Ignore or rules.When) are appended to the set.
func ignoreNameRequired(extra ...rules.Def) *rules.Set {
	defs := append([]rules.Def{
		rules.Field("name",
			rules.Assert("01", "name is required", is.Present),
		),
	}, extra...)
	return rules.For(new(ignoreItem), defs...)
}

// ignoreNameMustBeOK fails assertion "01" unless Name == "ok". Used to test
// When-gated ignores where the guard reads the same record.
func ignoreNameMustBeOK(extra ...rules.Def) *rules.Set {
	defs := append([]rules.Def{
		rules.Field("name",
			rules.Assert("01", "name must be 'ok'", is.Expr(`this == "ok"`)),
		),
	}, extra...)
	return rules.For(new(ignoreItem), defs...)
}

// ignoreGlobalDoc exercises rules.Ignore across two separately-registered
// namespaces through the global rules.Validate path.
type ignoreGlobalDoc struct {
	Name string `json:"name"`
}

func ignoreGlobalEmitterRules() *rules.Set {
	return rules.For(new(ignoreGlobalDoc),
		rules.Field("name",
			rules.Assert("01", "name is required", is.Present),
		),
	)
}

func ignoreGlobalSuppressorRules() *rules.Set {
	return rules.For(new(ignoreGlobalDoc),
		rules.Ignore("GOBL-IGN-EMIT-IGNOREGLOBALDOC-01"),
	)
}

func init() {
	// Suppressor is registered BEFORE the emitter on purpose: suppression
	// must not depend on the order in which namespaces run.
	rules.Register("ign-supp", rules.GOBL.Add("IGN-SUPP"), ignoreGlobalSuppressorRules())
	rules.Register("ign-emit", rules.GOBL.Add("IGN-EMIT"), ignoreGlobalEmitterRules())
}

func TestIgnore(t *testing.T) {
	t.Run("baseline: failing assert yields IGNOREITEM-01", func(t *testing.T) {
		faults := ignoreNameRequired().Validate(&ignoreItem{})
		require.NotNil(t, faults)
		assert.Equal(t, rules.Code("IGNOREITEM-01"), faults.First().Code())
	})

	t.Run("self-suppress: ignoring the failing code drops it", func(t *testing.T) {
		set := ignoreNameRequired(rules.Ignore("IGNOREITEM-01"))
		assert.Nil(t, set.Validate(&ignoreItem{}))
	})

	t.Run("non-match: ignoring a different code leaves the fault", func(t *testing.T) {
		set := ignoreNameRequired(rules.Ignore("GOBL-SOMETHING-ELSE-01"))
		faults := set.Validate(&ignoreItem{})
		require.NotNil(t, faults)
		assert.Equal(t, rules.Code("IGNOREITEM-01"), faults.First().Code())
	})

	t.Run("exact only: a prefix does not match", func(t *testing.T) {
		// Ignoring "IGNOREITEM" must NOT drop "IGNOREITEM-01".
		set := ignoreNameRequired(rules.Ignore("IGNOREITEM"))
		faults := set.Validate(&ignoreItem{})
		require.NotNil(t, faults)
		assert.Equal(t, rules.Code("IGNOREITEM-01"), faults.First().Code())
	})

	t.Run("when-gated off: ignore not collected when guard fails", func(t *testing.T) {
		set := ignoreNameRequired(
			rules.When(
				is.Expr(`Name == "extended"`),
				rules.Ignore("IGNOREITEM-01"),
			),
		)
		faults := set.Validate(&ignoreItem{Name: ""})
		require.NotNil(t, faults)
		assert.Equal(t, rules.Code("IGNOREITEM-01"), faults.First().Code())
	})

	t.Run("when-gated on: ignore collected when guard passes", func(t *testing.T) {
		set := ignoreNameMustBeOK(
			rules.When(
				is.Expr(`Name == "extended"`),
				rules.Ignore("IGNOREITEM-01"),
			),
		)
		// Name="extended" fails assert 01 (not "ok"), but the guard matches
		// so the fault is suppressed.
		assert.Nil(t, set.Validate(&ignoreItem{Name: "extended"}))
		// A different name still fails AND isn't suppressed (guard false).
		faults := set.Validate(&ignoreItem{Name: "nope"})
		require.NotNil(t, faults)
		assert.Equal(t, rules.Code("IGNOREITEM-01"), faults.First().Code())
	})

	t.Run("cross-namespace, order-independent suppression", func(t *testing.T) {
		// IGN-EMIT emits GOBL-IGN-EMIT-IGNOREGLOBALDOC-01; IGN-SUPP
		// (registered first) ignores that exact code. Global Validate runs
		// both and filters once at the end.
		assert.Nil(t, rules.Validate(&ignoreGlobalDoc{Name: ""}))
	})

	t.Run("WithIgnore option suppresses at the call site", func(t *testing.T) {
		set := ignoreNameRequired()
		// Without the option the fault is present.
		require.NotNil(t, set.Validate(&ignoreItem{}))
		// The option drops the matching code.
		assert.Nil(t, set.Validate(&ignoreItem{}, rules.WithIgnore("IGNOREITEM-01")))
		// A non-matching code leaves the fault in place.
		faults := set.Validate(&ignoreItem{}, rules.WithIgnore("GOBL-OTHER-01"))
		require.NotNil(t, faults)
		assert.Equal(t, rules.Code("IGNOREITEM-01"), faults.First().Code())
	})

	t.Run("sanity: emitter code is what the suppressor targets", func(t *testing.T) {
		// Confirm the registered emitter really produces the code the
		// suppressor ignores, so the cross-namespace test isn't a false pass.
		var emitter *rules.Set
		for _, ns := range rules.Registry() {
			if ns.ID == rules.Code("GOBL-IGN-EMIT") {
				emitter = ns
			}
		}
		require.NotNil(t, emitter)
		assert.Equal(t, rules.Code("GOBL-IGN-EMIT-IGNOREGLOBALDOC-01"),
			emitter.Subsets[0].Subsets[0].Assert[0].ID)
	})
}
