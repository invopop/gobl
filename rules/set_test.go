package rules_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetValidate(t *testing.T) {
	t.Run("skips set when type does not match", func(t *testing.T) {
		type Other struct{ X string }
		set := rules.For(new(Email),
			rules.Assert("01", "always fails",
				is.Expr(`false`),
			),
		)
		faults := set.Validate(&Other{X: "hello"})
		assert.Nil(t, faults)
	})

	t.Run("provides object of same pointer type to tests", func(t *testing.T) {
		set := rules.For(new(Person),
			rules.Assert("01", "expected a valid name",
				is.Func("valid", func(value any) bool {
					p, ok := value.(*Person)
					if !ok {
						return false
					}
					return p.Name != ""
				}),
			),
		)
		p := &Person{Name: "fooo"}
		faults := set.Validate(p)
		assert.NoError(t, faults)

		op := Person{Name: "fooo"}
		faults = set.Validate(op)
		assert.NoError(t, faults)
	})

	t.Run("when condition skips set when false", func(t *testing.T) {
		type Item struct {
			Active bool   `json:"active"`
			Name   string `json:"name"`
		}
		set := rules.For(new(Item),
			rules.When(is.Expr(`Active`),
				rules.Field("name",
					rules.Assert("01", "name required when active",
						is.Present,
					),
				),
			),
		)
		// Active is false: set should be skipped.
		faults := set.Validate(&Item{Active: false, Name: ""})
		assert.NoError(t, faults)
	})

	t.Run("when condition runs assertions when true", func(t *testing.T) {
		type Item struct {
			Active bool   `json:"active"`
			Name   string `json:"name"`
		}
		set := rules.For(new(Item),
			rules.When(is.Expr(`Active`),
				rules.Field("name",
					rules.Assert("01", "name required when active",
						is.Present,
					),
				),
			),
		)
		// Active is true, Name is blank: should fail.
		faults := set.Validate(&Item{Active: true, Name: ""})
		require.Error(t, faults)
	})

	t.Run("when condition runs assertions when true with expr", func(t *testing.T) {
		type Item struct {
			Active bool   `json:"active"`
			Name   string `json:"name"`
		}
		proto := new(Item)
		set := rules.For(proto,
			rules.When(is.Expr(`Active`),
				rules.Assert("01", "name required when active",
					is.Expr(`Name != ""`),
				),
			),
		)
		// Active is true, Name is blank: should fail.
		faults := set.Validate(&Item{Active: true, Name: ""})
		require.Error(t, faults)
	})
}

func TestNewSet(t *testing.T) {
	t.Run("codes are namespaced", func(t *testing.T) {
		// Use top-level helper so runtime.Caller(1) in For detects the
		// correct package and omits the package prefix from the set ID.
		set := rules.NewSet("MYAPP", emailRules())
		faults := set.Validate(&Email{Addr: ""})
		require.Error(t, faults)
		assert.Equal(t, rules.Code("MYAPP-EMAIL-01"), faults.First().Code())
		assert.True(t, faults.HasPath("$.addr"))
	})

	t.Run("multi-type validation", func(t *testing.T) {
		set := rules.NewSet("NS", personRules(), emailRules())

		// Validates Person with nested Email via type index traversal.
		p := &Person{
			Name:    "Alice",
			Address: new(Address),
			Emails: []Email{
				{Addr: "ok@example.com"},
				{Addr: ""},
			},
		}
		faults := set.Validate(p)
		require.Error(t, faults)
		assert.True(t, faults.HasCode("NS-PERSON-01"))
		assert.True(t, faults.HasCode("NS-EMAIL-01"))
		assert.True(t, faults.HasPath("$.emails[1].addr"))
	})

	t.Run("does not pollute global registry", func(t *testing.T) {
		before := len(rules.AllSets())
		rules.NewSet("STANDALONE", emailRules())
		after := len(rules.AllSets())
		assert.Equal(t, before, after)
	})

	t.Run("passes when valid", func(t *testing.T) {
		set := rules.NewSet("OK", emailRules())
		faults := set.Validate(&Email{Addr: "test@example.com"})
		assert.NoError(t, faults)
	})

	t.Run("same For output reused across multiple NewSet calls", func(t *testing.T) {
		shared := emailRules()
		a := rules.NewSet("AAA", shared)
		b := rules.NewSet("BBB", shared)

		fa := a.Validate(&Email{Addr: ""})
		fb := b.Validate(&Email{Addr: ""})
		require.Error(t, fa)
		require.Error(t, fb)
		assert.Equal(t, rules.Code("AAA-EMAIL-01"), fa.First().Code())
		assert.Equal(t, rules.Code("BBB-EMAIL-01"), fb.First().Code())
		// Original set is untouched.
		assert.Equal(t, rules.Code("EMAIL"), shared.ID)
	})
}

func docWithRegimeRules() *rules.Set {
	return rules.For(new(docWithRegime),
		rules.When(is.InContext(is.In("ES")),
			rules.Field("name",
				rules.Assert("01", "name required", is.Present),
			),
		),
	)
}

func TestNewSetValidateWithContext(t *testing.T) {
	set := rules.NewSet("CTX", docWithRegimeRules())

	t.Run("context-aware guard fires with matching context", func(t *testing.T) {
		doc := &docWithRegime{
			Regime: testRegime{Code: "ES"},
			Name:   "",
		}
		faults := set.Validate(doc)
		require.Error(t, faults)
		assert.True(t, faults.HasCode("CTX-DOCWITHREGIME-01"))
	})

	t.Run("context-aware guard skips with non-matching context", func(t *testing.T) {
		doc := &docWithRegime{
			Regime: testRegime{Code: "FR"},
			Name:   "",
		}
		faults := set.Validate(doc)
		assert.NoError(t, faults)
	})
}

func TestSetMarshalJSON(t *testing.T) {
	set := rules.For(new(Person),
		rules.When(is.Expr(`Age > 0`),
			rules.Assert("01", "name required",
				is.Expr(`Name != ""`),
			),
		),
		rules.Field("name",
			rules.Assert("02", "name not empty", is.Present),
		),
	)
	data, err := json.Marshal(set)
	require.NoError(t, err)

	var m map[string]any
	require.NoError(t, json.Unmarshal(data, &m))
	assert.Contains(t, m, "id")
	assert.Contains(t, m, "subsets")
}

func TestAssertionMarshalJSON(t *testing.T) {
	set := rules.For(new(Person),
		rules.Assert("01", "name and age valid",
			is.Expr(`Name != ""`),
			is.Expr(`Age > 0`),
		),
	)
	data, err := json.Marshal(set)
	require.NoError(t, err)

	var m map[string]any
	require.NoError(t, json.Unmarshal(data, &m))
	asserts, ok := m["assert"].([]any)
	require.True(t, ok)
	require.Len(t, asserts, 1)
	a := asserts[0].(map[string]any)
	// Tests should be comma-joined string
	tests, ok := a["tests"].(string)
	require.True(t, ok)
	assert.Contains(t, tests, ", ")
}

func TestAssertIfPresentMarshalJSON(t *testing.T) {
	// This triggers presentGuard.String() via Set.MarshalJSON.
	set := rules.For(new(Person),
		rules.Field("name",
			rules.AssertIfPresent("01", "name too short",
				is.Func("min 3", func(val any) bool {
					s, ok := val.(string)
					return ok && len(s) >= 3
				}),
			),
		),
	)
	data, err := json.Marshal(set)
	require.NoError(t, err)
	assert.Contains(t, string(data), "present")
}
