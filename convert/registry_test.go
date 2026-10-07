package convert

import (
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/schema"
	"github.com/stretchr/testify/assert"
)

type stubConverter struct {
	formats []*Format
}

func (c *stubConverter) Formats() []*Format                                 { return c.formats }
func (c *stubConverter) Detect(_ *Input) cbc.Key                            { return cbc.KeyEmpty }
func (c *stubConverter) Import(_ cbc.Key, _ []byte) (*gobl.Envelope, error) { return nil, nil }
func (c *stubConverter) Accepts(_ cbc.Key, _ *gobl.Envelope) bool           { return false }
func (c *stubConverter) Export(_ cbc.Key, _ *gobl.Envelope) ([]byte, error) { return nil, nil }

func stub(keys ...cbc.Key) *stubConverter {
	c := new(stubConverter)
	for _, k := range keys {
		c.formats = append(c.formats, &Format{Key: k})
	}
	return c
}

func TestRegistryAdd(t *testing.T) {
	t.Run("sorted formats", func(t *testing.T) {
		r := newRegistry()
		r.add(stub("b", "a"))
		r.add(stub("c"))
		keys := make([]cbc.Key, 0)
		for _, ctx := range r.formats() {
			keys = append(keys, ctx.Key)
		}
		assert.Equal(t, []cbc.Key{"a", "b", "c"}, keys)
		assert.Len(t, r.converters, 2)
	})
	t.Run("duplicate key", func(t *testing.T) {
		r := newRegistry()
		r.add(stub("a"))
		assert.PanicsWithValue(t, "convert: format a already registered", func() {
			r.add(stub("b", "a"))
		})
		assert.Nil(t, r.formatFor("b"), "nothing added from a rejected converter")
	})
	t.Run("duplicate key in one converter", func(t *testing.T) {
		r := newRegistry()
		assert.PanicsWithValue(t, "convert: format a already registered", func() {
			r.add(stub("a", "a"))
		})
		assert.Empty(t, r.formats())
	})
	t.Run("empty key", func(t *testing.T) {
		r := newRegistry()
		assert.PanicsWithValue(t, "convert: format key is empty", func() {
			r.add(stub(""))
		})
	})
	t.Run("no formats", func(t *testing.T) {
		r := newRegistry()
		assert.PanicsWithValue(t, "convert: converter has no formats", func() {
			r.add(stub())
		})
	})
}

// valueConverter is registered by value and holds a slice, so it is not
// comparable.
type valueConverter struct {
	formats []*Format
	detect  cbc.Key
}

func (c valueConverter) Formats() []*Format                                 { return c.formats }
func (c valueConverter) Detect(_ *Input) cbc.Key                            { return c.detect }
func (c valueConverter) Import(_ cbc.Key, _ []byte) (*gobl.Envelope, error) { return nil, nil }
func (c valueConverter) Accepts(_ cbc.Key, _ *gobl.Envelope) bool           { return false }
func (c valueConverter) Export(_ cbc.Key, _ *gobl.Envelope) ([]byte, error) { return nil, nil }

func TestRegistryDetectValueConverters(t *testing.T) {
	r := newRegistry()
	r.add(valueConverter{formats: []*Format{{Key: "a"}}, detect: "a"})
	r.add(valueConverter{formats: []*Format{{Key: "b"}}, detect: "a"})
	e, err := r.detect(nil, nil)
	assert.NoError(t, err, "format owned by another converter is ignored")
	assert.Equal(t, cbc.Key("a"), e.format.Key)
	e, err = r.detect(nil, []cbc.Key{"a"})
	assert.NoError(t, err)
	assert.Equal(t, cbc.Key("a"), e.format.Key)
}

func TestRegistryFormatFor(t *testing.T) {
	r := newRegistry()
	r.add(stub("a"))
	assert.Equal(t, cbc.Key("a"), r.formatFor("a").Key)
	assert.Nil(t, r.formatFor("x"))
}

func TestRegistryConversions(t *testing.T) {
	r := newRegistry()
	inv := schema.GOBL.Add("bill/invoice")
	st := schema.GOBL.Add("bill/status")
	r.add(&stubConverter{formats: []*Format{
		{Key: "b", Export: []schema.ID{inv}},
		{Key: "a", Import: []schema.ID{inv, st}, Export: []schema.ID{inv}},
	}})
	assert.Equal(t, []*Conversion{
		{Format: "a", Schema: inv, Direction: DirectionImport},
		{Format: "a", Schema: st, Direction: DirectionImport},
		{Format: "a", Schema: inv, Direction: DirectionExport},
		{Format: "b", Schema: inv, Direction: DirectionExport},
	}, r.conversions())
}
