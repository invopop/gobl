package rules_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
}

type Person struct {
	Name          string   `json:"name"`
	Age           int      `json:"age"`
	Address       *Address `json:"address"`
	SecondAddress *Address `json:"second_address,omitempty"`
	Emails        []Email  `json:"emails,omitempty"`
}

type Email struct {
	Addr string `json:"addr"`
}

func personRules() *rules.Set {
	return rules.For(new(Person),
		rules.Assert("01", "person address must have a city",
			is.Expr(`(address?.city ?? "") != ""`),
		),
	)
}

func emailRules() *rules.Set {
	return rules.For(new(Email),
		rules.Field("addr",
			rules.Assert("01", "expected a valid email address",
				is.Present,
				is.EmailFormat,
			),
		),
	)
}

func init() {
	rules.Register(
		"test",
		rules.GOBL.Add("TEST"),
		emailRules(),
		personRules(),
	)
}

func TestFor(t *testing.T) {
	t.Run("name includes go package name", func(t *testing.T) {
		set := emailRules()
		assert.Equal(t, "rules.Email", set.Name)
	})

	t.Run("id omits package when For is called from the same package", func(t *testing.T) {
		// emailRules() is defined and called in package rules_test, same as Email.
		set := emailRules()
		assert.Equal(t, rules.Code("EMAIL"), set.ID)
	})

	t.Run("id gets namespace prepended by register", func(t *testing.T) {
		// emailRules() is registered under GOBL-TEST in init(), so the
		// global registry holds a subset with ID "GOBL-TEST-EMAIL".
		var found *rules.Set
		for _, ns := range rules.Registry() {
			for _, sub := range ns.Subsets {
				if sub.Name == "rules.Email" {
					found = sub
					break
				}
			}
		}
		require.NotNil(t, found, "expected to find rules.Email set in registry")
		assert.Equal(t, rules.Code("GOBL-TEST-EMAIL"), found.ID)
	})

	t.Run("assertion code includes full namespace after registration", func(t *testing.T) {
		e := &Email{Addr: ""}
		faults := rules.Validate(e)
		require.NotNil(t, faults)
		assert.Equal(t, rules.Code("GOBL-TEST-EMAIL-01"), faults.First().Code())
	})

	t.Run("panics on invalid field name", func(t *testing.T) {
		assert.Panics(t, func() {
			rules.For(new(Email),
				rules.Field("nonexistent"),
			)
		})
	})
}

func TestValidate(t *testing.T) {
	t.Run("passes with valid email", func(t *testing.T) {
		e := &Email{Addr: "test@example.com"}
		faults := rules.Validate(e)
		assert.NoError(t, faults)
	})

	t.Run("fails with blank email address", func(t *testing.T) {
		e := &Email{Addr: ""}
		faults := rules.Validate(e)
		require.NotNil(t, faults)
		assert.Equal(t, "[GOBL-TEST-EMAIL-01] ($.addr) expected a valid email address", faults.Error())
		assert.True(t, faults.HasPath("$.addr"))
		f := faults.First()
		assert.Equal(t, "expected a valid email address", f.Message())
		assert.Equal(t, rules.Code("GOBL-TEST-EMAIL-01"), f.Code())
		assert.Equal(t, "$.addr", f.Path())
	})

	t.Run("recurses into struct fields", func(t *testing.T) {
		p := &Person{
			Name: "Alice",
			Emails: []Email{
				{Addr: "ok@example.com"},
				{Addr: ""},
			},
			Address: new(Address),
		}
		faults := rules.Validate(p)
		require.NotNil(t, faults)
		assert.True(t, faults.HasPath("$.emails[1].addr"), "expected fault at $.emails[1].addr")
		assert.Equal(t, "[GOBL-TEST-PERSON-01] person address must have a city; [GOBL-TEST-EMAIL-01] ($.emails[1].addr) expected a valid email address", faults.Error())
	})

	t.Run("recurses into pointer fields", func(t *testing.T) {
		p := &Person{
			Name:    "Bob",
			Address: &Address{Street: "123 Main St", City: "Anytown"},
		}
		// Address has no rules, should pass.
		faults := rules.Validate(p)
		assert.Nil(t, faults)
	})

	t.Run("nil pointer field is skipped", func(t *testing.T) {
		p := &Person{Name: "Carol", Address: &Address{City: "Anytown"}}
		faults := rules.Validate(p)
		assert.Nil(t, faults)
	})

}

