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
			rules.Expr(`(address?.city ?? "") != ""`),
		),
	)
}

func emailRules() *rules.Set {
	return rules.For(new(Email),
		rules.Field("addr",
			rules.Assert("01", "expected a valid email address",
				rules.Required,
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

	t.Run("id is pkg-typename before registration", func(t *testing.T) {
		// Email is in package rules_test (short name "rules"), so the base ID is RULES-EMAIL.
		set := emailRules()
		assert.Equal(t, rules.Code("RULES-EMAIL"), set.ID)
	})

	t.Run("id gets namespace prepended by register", func(t *testing.T) {
		// emailRules() is registered under GOBL-TEST in init(), so the
		// global registry holds a subset with ID "GOBL-TEST-RULES-EMAIL".
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
		assert.Equal(t, rules.Code("GOBL-TEST-RULES-EMAIL"), found.ID)
	})

	t.Run("assertion code includes full namespace after registration", func(t *testing.T) {
		e := &Email{Addr: ""}
		faults := rules.Validate(e)
		require.NotNil(t, faults)
		assert.Equal(t, rules.Code("GOBL-TEST-RULES-EMAIL-01"), faults.First().Code())
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
		assert.Equal(t, "[GOBL-TEST-RULES-EMAIL-01] addr: expected a valid email address", faults.Error())
		assert.True(t, faults.HasPath("addr"))
		f := faults.First()
		assert.Equal(t, "expected a valid email address", f.Message())
		assert.Equal(t, rules.Code("GOBL-TEST-RULES-EMAIL-01"), f.Code())
		assert.Equal(t, "addr", f.Path())
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
		assert.True(t, faults.HasPath("emails[1].addr"), "expected fault at emails[1].addr")
		assert.Equal(t, "[GOBL-TEST-RULES-PERSON-01] person address must have a city; [GOBL-TEST-RULES-EMAIL-01] emails[1].addr: expected a valid email address", faults.Error())
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
			rules.Required,
		),
		rules.Assert("02", "code must not exceed 10 characters",
			rules.Expr(`len(this) <= 10`),
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
		assert.Equal(t, rules.Code("RULES-TESTCODE-01"), faults.First().Code())
		assert.Equal(t, "code must not be empty", faults.First().Message())
	})

	t.Run("fails when too long", func(t *testing.T) {
		set := testCodeRules()
		faults := set.Validate(TestCode("ABCDEFGHIJK"))
		require.NotNil(t, faults)
		assert.Equal(t, 1, faults.Len())
		assert.Equal(t, rules.Code("RULES-TESTCODE-02"), faults.First().Code())
	})

	t.Run("global Validate finds value type rules", func(t *testing.T) {
		faults := rules.Validate(TestCode(""))
		require.NotNil(t, faults)
		assert.Equal(t, rules.Code("GOBL-TEST-RULES-TESTCODE-01"), faults.First().Code())
	})

	t.Run("global Validate passes valid value", func(t *testing.T) {
		faults := rules.Validate(TestCode("hello"))
		assert.Nil(t, faults)
	})

	t.Run("panics on invalid expression", func(t *testing.T) {
		assert.Panics(t, func() {
			rules.For(TestCode(""),
				rules.Assert("bad", "bad expr",
					rules.Expr(`this ===`),
				),
			)
		})
	})

}

func TestSetValidate(t *testing.T) {
	t.Run("skips set when type does not match", func(t *testing.T) {
		type Other struct{ X string }
		set := rules.For(new(Email),
			rules.Assert("01", "always fails",
				rules.Expr(`false`),
			),
		)
		faults := set.Validate(&Other{X: "hello"})
		assert.Nil(t, faults)
	})

	t.Run("provides object of same pointer type to tests", func(t *testing.T) {
		set := rules.For(new(Person),
			rules.Assert("01", "expected a valid name",
				rules.By("valid", func(value any) bool {
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
			rules.When(rules.Expr(`active`),
				rules.Field("name",
					rules.Assert("01", "name required when active",
						rules.Required,
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
			rules.When(rules.Expr(`active`),
				rules.Field("name",
					rules.Assert("01", "name required when active",
						rules.Required,
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
			rules.When(rules.Expr(`active`),
				rules.Assert("01", "name required when active",
					rules.Expr(`name != ""`),
				),
			),
		)
		// Active is true, Name is blank: should fail.
		faults := set.Validate(&Item{Active: true, Name: ""})
		require.Error(t, faults)
	})

}

func TestEach(t *testing.T) {
	t.Run("validates each element and reports indexed paths", func(t *testing.T) {
		set := rules.For(new(Person),
			rules.Field("emails",
				rules.Each(
					rules.Field("addr",
						rules.Assert("01", "email address is required", rules.Required),
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
		assert.True(t, faults.HasPath("emails[1].addr"))
		assert.True(t, faults.HasPath("emails[3].addr"))
		assert.False(t, faults.HasPath("emails[0].addr"))
		assert.False(t, faults.HasPath("emails[2].addr"))
	})

	t.Run("passes when all elements are valid", func(t *testing.T) {
		set := rules.For(new(Person),
			rules.Field("emails",
				rules.Each(
					rules.Field("addr",
						rules.Assert("01", "email address is required", rules.Required),
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
					rules.Assert("01", "required", rules.Required),
				),
			),
		)
		assert.NoError(t, set.Validate(&Person{}))
	})

	t.Run("whole-slice and per-element assertions coexist on same field", func(t *testing.T) {
		set := rules.For(new(Person),
			rules.Field("emails",
				rules.Assert("01", "no more than two emails",
					rules.By("max two", func(val any) bool {
						emails, ok := val.([]Email)
						return !ok || len(emails) <= 2
					}),
				),
				rules.Each(
					rules.Field("addr",
						rules.Assert("02", "email address is required", rules.Required),
					),
				),
			),
		)
		// Whole-slice violation: three emails.
		p := &Person{Emails: []Email{{Addr: "a@b.com"}, {Addr: "c@d.com"}, {Addr: "e@f.com"}}}
		faults := set.Validate(p)
		require.Error(t, faults)
		assert.True(t, faults.HasPath("emails"))

		// Per-element violation: empty addr.
		p2 := &Person{Emails: []Email{{Addr: "a@b.com"}, {Addr: ""}}}
		faults2 := set.Validate(p2)
		require.Error(t, faults2)
		assert.True(t, faults2.HasPath("emails[1].addr"))
	})

	t.Run("panics when used outside a slice field", func(t *testing.T) {
		assert.Panics(t, func() {
			rules.For(new(Email),
				rules.Each(
					rules.Assert("01", "required", rules.Required),
				),
			)
		})
	})
}
