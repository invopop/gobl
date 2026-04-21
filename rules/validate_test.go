package rules_test

import (
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
		assert.Equal(t, []string{"$.addr"}, f.Paths())
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
		p := &Person{Name: "Carol", Address: &Address{City: "Anytown"}, SecondAddress: nil}
		faults := rules.Validate(p)
		assert.NoError(t, faults)
	})
}

func TestValidateNil(t *testing.T) {
	faults := rules.Validate(nil)
	assert.Nil(t, faults)
}

func TestValidateWithNilPointerField(t *testing.T) {
	// Nil pointer field should be skipped without panic.
	p := &Person{Name: "Alice", Address: nil}
	// No rules fire since address is nil and name is set.
	assert.NotPanics(t, func() {
		rules.Validate(p)
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
		// Faults with same code+message are merged, so both paths appear in one fault
		assert.True(t, faults.HasPath("$.tags.alpha.addr"))
		assert.True(t, faults.HasPath("$.tags.beta.addr"))
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
		is.InContext(guardTest), nameSet)

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

func TestValidateWithContext(t *testing.T) {
	// Tests the WithContext option path in Validate.
	guardTest := is.In("DE")
	nameSet := rules.For(new(Email),
		rules.Field("addr",
			rules.Assert("01", "email required", is.Present),
		),
	)
	rules.RegisterWithGuard("ctx-opt-test", rules.GOBL.Add("CTXOPT"),
		is.InContext(guardTest), nameSet)

	t.Run("WithContext option injects context", func(t *testing.T) {
		e := &Email{Addr: ""}
		faults := rules.Validate(e, func(rc *rules.Context) {
			rc.Set("country", "DE")
		})
		require.Error(t, faults)
		assert.True(t, faults.HasCode("GOBL-CTXOPT-EMAIL-01"))
	})
}