// TestCode is a named string type used to test For with a value type.
type TestCode string

func testCodeRules() *rules.Set {
	return rules.For(TestCode(""),
		rules.Assert("01", "code must not be empty",
			is.Present,
		),
		rules.Assert("02", "code must not exceed 10 characters",
			is.Expr(`len(this) <= 10`),
		),
	)
}

func init() {
	rules.Register(
		"test",
		rules.GOBL.Add("TEST"),
		testCodeRules(),
	)
}

// Tag is a named string type used as a map key to test map key validation.
type Tag string

type Tagged struct {
	Tags map[Tag]Email `json:"tags,omitempty"`
}

func tagRules() *rules.Set {
	return rules.For(Tag(""),
		rules.Assert("01", "tag must not be empty", is.Present),
	)
}

func init() {
	rules.Register("test", rules.GOBL.Add("TEST"), tagRules())
}

func TestForValue(t *testing.T) {
	t.Run("passes with valid code", func(t *testing.T) {
		set := testCodeRules()
		faults := set.Validate(TestCode("ABC"))
		assert.Nil(t, faults)
	})

	t.Run("fails when empty", func(t *testing.T) {
		set := testCodeRules()
		faults := set.Validate(TestCode(""))
		require.Error(t, faults)
		assert.Equal(t, 1, faults.Len())
		assert.Equal(t, rules.Code("TESTCODE-01"), faults.First().Code())
		assert.Equal(t, "code must not be empty", faults.First().Message())
	})

	t.Run("fails when too long", func(t *testing.T) {
		set := testCodeRules()
		faults := set.Validate(TestCode("ABCDEFGHIJK"))
		require.NotNil(t, faults)
		assert.Equal(t, 1, faults.Len())
		assert.Equal(t, rules.Code("TESTCODE-02"), faults.First().Code())
	})

	t.Run("global Validate finds value type rules", func(t *testing.T) {
		faults := rules.Validate(TestCode(""))
		require.NotNil(t, faults)
		assert.Equal(t, rules.Code("GOBL-TEST-TESTCODE-01"), faults.First().Code())
	})

	t.Run("global Validate passes valid value", func(t *testing.T) {
		faults := rules.Validate(TestCode("hello"))
		assert.Nil(t, faults)
	})

	t.Run("panics on invalid expression", func(t *testing.T) {
		assert.Panics(t, func() {
			rules.For(TestCode(""),
				rules.Assert("bad", "bad expr",
					is.Expr(`this ===`),
				),
			)
		})
	})

}

func TestMapValidation(t *testing.T) {
	t.Run("traverses map values and reports keyed paths", func(t *testing.T) {
		obj := &Tagged{
			Tags: map[Tag]Email{
				"work": {Addr: "work@example.com"},
				"home": {Addr: ""},
			},
		}
		faults := rules.Validate(obj)
		require.NotNil(t, faults)
		assert.True(t, faults.HasPath("$.tags.home.addr"), "expected fault at $.tags.home.addr")
		assert.False(t, faults.HasPath("$.tags.work.addr"), "expected no fault at $.tags.work.addr")
	})

	t.Run("validates named map keys", func(t *testing.T) {
		obj := &Tagged{
			Tags: map[Tag]Email{
				"":     {Addr: "ok@example.com"},
				"work": {Addr: "work@example.com"},
			},
		}
		faults := rules.Validate(obj)
		require.NotNil(t, faults)
		assert.True(t, faults.HasPath("$.tags"), "expected fault at $.tags for empty key")
	})

	t.Run("key order is deterministic across multiple failing entries", func(t *testing.T) {
		obj := &Tagged{
			Tags: map[Tag]Email{
				"beta":  {Addr: ""},
				"alpha": {Addr: ""},
			},
		}
		faults := rules.Validate(obj)
		require.NotNil(t, faults)
		require.GreaterOrEqual(t, faults.Len(), 2)
		assert.Equal(t, "$.tags.alpha.addr", faults.At(0).Path())
		assert.Equal(t, "$.tags.beta.addr", faults.At(1).Path())
	})

	t.Run("nil map is skipped", func(t *testing.T) {
		obj := &Tagged{Tags: nil}
		faults := rules.Validate(obj)
		assert.Nil(t, faults)
	})

	t.Run("empty map passes", func(t *testing.T) {
		obj := &Tagged{Tags: map[Tag]Email{}}
		faults := rules.Validate(obj)
		assert.Nil(t, faults)
	})
}

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
			rules.When(is.Expr(`active`),
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
			rules.When(is.Expr(`active`),
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
			rules.When(is.Expr(`active`),
				rules.Assert("01", "name required when active",
					is.Expr(`name != ""`),
				),
			),
		)
		// Active is true, Name is blank: should fail.
		faults := set.Validate(&Item{Active: true, Name: ""})
		require.Error(t, faults)
	})

}

