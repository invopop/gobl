package head

import (
	"errors"
	"time"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/uuid"
)

var (
	// ErrSignaturePayload is returned when a signature's payload cannot be parsed.
	ErrSignaturePayload = errors.New("head: invalid signature payload")
	// ErrSignatureMismatch is returned when a signature's payload does not match the header.
	ErrSignatureMismatch = errors.New("head: signature payload mismatch")
	// ErrSignatureKeyMismatch is returned when no provided key matches the signature.
	ErrSignatureKeyMismatch = errors.New("head: no key match found")
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

	// From is the URI-form transport address of the envelope's issuer,
	// e.g. "gobl:samlown.example.com" or
	// "iso6523-actorid-upis::9920:b123123123".
	From cbc.URI `json:"from,omitempty" jsonschema:"title=From"`

	// To is the URI-form transport address of the envelope's intended
	// receiver.
	To cbc.URI `json:"to,omitempty" jsonschema:"title=To"`

	// Ignore lists fully-qualified validation fault codes to suppress when
	// validating this envelope, e.g. "GOBL-EU-EN16931-ORG-ITEM-01". Intended
	// for format-conversion cases where a specific, known fault is acceptable.
	// Covered by the envelope signature when sealed. NOTE: any code may be
	// listed, including structural envelope/header codes — use deliberately.
	Ignore []rules.Code `json:"ignore,omitempty" jsonschema:"title=Ignore"`
}

// RulesContext injects this header's Ignore codes into the validation context
// so that matching faults are dropped from the result.
func (h *Header) RulesContext() rules.WithContext {
	if h == nil || len(h.Ignore) == 0 {
		return func(*rules.Context) {}
	}
	return rules.WithIgnore(h.Ignore...)
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
				is.Present,
				uuid.HasTimestamp,
			),
		),
		rules.Field("dig",
			rules.Assert("02", "header must have a digest",
				is.Present,
			),
		),
		rules.Field("stamps",
			rules.Assert("03", "duplicate stamp providers are not allowed",
				is.Func("no duplicate stamps", hasNoDuplicateStamps),
			),
		),
		rules.Field("links",
			rules.Assert("04", "duplicate link keys are not allowed",
				is.Func("no duplicate links", hasNoDuplicateLinks),
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

// SigningPayload defines the fields locked by a signature. UUID and
// Digest identify the document; Iss and Aud are the verifiable origin
// and audience of *this* signature (https URLs); IssuedAt is the time
// the signature was produced as a JWT-standard NumericDate (Unix
// seconds, per RFC 7519 §2). Header stamps, links, tags, meta, notes
// and the (unsigned, intent-level) From/To fields can still be
// modified after signing.
type SigningPayload struct {
	UUID     uuid.UUID    `json:"uuid"`
	Digest   *dsig.Digest `json:"dig"`
	Iss      cbc.URI      `json:"iss,omitempty"`
	Aud      cbc.URI      `json:"aud,omitempty"`
	IssuedAt int64        `json:"iat,omitempty"`
}

func (h *Header) payload(iss, aud cbc.URI, iat int64) *SigningPayload {
	return &SigningPayload{
		UUID:     h.UUID,
		Digest:   h.Digest,
		Iss:      iss,
		Aud:      aud,
		IssuedAt: iat,
	}
}

// Sign creates a JWS signature over the header's document identity
// (UUID + Digest) together with the signer's GOBL Net identity (iss),
// the optional audience (aud) it is bound to, and the current UTC
// time as a JWT-standard `iat` claim (Unix seconds). Generic JWT
// verifiers resolve the public keys by fetching
// `<iss>/.well-known/jwks.json` from the HTTPS iss URL — no `jku`
// header is needed.
func (h *Header) Sign(key *dsig.PrivateKey, iss, aud cbc.URI, opts ...dsig.SignerOption) (*dsig.Signature, error) {
	iat := time.Now().UTC().Unix()
	return dsig.NewSignature(key, h.payload(iss, aud, iat), opts...)
}

// Verify checks that the signature covers this header's document
// identity (UUID + Digest). If public keys are provided, the signature
// must also be cryptographically valid against at least one of them,
// and when that key declares a validity window the signed `ts` MUST
// fall within it. The signed iss/aud are part of the statement and are
// not validated here — read them with SignedPayload.
func (h *Header) Verify(sig *dsig.Signature, keys ...*dsig.PublicKey) error {
	if len(keys) == 0 {
		p := new(SigningPayload)
		if err := sig.UnsafePayload(p); err != nil {
			return ErrSignaturePayload
		}
		return h.matchPayload(p)
	}
	for _, k := range keys {
		p := new(SigningPayload)
		if err := sig.VerifyPayload(k, p); err != nil {
			continue
		}
		if err := h.matchPayload(p); err != nil {
			return err
		}
		var iat time.Time
		if p.IssuedAt > 0 {
			iat = time.Unix(p.IssuedAt, 0).UTC()
		}
		if err := k.Allows(iat); err != nil {
			return err
		}
		return nil
	}
	return ErrSignatureKeyMismatch
}

// SignedPayload extracts the (unverified) signed payload from a
// signature, used to read iss before fetching the signer's keys.
func SignedPayload(sig *dsig.Signature) (*SigningPayload, error) {
	p := new(SigningPayload)
	if err := sig.UnsafePayload(p); err != nil {
		return nil, ErrSignaturePayload
	}
	return p, nil
}

func (h *Header) matchPayload(actual *SigningPayload) error {
	if h.UUID.String() != actual.UUID.String() {
		return ErrSignatureMismatch
	}
	if h.Digest == nil || actual.Digest == nil {
		if h.Digest != actual.Digest {
			return ErrSignatureMismatch
		}
		return nil
	}
	if err := h.Digest.Equals(actual.Digest); err != nil {
		return ErrSignatureMismatch
	}
	return nil
}

// Contains compares the provided header to ensure that all the fields
// and properties are contained within the base header. Only a subset of
// the most important fields are compared.
//
// Deprecated: Use Verify with a signature instead.
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
