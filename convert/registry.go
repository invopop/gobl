package convert

import (
	"fmt"
	"sort"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
)

var converters = newRegistry()

// registry holds the registered converters and the contexts they handle.
type registry struct {
	converters []Converter // in registration order
	keys       []cbc.Key   // sorted context keys
	list       map[cbc.Key]*entry
}

type entry struct {
	context   *Context
	converter Converter
}

func newRegistry() *registry {
	return &registry{
		list: make(map[cbc.Key]*entry),
	}
}

func (r *registry) add(c Converter) {
	cs := c.Contexts()
	if len(cs) == 0 {
		panic("convert: converter has no contexts")
	}
	for _, ctx := range cs {
		if ctx.Key == cbc.KeyEmpty {
			panic("convert: context key is empty")
		}
		if _, ok := r.list[ctx.Key]; ok {
			panic(fmt.Sprintf("convert: context %s already registered", ctx.Key))
		}
	}
	for _, ctx := range cs {
		r.keys = append(r.keys, ctx.Key)
		r.list[ctx.Key] = &entry{context: ctx, converter: c}
	}
	sort.Slice(r.keys, func(i, j int) bool {
		return r.keys[i].String() < r.keys[j].String()
	})
	r.converters = append(r.converters, c)
}

func (r *registry) entryFor(key cbc.Key) *entry {
	return r.list[key]
}

func (r *registry) contextFor(key cbc.Key) *Context {
	if e := r.list[key]; e != nil {
		return e.context
	}
	return nil
}

func (r *registry) contexts() []*Context {
	all := make([]*Context, len(r.keys))
	for i, k := range r.keys {
		all[i] = r.list[k].context
	}
	return all
}

func (r *registry) contextsFor(country l10n.Code) []*Context {
	list := make([]*Context, 0)
	for _, k := range r.keys {
		ctx := r.list[k].context
		if ctx.appliesTo(country) {
			list = append(list, ctx)
		}
	}
	return list
}

func (r *registry) conversions() []*Conversion {
	list := make([]*Conversion, 0)
	for _, k := range r.keys {
		ctx := r.list[k].context
		for _, s := range ctx.Import {
			list = append(list, &Conversion{Context: k, Schema: s, Direction: DirectionImport})
		}
		for _, s := range ctx.Export {
			list = append(list, &Conversion{Context: k, Schema: s, Direction: DirectionExport})
		}
	}
	return list
}

// candidates provides the converters that handle any of the keys, in
// registration order, or all of them if no keys are given.
func (r *registry) candidates(keys []cbc.Key) ([]Converter, error) {
	if len(keys) == 0 {
		return r.converters, nil
	}
	for _, k := range keys {
		if r.list[k] == nil {
			return nil, ErrUnknownContext.WithReason("context %s not registered", k)
		}
	}
	list := make([]Converter, 0, len(keys))
	for _, c := range r.converters {
		for _, k := range keys {
			if r.list[k].converter == c {
				list = append(list, c)
				break
			}
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
	for _, c := range cs {
		k := c.Detect(in)
		if k == cbc.KeyEmpty {
			continue
		}
		e := r.list[k]
		if e == nil || e.converter != c {
			continue // not a context owned by this converter
		}
		if len(keys) > 0 && !k.In(keys...) {
			continue
		}
		if match != nil {
			return nil, ErrAmbiguous.WithReason("detected as %s and %s", match.context.Key, k)
		}
		match = e
	}
	if match == nil {
		return nil, ErrUnknownContext.WithReason("data not recognized")
	}
	return match, nil
}
