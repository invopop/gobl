package gobl

import (
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/org"
)

// Header defines the meta data of the body. The header is used as the payload
// for the JSON Web Signatures, so we want this to be as compact as possible.
type Header struct {
	Type   string       `json:"typ" jsonschema:"title=Type,description=Body type of the document contents"`
	Digest *dsig.Digest `json:"dig" jsonschema:"title=Digest,description=Digest of the canonical JSON body"`
	Stamps []*Stamp     `json:"stamps,omitempty"`
	Meta   org.Meta     `json:"meta,omitempty"`
}

// Stamp defines an official seal of approval from a third party like a governmental agency
// or intermediary and should thus be included in any official envelopes.
type Stamp struct {
	Provider string `json:"prv" jsonschema:"title=Provider,description=Identity of the agency used to create the stamp"`
	Value    string `json:"val" jsonschema:"title=Value,description=The serialized stamp value generated for or by the external agency"`
}