func TestFieldEmpty(t *testing.T) {
	set := rules.For(new(Person),
		rules.Field("name",
			rules.Assert("01", "name must not be set", is.Empty),
		),
	)

	t.Run("passes when field is empty string", func(t *testing.T) {
		faults := set.Validate(&Person{Name: ""})
		assert.NoError(t, faults)
	})

	t.Run("fails when field has a value", func(t *testing.T) {
		faults := set.Validate(&Person{Name: "Alice"})
		require.Error(t, faults)
		assert.True(t, faults.HasPath("$.name"))
	})
}

func TestFieldNil(t *testing.T) {
	set := rules.For(new(Person),
		rules.Field("second_address",
			rules.Assert("01", "second address must not be set", is.Nil),
		),
	)

	t.Run("passes when pointer field is nil", func(t *testing.T) {
		faults := set.Validate(&Person{SecondAddress: nil})
		assert.NoError(t, faults)
	})

	t.Run("fails when pointer field is non-nil", func(t *testing.T) {
		faults := set.Validate(&Person{SecondAddress: &Address{City: "London"}})
		require.Error(t, faults)
		assert.True(t, faults.HasPath("$.second_address"))
	})

	t.Run("fails when pointer field is non-nil but points to empty value", func(t *testing.T) {
		faults := set.Validate(&Person{SecondAddress: new(Address)})
		require.Error(t, faults)
		assert.True(t, faults.HasPath("$.second_address"))
	})
}

func TestFieldRulesDoNotBleedToSameType(t *testing.T) {
	// Verifies that rules defined inside a When+Field block for a specific field
	// (e.g. "address") are NOT applied to other fields of the same type
	// (e.g. "second_address"). This guards against the namespace traversal
	// incorrectly matching field subsets by type rather than by field name.
	set := rules.For(new(Person),
		rules.When(
			is.Expr(`name != ""`),
			rules.Field("address",
				rules.Assert("01", "city is required",
					is.Func("city required", func(val any) bool {
						a, ok := val.(*Address)
						return !ok || a == nil || a.City != ""
					}),
				),
			),
		),
	)

	t.Run("rule is not applied to second_address field of same type", func(t *testing.T) {
		p := &Person{
			Name:          "Alice",
			Address:       &Address{City: "London"}, // valid
			SecondAddress: &Address{City: ""},       // would fail city check if rule bled
		}
		faults := set.Validate(p)
		assert.NoError(t, faults)
	})

	t.Run("rule still fails when the scoped address field fails", func(t *testing.T) {
		p := &Person{
			Name:          "Alice",
			Address:       &Address{City: ""},       // should fail
			SecondAddress: &Address{City: "London"}, // should not be checked
		}
		faults := set.Validate(p)
		require.Error(t, faults)
		assert.True(t, faults.HasPath("$.address"))
		assert.False(t, faults.HasPath("$.second_address"))
	})
}

