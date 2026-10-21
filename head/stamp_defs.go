package head

import (
	"fmt"
	"sort"

	"github.com/invopop/gobl/cbc"
)

// stampDefs is a global register of all stamp definitions.
var stampDefs = newStampDefCollection()

type stampDefCollection struct {
	keys []cbc.Key
	list map[cbc.Key]*cbc.Definition
}

func newStampDefCollection() *stampDefCollection {
	return &stampDefCollection{
		list: make(map[cbc.Key]*cbc.Definition),
	}
}

func (c *stampDefCollection) add(sd *cbc.Definition) {
	if prev, ok := c.list[sd.Key]; ok {
		if prev != sd {
			panic(fmt.Sprintf("stamp '%s' already registered", sd.Key))
		}
		return
	}
	c.keys = append(c.keys, sd.Key)
	sort.Slice(c.keys, func(i, j int) bool {
		return c.keys[i].String() < c.keys[j].String()
	})
	c.list[sd.Key] = sd
}

// RegisterStamps adds stamp definitions to the global register. Each
// definition's key is the value expected in a Stamp's Provider field.
// Registering a different definition for a key that already exists will
// panic.
func RegisterStamps(defs ...*cbc.Definition) {
	for _, sd := range defs {
		stampDefs.add(sd)
	}
}

// StampDefForKey returns the stamp definition for the given provider key,
// or nil if none was registered.
func StampDefForKey(key cbc.Key) *cbc.Definition {
	return stampDefs.list[key]
}

// AllStampDefs provides all the registered stamp definitions ordered by key.
func AllStampDefs() []*cbc.Definition {
	all := make([]*cbc.Definition, len(stampDefs.keys))
	for i, k := range stampDefs.keys {
		all[i] = stampDefs.list[k]
	}
	return all
}
