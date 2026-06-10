package tax

import (
	"sort"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// ExternalAddon describes an addon whose implementation lives in a separate,
// optional Go module rather than in GOBL core (for example
// github.com/invopop/gobl.fr.ctc).
//
// Registering an addon here does NOT make it functional at runtime. It only
// records that the key is a recognised, *approved* addon so that:
//
//   - the key is accepted in the JSON Schema `$addons` enum (see
//     [AddonList.JSONSchemaExtend]), and
//   - GOBL keeps a curated, reviewed list of the external addons it endorses.
//
// At Validate/Calculate time the addon must still be actually loaded — its
// module must have been imported so its init() called [RegisterAddonDef] —
// otherwise the "$addons must be registered" rule fails. The approved list is
// recognition and governance, never a runtime bypass: a document is never
// silently processed without the addon's normalizers and rules.
type ExternalAddon struct {
	// Key is the addon key, e.g. "fr-ctc-v1".
	Key cbc.Key
	// Name is a human-readable description, surfaced in the JSON Schema enum.
	Name i18n.String
	// Module is the Go module path that implements the addon, kept as the
	// approval record, e.g. "github.com/invopop/gobl.fr.ctc".
	Module string
}

type externalAddonCollection struct {
	keys []cbc.Key
	list map[cbc.Key]*ExternalAddon
}

var approvedAddons = &externalAddonCollection{list: make(map[cbc.Key]*ExternalAddon)}

func (c *externalAddonCollection) add(ea *ExternalAddon) {
	if _, ok := c.list[ea.Key]; !ok {
		c.keys = append(c.keys, ea.Key)
		sort.Slice(c.keys, func(i, j int) bool {
			return c.keys[i].String() < c.keys[j].String()
		})
	}
	c.list[ea.Key] = ea
}

// RegisterApprovedAddon adds an external addon to the curated list of addons
// recognised by GOBL. See [ExternalAddon] for the semantics. The curated list
// is maintained in the addons package (addons/external.go); adding an entry
// there is the approval step, reviewed via pull request.
func RegisterApprovedAddon(ea *ExternalAddon) {
	approvedAddons.add(ea)
}

// ApprovedAddons returns the curated list of approved external addons, ordered
// by key. Being on this list does not imply the addon is loaded at runtime.
func ApprovedAddons() []*ExternalAddon {
	list := make([]*ExternalAddon, len(approvedAddons.keys))
	for i, k := range approvedAddons.keys {
		list[i] = approvedAddons.list[k]
	}
	return list
}