func TestEach(t *testing.T) {
	t.Run("validates each element and reports indexed paths", func(t *testing.T) {
		set := rules.For(new(Person),
			rules.Field("emails",
				rules.Each(
					rules.Field("addr",
						rules.Assert("01", "email address is required", is.Present),
					),
				),
			),
		)
		p := &Person{
			Name: "Alice",
			Emails: []Email{
				{Addr: "a@example.com"},
				{Addr: ""},
				{Addr: "b@example.com"},
				{Addr: ""},
			},
		}
		faults := set.Validate(p)
		require.Error(t, faults)
		assert.True(t, faults.HasPath("$.emails[1].addr"))
		assert.True(t, faults.HasPath("$.emails[3].addr"))
		assert.False(t, faults.HasPath("$.emails[0].addr"))
		assert.False(t, faults.HasPath("$.emails[2].addr"))
	})

	t.Run("passes when all elements are valid", func(t *testing.T) {
		set := rules.For(new(Person),
			rules.Field("emails",
				rules.Each(
					rules.Field("addr",
						rules.Assert("01", "email address is required", is.Present),
					),
				),
			),
		)
		p := &Person{
			Emails: []Email{{Addr: "a@example.com"}, {Addr: "b@example.com"}},
		}
		assert.NoError(t, set.Validate(p))
	})

	t.Run("passes with empty slice", func(t *testing.T) {
		set := rules.For(new(Person),
			rules.Field("emails",
				rules.Each(
					rules.Assert("01", "required", is.Present),
				),
			),
		)
		assert.NoError(t, set.Validate(&Person{}))
	})

	t.Run("whole-slice and per-element assertions coexist on same field", func(t *testing.T) {
		set := rules.For(new(Person),
			rules.Field("emails",
				rules.Assert("01", "no more than two emails",
					is.Func("max two", func(val any) bool {
						emails, ok := val.([]Email)
						return !ok || len(emails) <= 2
					}),
				),
				rules.Each(
					rules.Field("addr",
						rules.Assert("02", "email address is required", is.Present),
					),
				),
			),
		)
		// Whole-slice violation: three emails.
		p := &Person{Emails: []Email{{Addr: "a@b.com"}, {Addr: "c@d.com"}, {Addr: "e@f.com"}}}
		faults := set.Validate(p)
		require.Error(t, faults)
		assert.True(t, faults.HasPath("$.emails"))

		// Per-element violation: empty addr.
		p2 := &Person{Emails: []Email{{Addr: "a@b.com"}, {Addr: ""}}}
		faults2 := set.Validate(p2)
		require.Error(t, faults2)
		assert.True(t, faults2.HasPath("$.emails[1].addr"))
	})

	t.Run("panics when used outside a slice field", func(t *testing.T) {
		assert.Panics(t, func() {
			rules.For(new(Email),
				rules.Each(
					rules.Assert("01", "required", is.Present),
				),
			)
		})
	})
}

// Wrapper and Inner types for Embeddable tests.
type Inner struct {
	Name string `json:"name"`
}

type Wrapper struct {
	inner *Inner
}

func (w *Wrapper) Embedded() any {
	return w.inner
}

func innerRules() *rules.Set {
	return rules.For(new(Inner),
		rules.Field("name",
			rules.Assert("01", "inner name is required", is.Present),
		),
	)
}

// Container holds a Wrapper in a JSON field to test path prefixing.
type Container struct {
	Doc *Wrapper `json:"doc"`
}

func init() {
	rules.Register("embeddable-test", rules.GOBL.Add("EMBTEST"), innerRules())
}

func TestEmbeddable(t *testing.T) {
	t.Run("faults from embedded payload carry no extra prefix", func(t *testing.T) {
		w := &Wrapper{inner: &Inner{Name: ""}}
		set := innerRules()
		faults := set.Validate(w)
		// The set is scoped to Inner; Wrapper itself has no matching rules.
		assert.Nil(t, faults)
	})

	t.Run("global Validate traverses Embeddable and prefixes path correctly", func(t *testing.T) {
		c := &Container{Doc: &Wrapper{inner: &Inner{Name: ""}}}
		faults := rules.Validate(c)
		require.Error(t, faults)
		assert.True(t, faults.HasPath("$.doc.name"), "expected fault at $.doc.name, got: %v", faults)
	})

	t.Run("global Validate passes when embedded payload is valid", func(t *testing.T) {
		c := &Container{Doc: &Wrapper{inner: &Inner{Name: "Alice"}}}
		faults := rules.Validate(c)
		assert.Nil(t, faults)
	})

	t.Run("nil embedded payload is skipped", func(t *testing.T) {
		c := &Container{Doc: &Wrapper{inner: nil}}
		faults := rules.Validate(c)
		assert.Nil(t, faults)
	})
}

