package cbc_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/stretchr/testify/assert"
)

func TestDefinitionRegistry(t *testing.T) {
	bar := &cbc.Definition{
		Key:   "bar",
		Name:  i18n.NewString("Bar"),
		Usage: &cbc.Usage{Keys: []cbc.Key{"addon-a"}},
	}
	foo := &cbc.Definition{
		Key:   "foo",
		Name:  i18n.NewString("Foo"),
		Usage: &cbc.Usage{Codes: []cbc.Code{"ES"}},
	}
	gen := &cbc.Definition{
		Key:  "any",
		Name: i18n.NewString("Any"),
	}
	r := cbc.NewDefinitionRegistry()
	r.Add(foo, bar)
	r.Add(gen)

	t.Run("for", func(t *testing.T) {
		assert.Equal(t, foo, r.For("foo"))
		assert.Nil(t, r.For("missing"))
	})
	t.Run("all ordered", func(t *testing.T) {
		assert.Equal(t, []*cbc.Definition{gen, bar, foo}, r.All())
	})
	t.Run("filter", func(t *testing.T) {
		assert.Equal(t, []*cbc.Definition{gen, bar}, r.Filter([]cbc.Key{"addon-a"}, nil))
		assert.Equal(t, []*cbc.Definition{gen, foo}, r.Filter(nil, []cbc.Code{"ES"}))
		assert.Equal(t, []*cbc.Definition{gen}, r.Filter(nil, nil))
	})
	t.Run("same definition twice", func(t *testing.T) {
		assert.NotPanics(t, func() {
			r.Add(foo)
		})
		assert.Len(t, r.All(), 3)
	})
	t.Run("conflicting definition", func(t *testing.T) {
		assert.PanicsWithValue(t, "definition 'foo' already registered", func() {
			r.Add(&cbc.Definition{Key: "foo", Name: i18n.NewString("Other")})
		})
	})
}
