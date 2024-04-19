package head

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

	// Seals of approval from other organisations that can only be added to
	// non-draft envelopes.
	Stamps []*Stamp `json:"stamps,omitempty" jsonschema:"title=Stamps"`

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
	h.UUID = uuid.V1()
	h.Meta = make(cbc.Meta)
	h.Draft = true // default to draft
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
		validation.Field(&h.Stamps,
			validation.When(h.Draft, validation.Empty),
			DetectDuplicateStamps,
		),
	)
}

// AddStamp adds a new stamp to the header. If the stamp already exists,
// it will be overwritten.
func (h *Header) AddStamp(s *Stamp) {
	h.Stamps = AddStamp(h.Stamps, s)
}

// Contains compares the provided header to ensure that all the fields
// and properties are contained within the base header. Only a subset of
// the most important fields are compared.
func (h *Header) Contains(h2 *Header) bool {
	if h.UUID.String() != h2.UUID.String() {
		return false
	}
	if h2.Digest != nil && h.Digest.String() != h2.Digest.String() {
		return false
	}
	for _, s2 := range h2.Stamps {
		match := false
		for _, s := range h.Stamps {
			if s.Provider == s2.Provider && s.Value == s2.Value {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}
	for _, t2 := range h2.Tags {
		match := false
		for _, t := range h.Tags {
			if t == t2 {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}
	for k2, v2 := range h2.Meta {
		v, ok := h.Meta[k2]
		if !ok || v != v2 {
			return false
		}
	}
	if h2.Notes != "" && h2.Notes != h.Notes {
		return false
	}
	return true // all comparisons have passed!
}
