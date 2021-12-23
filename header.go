package gobl

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/region"
	"github.com/invopop/gobl/uuid"
)

// Header defines the meta data of the body. The header is used as the payload
// for the JSON Web Signatures, so we want this to be as compact as possible.
type Header struct {
	// Unique UUIDv1 identifier for the envelope.
	UUID uuid.UUID `json:"uuid" jsonschema:"title=UUID"`

	// Body type of the document contents.
	Type string `json:"typ" jsonschema:"title=Type"`

	// Code for the region the document should be validated with.
	Region region.Code `json:"rgn" jsonschema:"title=Region"`

	// Digest of the canonical JSON body.
	Digest *dsig.Digest `json:"dig" jsonschema:"title=Digest"`

	// Seals of approval from other organisations.
	Stamps []*Stamp `json:"stamps,omitempty" jsonschema:"title=Stamps"`

	// Set of labels that describe but have no influence on the data.
	Tags []string `json:"tags,omitempty" jsonschema:"title=Tags"`

	// Additional semi-structured information about this envelope.
	Meta org.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`

	// Any information that may be relevant to other humans about this envelope
	Notes string `json:"notes,omitempty" jsonschema:"title=Notes"`

	// When true, implies that this document should not be considered final. Digital signatures are optional.
	Draft bool `json:"draft,omitempty" jsonschema:"title=Draft"`
}

// NewHeader creates a new header and automatically assigns a UUIDv1.
func NewHeader(rc region.Code) *Header {
	h := new(Header)
	h.UUID = uuid.NewV1()
	h.Region = rc
	h.Meta = make(org.Meta)
	return h
}

// Stamp defines an official seal of approval from a third party like a governmental agency
// or intermediary and should thus be included in any official envelopes.
type Stamp struct {
	Provider string `json:"prv" jsonschema:"title=Provider,description=Identity of the agency used to create the stamp"`
	Value    string `json:"val" jsonschema:"title=Value,description=The serialized stamp value generated for or by the external agency"`
}

// Validate checks that the header contains the basic information we need to function.
func (h *Header) Validate() error {
	return validation.ValidateStruct(h,
		validation.Field(&h.UUID, validation.Required, uuid.IsV1),
		validation.Field(&h.Type, validation.Required),
		validation.Field(&h.Region, validation.Required),
		validation.Field(&h.Digest, validation.Required),
	)
}
