package legal

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/uuid"
)

// Assent is a separately signed declaration of intent concerning one exact
// contract artifact. Keeping it separate avoids embedding a signature inside
// the digest it is intended to sign and supports counterparts.
type Assent struct {
	uuid.Identify
	Contract *Reference `json:"contract" jsonschema:"title=Contract"`
	Party    cbc.Key    `json:"party" jsonschema:"title=Party"`
	// Entity identifies the concrete legal person taking an open party role.
	// It may be omitted when the Contract already identifies that entity.
	Entity    *org.Party `json:"entity,omitempty" jsonschema:"title=Entity"`
	Signatory cbc.Key    `json:"signatory,omitempty" jsonschema:"title=Signatory"`
	// Representative identifies a person acting for Entity when the contract
	// did not name that person in advance.
	Representative *org.Person `json:"representative,omitempty" jsonschema:"title=Representative"`
	Intent         cbc.Key     `json:"intent" jsonschema:"title=Intent"`
	Capacity       string      `json:"capacity,omitempty" jsonschema:"title=Capacity"`
	Statement      string      `json:"statement,omitempty" jsonschema:"title=Statement"`
}

// Calculate prepares the assent for use as a GOBL document payload.
func (a *Assent) Calculate() error {
	norm.Normalize(a)
	return nil
}