func TestAssertIfPresent(t *testing.T) {
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

	t.Run("non-empty value runs assertion and passes", func(t *testing.T) {
		faults := set.Validate(&Person{Name: "Alice"})
		assert.NoError(t, faults)
	})

	t.Run("non-empty value runs assertion and fails", func(t *testing.T) {
		faults := set.Validate(&Person{Name: "Al"})
		require.Error(t, faults)
		assert.True(t, faults.HasPath("$.name"))
	})

	t.Run("empty string skips assertion", func(t *testing.T) {
		faults := set.Validate(&Person{Name: ""})
		assert.NoError(t, faults)
	})

	t.Run("nil pointer skips assertion", func(t *testing.T) {
		ptrSet := rules.For(new(Person),
			rules.Field("address",
				rules.AssertIfPresent("01", "street required",
					is.Func("has street", func(val any) bool {
						a, ok := val.(*Address)
						return ok && a.Street != ""
					}),
				),
			),
		)
		faults := ptrSet.Validate(&Person{Address: nil})
		assert.NoError(t, faults)
	})
}

func TestObject(t *testing.T) {
	set := rules.For(new(Person),
		rules.Object(
			rules.Assert("01", "name required",
				is.Expr(`name != ""`),
			),
			rules.Assert("02", "age positive",
				is.Expr(`age > 0`),
			),
		),
	)

	t.Run("both pass", func(t *testing.T) {
		faults := set.Validate(&Person{Name: "Alice", Age: 30})
		assert.NoError(t, faults)
	})

	t.Run("both fail", func(t *testing.T) {
		faults := set.Validate(&Person{Name: "", Age: 0})
		require.Error(t, faults)
		assert.GreaterOrEqual(t, faults.Len(), 2)
	})
}

func TestAllSets(t *testing.T) {
	sets := rules.AllSets()
	assert.NotEmpty(t, sets)
	// Should be same as Registry.
	assert.Equal(t, rules.Registry(), sets)
}

