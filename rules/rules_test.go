package rules_test

import (
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
			is.Expr(`(Address?.City ?? "") != ""`),
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
	rules.Register(
		"test",
		rules.GOBL.Add("TEST"),
		emailRules(),
		personRules(),
		testCodeRules(),
		tagRules(),
	)
	rules.Register("embeddable-test", rules.GOBL.Add("EMBTEST"), innerRules())
}

func TestFor(t *testing.T) {
	t.Run("object includes go package name", func(t *testing.T) {
		set := emailRules()
		assert.Equal(t, "rules.Email", set.Object)
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
				if sub.Object == "rules.Email" {
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
				is.Expr(`Name != ""`),
			),
			rules.Assert("02", "age positive",
				is.Expr(`Age > 0`),
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
	// Should contain the same sets as Registry (order may vary due to map iteration).
	reg := rules.Registry()
	assert.Equal(t, len(reg), len(sets))
	for _, s := range reg {
		assert.Contains(t, sets, s)
	}
}
