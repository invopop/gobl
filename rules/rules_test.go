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
	p := new(Person)
	return rules.ForStruct(p,
		rules.Object(
			rules.Assert("001", "person address must have a city",
				rules.Expr(`(address?.city ?? "") != ""`),
			),
		),
	)
}

func emailRules() *rules.Set {
	e := new(Email)
	return rules.ForStruct(e,
		rules.Field(&e.Addr,
			rules.Assert("010", "expected a valid email address",
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
		assert.Equal(t, "[GOBL-TEST-EMAIL-010] addr: expected a valid email address", faults.Error())
		assert.True(t, faults.HasPath("addr"))
		f := faults.First()
		assert.Equal(t, "expected a valid email address", f.Message())
		assert.Equal(t, rules.Code("GOBL-TEST-EMAIL-010"), f.Code())
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
		assert.Equal(t, "[GOBL-TEST-PERSON-001] person address must have a city; [GOBL-TEST-EMAIL-010] emails[1].addr: expected a valid email address", faults.Error())
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

// TestCode is a named string type used to test ForValue.
type TestCode string

func testCodeRules() *rules.Set {
	return rules.ForValue(TestCode(""),
		rules.Assert("001", "code must not be empty",
			rules.Required,
		),
		rules.Assert("002", "code must not exceed 10 characters",
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
		assert.Equal(t, rules.Code("TESTCODE-001"), faults.First().Code())
		assert.Equal(t, "code must not be empty", faults.First().Message())
	})

	t.Run("fails when too long", func(t *testing.T) {
		set := testCodeRules()
		faults := set.Validate(TestCode("ABCDEFGHIJK"))
		require.NotNil(t, faults)
		assert.Equal(t, 1, faults.Len())
		assert.Equal(t, rules.Code("TESTCODE-002"), faults.First().Code())
	})

	t.Run("global Validate finds value type rules", func(t *testing.T) {
		faults := rules.Validate(TestCode(""))
		require.NotNil(t, faults)
		assert.Equal(t, rules.Code("GOBL-TEST-TESTCODE-001"), faults.First().Code())
	})

	t.Run("global Validate passes valid value", func(t *testing.T) {
		faults := rules.Validate(TestCode("hello"))
		assert.Nil(t, faults)
	})

	t.Run("panics on invalid expression", func(t *testing.T) {
		assert.Panics(t, func() {
			rules.ForValue(TestCode(""),
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
		set := rules.ForStruct(new(Email),
			rules.Field(new(string),
				rules.Assert("001", "always fails",
					rules.Expr(`false`),
				),
			),
		)
		faults := set.Validate(&Other{X: "hello"})
		assert.Nil(t, faults)
	})

	t.Run("when condition skips set when false", func(t *testing.T) {
		type Item struct {
			Active bool   `json:"active"`
			Name   string `json:"name"`
		}
		proto := new(Item)
		set := rules.ForStruct(proto,
			rules.Field(&proto.Name,
				rules.Assert("001", "name required when active",
					rules.Required,
				),
			),
		).When(rules.Expr(`active`))
		// Active is false: set should be skipped.
		faults := set.Validate(&Item{Active: false, Name: ""})
		assert.NoError(t, faults)
	})

	t.Run("when condition runs assertions when true", func(t *testing.T) {
		type Item struct {
			Active bool   `json:"active"`
			Name   string `json:"name"`
		}
		proto := new(Item)
		inner := rules.ForStruct(proto,
			rules.Field(&proto.Name,
				rules.Assert("001", "name required when active",
					rules.Required,
				),
			),
		)
		set := inner.When(rules.Expr(`active`))
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
		inner := rules.ForStruct(proto,
			rules.Object(
				rules.Assert("001", "name required when active",
					rules.Expr(`name != ""`),
				),
			),
		)
		set := inner.When(rules.Expr(`active`))
		// Active is true, Name is blank: should fail.
		faults := set.Validate(&Item{Active: true, Name: ""})
		require.Error(t, faults)
	})
}
