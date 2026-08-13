package rules_test

import (
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			is.Expr(`Name != ""`),
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

func TestFieldNested(t *testing.T) {
	set := rules.For(new(Person),
		rules.Field("address",
			rules.Field("city",
				rules.Assert("01", "city is required", is.Present),
			),
			rules.Field("street",
				rules.Assert("02", "street is required", is.Present),
			),
		),
	)

	t.Run("reports fault at nested field path", func(t *testing.T) {
		p := &Person{Address: &Address{Street: "1 Main", City: ""}}
		faults := set.Validate(p)
		require.Error(t, faults)
		assert.True(t, faults.HasPath("$.address.city"))
	})

	t.Run("passes when nested fields are set", func(t *testing.T) {
		p := &Person{Address: &Address{Street: "1 Main", City: "London"}}
		assert.NoError(t, set.Validate(p))
	})

	t.Run("skips nested field subset when parent pointer is nil", func(t *testing.T) {
		p := &Person{Address: nil}
		assert.NoError(t, set.Validate(p))
	})

	t.Run("sibling nested fields are independent", func(t *testing.T) {
		p := &Person{Address: &Address{Street: "1 Main", City: ""}}
		faults := set.Validate(p)
		require.Error(t, faults)
		assert.True(t, faults.HasPath("$.address.city"))
		assert.False(t, faults.HasPath("$.address.street"))
	})
}

func TestEmbeddedStruct(t *testing.T) {
	type Individual struct {
		Person
		Testing bool `json:"testing"`
	}
	set := rules.For(new(Individual),
		rules.Field("address",
			rules.Field("city",
				rules.Assert("01", "city is required", is.Present),
			),
		),
		rules.Field("testing",
			rules.Assert("02", "testing must be true", is.In(true)),
		),
	)
	t.Run("embedded struct fields are validated", func(t *testing.T) {
		i := new(Individual)
		i.Address = &Address{Street: "1 Main", City: "Foo"}
		i.Testing = true
		assert.NoError(t, set.Validate(i))
	})
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
