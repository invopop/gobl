package gobl

import (
	"bytes"
	"encoding/json"

	"github.com/alecthomas/jsonschema"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/c14n"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/region"
)

// Envelope wraps around a gobl document and provides support for digest creation
// and digital signatures.
type Envelope struct {
	// The GOBL document version used to generate the envelope
	Version Version `json:"ver" jsonschema:"title=Version"`
	// Details on what the contents are
	Head *Header `json:"head" jsonschema:"title=Header"`
	// The data inside the envelope
	Document *Payload `json:"doc" jsonschema:"title=Document,description="`
	// JSON Web Signatures of the header
	Signatures []*dsig.Signature `json:"sigs" jsonschema:"title=Signatures"`
}

// Document defines what we expect from a document to be able to be included in an envelope.
type Document interface {
	Type() string
}

// Calculable defines the methods expected of a document payload that contains a `Calculate`
// method to be used to perform any additional calculations.
type Calculable interface {
	Calculate(r region.Region) error
}

// Validatable describes a document that can be validated.
type Validatable interface {
	Validate(r region.Region) error
}

// NewEnvelope builds a new envelope object ready for data to be inserted
// and signed. If you are loading data from json, you can safely use a regular
// `new(Envelope)` call directly.
//
// A known region code is required as this will be used for any calculations and
// validations that need to be performed on the document to be inserted.
func NewEnvelope(rc region.Code) *Envelope {
	e := new(Envelope)
	e.Version = VERSION
	e.Head = NewHeader(rc)
	e.Document = new(Payload)
	e.Signatures = make([]*dsig.Signature, 0)
	return e
}

// Region extracts the region from the header and provides a complete region object.
func (e *Envelope) Region() region.Region {
	if e.Head == nil || e.Head.Region == "" {
		return nil
	}
	return Regions().For(e.Head.Region)
}

// Validate ensures that the envelope contains everything it should to be considered valid GoBL.
func (e *Envelope) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.Version, validation.Required),
		validation.Field(&e.Head, validation.Required),
		validation.Field(&e.Document, validation.Required),
		validation.Field(&e.Signatures, validation.When(e.Head != nil && !e.Head.Draft, validation.Required)),
	)
}

// Verify ensures the digest headers still match the document contents.
func (e *Envelope) Verify() error {
	d1 := e.Head.Digest
	d2, err := e.Document.digest()
	if err != nil {
		return err
	}
	return d1.Equals(d2)
}

// Sign uses the private key to the envelope headers.
func (e *Envelope) Sign(key *dsig.PrivateKey) error {
	sig, err := key.Sign(e.Head)
	if err != nil {
		return ErrSignature.WithCause(err)
	}
	e.Signatures = append(e.Signatures, sig)
	return nil
}

// Insert takes the provided document, performs any calculations, validates, then
// serializes it ready for use.
func (e *Envelope) Insert(doc Document) error {
	if e.Head == nil {
		return ErrInternal.WithErrorf("missing head")
	}
	if e.Version == "" {
		e.Version = VERSION
	}

	// arm doors and cross check
	r := e.Region()
	if r == nil {
		return ErrNoRegion
	}
	if obj, ok := doc.(Calculable); ok {
		if err := obj.Calculate(r); err != nil {
			return ErrCalculation.WithCause(err)
		}
	}
	if obj, ok := doc.(Validatable); ok {
		if err := obj.Validate(r); err != nil {
			return ErrValidation.WithCause(err)
		}
	}

	if e.Document == nil {
		e.Document = new(Payload)
	}
	if err := e.Document.insert(doc); err != nil {
		return err
	}
	e.Head.Type = doc.Type()

	var err error
	e.Head.Digest, err = e.Document.digest()
	if err != nil {
		return err
	}

	return nil
}

// Extract the contents of the envelope into the provided document type.
func (e *Envelope) Extract(doc Document) error {
	if e.Document == nil {
		return ErrNoDocument.WithErrorf("cannot extract document from empty envelope")
	}
	return e.Document.extract(doc)
}

// Payload helps us handle the document's contents by essentially wrapping around
// the json RawMessage.
type Payload struct {
	data json.RawMessage
}

func (p *Payload) insert(doc Document) error {
	var err error
	p.data, err = json.Marshal(doc)
	if err != nil {
		return ErrMarshal.WithCause(err)
	}
	return nil
}

func (p *Payload) extract(doc Document) error {
	return json.Unmarshal(p.data, doc)
}

func (p *Payload) digest() (*dsig.Digest, error) {
	r := bytes.NewReader(p.data)
	cd, err := c14n.CanonicalJSON(r)
	if err != nil {
		return nil, ErrInternal.WithErrorf("canonical JSON error: %w", err)
	}
	return dsig.NewSHA256Digest(cd), nil
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
func (p *Payload) UnmarshalJSON(data []byte) error {
	p.data = json.RawMessage(data)
	return nil
}

// MarshalJSON satisfies the json.Marshaler interface.
func (p *Payload) MarshalJSON() ([]byte, error) {
	return p.data, nil
}

// JSONSchemaType returns a jsonschema.Type object.
func (Payload) JSONSchemaType() *jsonschema.Type {
	return &jsonschema.Type{
		Type:        "object",
		Title:       "Payload",
		Description: "Contents of the envelope",
	}
}
