package legal

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
)

// Party identifies a legal person and the role it has in an agreement.
// Key is the stable handle used by execution requirements and annotations.
type Party struct {
	Key cbc.Key `json:"key" jsonschema:"title=Key"`
	// Role is a machine-readable description such as buyer, seller, employer,
	// licensor, or processor. Contract-specific roles are permitted.
	Role cbc.Key `json:"role" jsonschema:"title=Role"`
	// DefinedAs is the human label used for the party in the contract text.
	DefinedAs string `json:"defined_as" jsonschema:"title=Defined As"`
	// Description identifies an open class of parties in standard terms or an
	// offer before a specific legal person has adhered to it.
	Description string `json:"description,omitempty" jsonschema:"title=Description"`
	// Entity contains the party's legal identity and contact details when the
	// party is already known. An open role may instead use Description and be
	// bound to a concrete entity by a later Assent.
	Entity *org.Party `json:"entity,omitempty" jsonschema:"title=Entity"`
	// Signatories are people authorised or expected to act for this party.
	Signatories []*Signatory `json:"signatories,omitempty" jsonschema:"title=Signatories"`
}

// Signatory describes a human who may provide assent on behalf of a party.
// It records the asserted capacity; proof of authority remains separate
// evidence and may be linked using Authority.
type Signatory struct {
	Key       cbc.Key     `json:"key" jsonschema:"title=Key"`
	Person    *org.Person `json:"person" jsonschema:"title=Person"`
	Capacity  string      `json:"capacity,omitempty" jsonschema:"title=Capacity"`
	Authority *Reference  `json:"authority,omitempty" jsonschema:"title=Authority Evidence"`
}
