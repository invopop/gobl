package convert

import (
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/schema"
	"github.com/stretchr/testify/assert"
)

type stubConverter struct {
	contexts []*Context
}

func (c *stubConverter) Contexts() []*Context                               { return c.contexts }
func (c *stubConverter) Detect(_ *Input) cbc.Key                            { return cbc.KeyEmpty }
func (c *stubConverter) Import(_ cbc.Key, _ []byte) (*gobl.Envelope, error) { return nil, nil }
func (c *stubConverter) Accepts(_ cbc.Key, _ *gobl.Envelope) bool           { return false }
func (c *stubConverter) Export(_ cbc.Key, _ *gobl.Envelope) ([]byte, error) { return nil, nil }

func stub(keys ...cbc.Key) *stubConverter {
	c := new(stubConverter)
	for _, k := range keys {
		c.contexts = append(c.contexts, &Context{Key: k})
	}
	return c
}

func TestRegistryAdd(t *testing.T) {
	t.Run("sorted contexts", func(t *testing.T) {
		r := newRegistry()
		r.add(stub("b", "a"))
		r.add(stub("c"))
		keys := make([]cbc.Key, 0)
		for _, ctx := range r.contexts() {
			keys = append(keys, ctx.Key)
		}
		assert.Equal(t, []cbc.Key{"a", "b", "c"}, keys)
		assert.Len(t, r.converters, 2)
	})
	t.Run("duplicate key", func(t *testing.T) {
		r := newRegistry()
		r.add(stub("a"))
		assert.PanicsWithValue(t, "convert: context a already registered", func() {
			r.add(stub("b", "a"))
		})
		assert.Nil(t, r.contextFor("b"), "nothing added from a rejected converter")
	})
	t.Run("empty key", func(t *testing.T) {
		r := newRegistry()
		assert.PanicsWithValue(t, "convert: context key is empty", func() {
			r.add(stub(""))
		})
	})
	t.Run("no contexts", func(t *testing.T) {
		r := newRegistry()
		assert.PanicsWithValue(t, "convert: converter has no contexts", func() {
			r.add(stub())
		})
	})
}

func TestRegistryContextFor(t *testing.T) {
	r := newRegistry()
	r.add(stub("a"))
	assert.Equal(t, cbc.Key("a"), r.contextFor("a").Key)
	assert.Nil(t, r.contextFor("x"))
}

func TestRegistryConversions(t *testing.T) {
	r := newRegistry()
	inv := schema.GOBL.Add("bill/invoice")
	st := schema.GOBL.Add("bill/status")
	r.add(&stubConverter{contexts: []*Context{
		{Key: "b", Export: []schema.ID{inv}},
		{Key: "a", Import: []schema.ID{inv, st}, Export: []schema.ID{inv}},
	}})
	assert.Equal(t, []*Conversion{
		{Context: "a", Schema: inv, Direction: DirectionImport},
		{Context: "a", Schema: st, Direction: DirectionImport},
		{Context: "a", Schema: inv, Direction: DirectionExport},
		{Context: "b", Schema: inv, Direction: DirectionExport},
	}, r.conversions())
}
