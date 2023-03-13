package gobl

import (
	"context"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
)

// Header defines the metadata of the body. The header is used as the payload
// for the JSON Web Signatures, so we want this to be as compact as possible.
type Header struct {
	// Unique UUIDv1 identifier for the envelope.
	UUID uuid.UUID `json:"uuid" jsonschema:"title=UUID"`

	// Digest of the canonical JSON body.
	Digest *dsig.Digest `json:"dig" jsonschema:"title=Digest"`

	// Seals of approval from other organisations.
	Stamps []*cbc.Stamp `json:"stamps,omitempty" jsonschema:"title=Stamps"`

	// Set of labels that describe but have no influence on the data.
	Tags []string `json:"tags,omitempty" jsonschema:"title=Tags"`

	// Additional semi-structured information about this envelope.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`

	// Any information that may be relevant to other humans about this envelope
	Notes string `json:"notes,omitempty" jsonschema:"title=Notes"`

	// When true, implies that this document should not be considered final. Digital signatures are optional.
	Draft bool `json:"draft,omitempty" jsonschema:"title=Draft"`
}

// NewHeader creates a new header and automatically assigns a UUIDv1.
func NewHeader() *Header {
	h := new(Header)
	h.UUID = uuid.MakeV1()
	h.Meta = make(cbc.Meta)
	return h
}

// Validate checks that the header contains the basic information we need to function.
func (h *Header) Validate() error {
	return h.ValidateWithContext(context.Background())
}

// ValidateWithContext checks that the header contains the basic information we need to function.
func (h *Header) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, h,
		validation.Field(&h.UUID, validation.Required, uuid.IsV1),
		validation.Field(&h.Digest, validation.Required),
		validation.Field(&h.Stamps),
	)
}

// AddStamp adds a new stamp to the header. If the stamp already exists,
// it will be overwritten.
func (h *Header) AddStamp(s *cbc.Stamp) {
	if h.Stamps == nil {
		h.Stamps = make([]*cbc.Stamp, 0)
	}
	for _, v := range h.Stamps {
		if v.Provider == s.Provider {
			v.Value = s.Value
			return
		}
	}
	h.Stamps = append(h.Stamps, s)
}
