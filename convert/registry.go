package convert

import (
	"fmt"
	"sort"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
)

var converters = newRegistry()

// registry holds the registered converters and the formats they handle.
type registry struct {
	converters []Converter // in registration order
	keys       []cbc.Key   // sorted format keys
	list       map[cbc.Key]*entry
}

type entry struct {
	format    *Format
	converter Converter
	index     int // position of the converter in registration order
}

func newRegistry() *registry {
	return &registry{
		list: make(map[cbc.Key]*entry),
	}
}

func (r *registry) add(c Converter) {
	cs := c.Formats()
	if len(cs) == 0 {
		panic("convert: converter has no formats")
	}
	seen := make(map[cbc.Key]bool, len(cs))
	for _, ctx := range cs {
		if ctx.Key == cbc.KeyEmpty {
			panic("convert: format key is empty")
		}
		if _, ok := r.list[ctx.Key]; ok || seen[ctx.Key] {
			panic(fmt.Sprintf("convert: format %s already registered", ctx.Key))
		}
		seen[ctx.Key] = true
	}
	index := len(r.converters)
	for _, ctx := range cs {
		r.keys = append(r.keys, ctx.Key)
		r.list[ctx.Key] = &entry{format: ctx, converter: c, index: index}
	}
	sort.Slice(r.keys, func(i, j int) bool {
		return r.keys[i].String() < r.keys[j].String()
	})
	r.converters = append(r.converters, c)
}

func (r *registry) entryFor(key cbc.Key) *entry {
	return r.list[key]
}

func (r *registry) formatFor(key cbc.Key) *Format {
	if e := r.list[key]; e != nil {
		return e.format
	}
	return nil
}

func (r *registry) formats() []*Format {
	all := make([]*Format, len(r.keys))
	for i, k := range r.keys {
		all[i] = r.list[k].format
	}
	return all
}

func (r *registry) formatsFor(country l10n.Code) []*Format {
	list := make([]*Format, 0)
	for _, k := range r.keys {
		ctx := r.list[k].format
		if ctx.appliesTo(country) {
			list = append(list, ctx)
		}
	}
	return list
}

func (r *registry) conversions() []*Conversion {
	list := make([]*Conversion, 0)
	for _, k := range r.keys {
		ctx := r.list[k].format
		for _, s := range ctx.Import {
			list = append(list, &Conversion{Format: k, Schema: s, Direction: DirectionImport})
		}
		for _, s := range ctx.Export {
			list = append(list, &Conversion{Format: k, Schema: s, Direction: DirectionExport})
		}
	}
	return list
}

// candidates provides the registration indexes of the converters that
// handle any of the keys, in order, or of all of them if no keys are given.
func (r *registry) candidates(keys []cbc.Key) ([]int, error) {
	match := make([]bool, len(r.converters))
	for _, k := range keys {
		e := r.list[k]
		if e == nil {
			return nil, ErrUnknownFormat.WithReason("format %s not registered", k)
		}
		match[e.index] = true
	}
	list := make([]int, 0, len(r.converters))
	for i := range r.converters {
		if len(keys) == 0 || match[i] {
			list = append(list, i)
		}
	}
	return list, nil
}

func (r *registry) detect(data []byte, keys []cbc.Key) (*entry, error) {
	cs, err := r.candidates(keys)
	if err != nil {
		return nil, err
	}
	in := NewInput(data)
	var match *entry
	for _, i := range cs {
		k := r.converters[i].Detect(in)
		if k == cbc.KeyEmpty {
			continue
		}
		e := r.list[k]
		if e == nil || e.index != i {
			continue // not a format owned by this converter
		}
		if len(keys) > 0 && !k.In(keys...) {
			continue
		}
		if match != nil {
			return nil, ErrAmbiguous.WithReason("detected as %s and %s", match.format.Key, k)
		}
		match = e
	}
	if match == nil {
		return nil, ErrUnknownFormat.WithReason("data not recognized")
	}
	return match, nil
}