func TestSetMarshalJSON(t *testing.T) {
	set := rules.For(new(Person),
		rules.When(is.Expr(`age > 0`),
			rules.Assert("01", "name required",
				is.Expr(`name != ""`),
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
			is.Expr(`name != ""`),
			is.Expr(`age > 0`),
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

func TestValidateNil(t *testing.T) {
	faults := rules.Validate(nil)
	assert.Nil(t, faults)
}

// testRegime implements ContextAdder for guard testing.
type testRegime struct {
	Code string
}

func (r testRegime) RulesContext() rules.WithContext {
	return func(rc *rules.Context) {
		rc.Set("regime", r.Code)
	}
}

type docWithRegime struct {
	Regime testRegime `json:"regime"`
	Name   string     `json:"name"`
}

func TestRegisterWithGuard(t *testing.T) {
	guardTest := is.In("ES", "PT")
	nameSet := rules.For(new(docWithRegime),
		rules.Field("name",
			rules.Assert("01", "name required", is.Present),
		),
	)

	rules.RegisterWithGuard("guard-test", rules.GOBL.Add("GUARDTEST"),
		is.HasContext(guardTest), nameSet)

	t.Run("guard passes - set applied", func(t *testing.T) {
		doc := &docWithRegime{
			Regime: testRegime{Code: "ES"},
			Name:   "",
		}
		faults := rules.Validate(doc)
		require.Error(t, faults)
		assert.True(t, faults.HasCode("GOBL-GUARDTEST-DOCWITHREGIME-01"))
	})

	t.Run("guard fails - set skipped", func(t *testing.T) {
		doc := &docWithRegime{
			Regime: testRegime{Code: "FR"},
			Name:   "",
		}
		faults := rules.Validate(doc)
		// Should not have the guard-test fault (guard doesn't match FR)
		if faults != nil {
			assert.False(t, faults.HasCode("GOBL-GUARDTEST-DOCWITHREGIME-01"))
		}
	})
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

func TestValidateWithContext(t *testing.T) {
	// Tests the WithContext option path in Validate.
	guardTest := is.In("DE")
	nameSet := rules.For(new(Email),
		rules.Field("addr",
			rules.Assert("01", "email required", is.Present),
		),
	)
	rules.RegisterWithGuard("ctx-opt-test", rules.GOBL.Add("CTXOPT"),
		is.HasContext(guardTest), nameSet)

	t.Run("WithContext option injects context", func(t *testing.T) {
		e := &Email{Addr: ""}
		faults := rules.Validate(e, func(rc *rules.Context) {
			rc.Set("country", "DE")
		})
		require.Error(t, faults)
		assert.True(t, faults.HasCode("GOBL-CTXOPT-EMAIL-01"))
	})
}

func TestEachWithPointerElements(t *testing.T) {
	type Item struct {
		Value string `json:"value"`
	}
	type Container struct {
		Items []*Item `json:"items"`
	}

	set := rules.For(new(Container),
		rules.Field("items",
			rules.Each(
				rules.Field("value",
					rules.Assert("01", "value required", is.Present),
				),
			),
		),
	)

	t.Run("pointer elements validated", func(t *testing.T) {
		c := &Container{Items: []*Item{{Value: "ok"}, {Value: ""}}}
		faults := set.Validate(c)
		require.Error(t, faults)
		assert.True(t, faults.HasPath("$.items[1].value"))
	})

	t.Run("all valid", func(t *testing.T) {
		c := &Container{Items: []*Item{{Value: "a"}, {Value: "b"}}}
		faults := set.Validate(c)
		assert.NoError(t, faults)
	})

	t.Run("nil slice field", func(t *testing.T) {
		c := &Container{Items: nil}
		faults := set.Validate(c)
		assert.NoError(t, faults)
	})
}

func TestEachTopLevelSlice(t *testing.T) {
	// Tests validateEachValue with pointer-to-slice and Each at top-level of When.
	type Row struct {
		Name string `json:"name"`
	}
	type Table struct {
		Rows []Row `json:"rows"`
	}
	set := rules.For(new(Table),
		rules.Field("rows",
			rules.Each(
				rules.Assert("01", "name required",
					is.Func("has name", func(val any) bool {
						r, ok := val.(*Row)
						return ok && r.Name != ""
					}),
				),
			),
		),
	)
	t.Run("each validates all rows", func(t *testing.T) {
		tb := &Table{Rows: []Row{{Name: "a"}, {Name: ""}}}
		faults := set.Validate(tb)
		require.Error(t, faults)
		assert.True(t, faults.HasPath("$.rows[1]"))
	})
}

func TestValidateWithNilPointerField(_ *testing.T) {
	// Nil pointer field should be skipped without panic.
	p := &Person{Name: "Alice", Address: nil}
	faults := rules.Validate(p)
	// No rules fire since address is nil and name is set.
	_ = faults
}

type noJSONTag struct {
	Bare string
}

func TestFieldWithNoJSONTag(t *testing.T) {
	// Tests the jsonFieldName fallback to field name when no json tag.
	set := rules.For(new(noJSONTag),
		rules.Field("Bare",
			rules.Assert("01", "bare required", is.Present),
		),
	)
	faults := set.Validate(&noJSONTag{Bare: ""})
	require.Error(t, faults)
	assert.True(t, faults.HasPath("$.Bare"))
}

type dashTag struct {
	Hidden string `json:"-"`
	Shown  string `json:"shown"`
}

func TestFieldWithDashJSONTag(t *testing.T) {
	// Exercises the json:"-" path in jsonFieldName.
	set := rules.For(new(dashTag),
		rules.Field("shown",
			rules.Assert("01", "shown required", is.Present),
		),
	)
	faults := set.Validate(&dashTag{Shown: ""})
	require.Error(t, faults)
}
