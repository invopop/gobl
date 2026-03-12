package head

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/uuid"
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

	// Links provide URLs to other resources that are related to this envelope
	// and unlike stamps can be added even in the draft state.
	Links []*Link `json:"links,omitempty" jsonschema:"title=Links"`

	// Set of labels that describe but have no influence on the data.
	Tags []string `json:"tags,omitempty" jsonschema:"title=Tags"`

	// Additional semi-structured information about this envelope.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`

	// Any information that may be relevant to other humans about this envelope
	Notes string `json:"notes,omitempty" jsonschema:"title=Notes"`
}

// NewHeader creates a new header and automatically assigns a UUIDv1.
func NewHeader() *Header {
	h := new(Header)
	h.UUID = uuid.V7()
	h.Meta = make(cbc.Meta)
	return h
}

func headerRules() *rules.Set {
	return rules.For(new(Header),
		rules.Field("uuid",
			rules.Assert("01", "header must contain a UUID v1 or v7 with timestamp",
				rules.Required,
				uuid.HasTimestamp,
			),
		),
		rules.Field("dig",
			rules.Assert("02", "header must have a digest",
				rules.Required,
			),
		),
		rules.Field("stamps",
			rules.Assert("03", "duplicate stamp providers are not allowed",
				rules.By("no duplicate stamps", hasNoDuplicateStamps),
			),
		),
		rules.Field("links",
			rules.Assert("04", "duplicate link keys are not allowed",
				rules.By("no duplicate links", hasNoDuplicateLinks),
			),
		),
	)
}

func hasNoDuplicateStamps(val any) bool {
	stamps, ok := val.([]*Stamp)
	if !ok {
		return true
	}
	seen := make([]*Stamp, 0, len(stamps))
	for _, s := range stamps {
		if s.In(seen) {
			return false
		}
		seen = append(seen, s)
	}
	return true
}

func hasNoDuplicateLinks(val any) bool {
	links, ok := val.([]*Link)
	if !ok {
		return true
	}
	seen := make([]*Link, 0, len(links))
	for _, l := range links {
		if LinkByCategoryAndKey(seen, l.Category, l.Key) != nil {
			return false
		}
		seen = append(seen, l)
	}
	return true
}

// AddStamp adds a new stamp to the header. If the stamp already exists,
// it will be overwritten.
func (h *Header) AddStamp(s *Stamp) {
	h.Stamps = AddStamp(h.Stamps, s)
}

// Stamp provides the stamp for the given provider or nil.
func (h *Header) Stamp(provider cbc.Key) *Stamp {
	return GetStamp(h.Stamps, provider)
}

// GetStamp provides the stamp for the given provider or nil.
// Deprecated: use Stamp instead.
func (h *Header) GetStamp(provider cbc.Key) *Stamp {
	return h.Stamp(provider)
}

// AddLink will add the link to the header, or update a link with the same
// category and key.
func (h *Header) AddLink(l *Link) {
	h.Links = AppendLink(h.Links, l)
}

// Link provides the link with the matching category and key in the header, or nil.
func (h *Header) Link(category, key cbc.Key) *Link {
	return LinkByCategoryAndKey(h.Links, category, key)
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
	for _, l2 := range h2.Links {
		match := false
		for _, l := range h.Links {
			if l.Key == l2.Key && l.URL == l2.URL {
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
