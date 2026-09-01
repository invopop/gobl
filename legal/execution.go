package legal

import "github.com/invopop/gobl/cbc"

// Execution intent keys describe the legal meaning attributed to an assent.
const (
	IntentExecute     cbc.Key = "execute"
	IntentApprove     cbc.Key = "approve"
	IntentAcknowledge cbc.Key = "acknowledge"
	IntentWitness     cbc.Key = "witness"
)

// Execution describes what assent is required for a contract version. Actual
// assent is represented by separately signed Assent documents.
type Execution struct {
	Method       cbc.Key                 `json:"method,omitempty" jsonschema:"title=Method"`
	Counterparts bool                    `json:"counterparts,omitempty" jsonschema:"title=Counterparts"`
	Statement    string                  `json:"statement,omitempty" jsonschema:"title=Execution Statement"`
	Requirements []*SignatureRequirement `json:"requirements" jsonschema:"title=Signature Requirements"`
}

// SignatureRequirement identifies whose assent is required and in what role.
type SignatureRequirement struct {
	Party     cbc.Key `json:"party" jsonschema:"title=Party"`
	Signatory cbc.Key `json:"signatory,omitempty" jsonschema:"title=Signatory"`
	Intent    cbc.Key `json:"intent" jsonschema:"title=Intent"`
	Witnesses int     `json:"witnesses,omitempty" jsonschema:"title=Witnesses"`
}
