package head

import (
	"fmt"
	"sort"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
)

// StampDef describes a stamp provider key that may be used in a Stamp,
// alongside the regime or addon that defines it. Documents with a matching
// `$regime` or `$addons` entry can be expected to receive the stamp.
type StampDef struct {
	// Regime code that defines the stamp, matching a document's `$regime`.
	Regime l10n.TaxCountryCode `json:"regime,omitempty" jsonschema:"title=Regime"`
	// Addon key that defines the stamp, matching an entry in a document's `$addons`.
	Addon cbc.Key `json:"addon,omitempty" jsonschema:"title=Addon"`
	*cbc.Definition
}

func (sd *StampDef) owner() string {
	if sd.Addon != cbc.KeyEmpty {
		return "addon '" + sd.Addon.String() + "'"
	}
	return "regime '" + sd.Regime.String() + "'"
}

// stampDefs is a global register of all stamp definitions that have been
// registered by regimes and addons.
var stampDefs = newStampDefCollection()

type stampDefCollection struct {
	keys []cbc.Key
	list map[cbc.Key]*StampDef
}

func newStampDefCollection() *stampDefCollection {
	return &stampDefCollection{
		list: make(map[cbc.Key]*StampDef),
	}
}

func (c *stampDefCollection) add(sd *StampDef) {
	if prev, ok := c.list[sd.Key]; ok {
		if prev.Regime != sd.Regime || prev.Addon != sd.Addon {
			panic(fmt.Sprintf("stamp '%s' from %s already registered by %s", sd.Key, sd.owner(), prev.owner()))
		}
	} else {
		c.keys = append(c.keys, sd.Key)
		sort.Slice(c.keys, func(i, j int) bool {
			return c.keys[i].String() < c.keys[j].String()
		})
	}
	c.list[sd.Key] = sd
}

// RegisterRegimeStamps adds the stamp definitions to the global register
// as belonging to the regime with the given code. This is expected to be
// called from the regime package's init function. Each definition's key is
// the value expected in a Stamp's Provider field. A stamp key may only belong
// to a single regime or addon: registering it from another will panic.
func RegisterRegimeStamps(regime l10n.TaxCountryCode, defs ...*cbc.Definition) {
	for _, d := range defs {
		stampDefs.add(&StampDef{Regime: regime, Definition: d})
	}
}

// RegisterAddonStamps adds the stamp definitions to the global register as
// belonging to the addon with the given key. This is expected to be called
// from the addon package's init function. A stamp key may only belong to a
// single regime or addon: registering it from another will panic.
func RegisterAddonStamps(addon cbc.Key, defs ...*cbc.Definition) {
	for _, d := range defs {
		stampDefs.add(&StampDef{Addon: addon, Definition: d})
	}
}

// StampDefForKey returns the stamp definition for the given provider key,
// or nil if none was registered.
func StampDefForKey(key cbc.Key) *StampDef {
	return stampDefs.list[key]
}

// AllStampDefs provides all the registered stamp definitions ordered by key.
func AllStampDefs() []*StampDef {
	all := make([]*StampDef, len(stampDefs.keys))
	for i, k := range stampDefs.keys {
		all[i] = stampDefs.list[k]
	}
	return all
}
