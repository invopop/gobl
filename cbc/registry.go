package cbc

import (
	"fmt"
	"slices"
)

// DefinitionRegistry holds a set of definitions identified by their keys,
// such as the extensions or stamps offered by regimes and addons. Definitions
// are expected to be added from package init functions.
type DefinitionRegistry struct {
	keys []Key // ordered list
	list map[Key]*Definition
}

// NewDefinitionRegistry prepares an empty definition registry.
func NewDefinitionRegistry() *DefinitionRegistry {
	return &DefinitionRegistry{
		list: make(map[Key]*Definition),
	}
}

// Add registers the definitions. Adding the same definition twice has no
// effect, while adding a different definition with an existing key panics.
func (r *DefinitionRegistry) Add(defs ...*Definition) {
	for _, d := range defs {
		if ed, ok := r.list[d.Key]; ok {
			if ed != d {
				panic(fmt.Sprintf("definition '%s' already registered", d.Key))
			}
			continue
		}
		i, _ := slices.BinarySearch(r.keys, d.Key)
		r.keys = slices.Insert(r.keys, i, d.Key)
		r.list[d.Key] = d
	}
}

// For provides the definition with the given key, or nil.
func (r *DefinitionRegistry) For(key Key) *Definition {
	return r.list[key]
}

// All provides the registered definitions ordered by key.
func (r *DefinitionRegistry) All() []*Definition {
	all := make([]*Definition, len(r.keys))
	for i, k := range r.keys {
		all[i] = r.list[k]
	}
	return all
}

// Filter provides the registered definitions, ordered by key, whose usage
// matches the provided keys and codes.
func (r *DefinitionRegistry) Filter(keys []Key, codes []Code) []*Definition {
	list := make([]*Definition, 0)
	for _, k := range r.keys {
		if d := r.list[k]; d.Usage.Matches(keys, codes) {
			list = append(list, d)
		}
	}
	return list
}
