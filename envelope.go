package gobl

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/schema"
)

// Envelope wraps around a gobl document and provides support for digest creation
// and digital signatures.
type Envelope struct {
	// Schema identifies the schema that should be used to understand this document
	Schema schema.ID `json:"$schema" jsonschema:"title=JSON Schema ID"`
	// Details on what the contents are
	Head *Header `json:"head" jsonschema:"title=Header"`
	// The data inside the envelope
	Document *Document `json:"doc" jsonschema:"title=Document"`
	// JSON Web Signatures of the header
	Signatures []*dsig.Signature `json:"sigs" jsonschema:"title=Signatures"`
}

// EnvelopeSchema sets the general definition of the schema ID for this version of the
// envelope.
var EnvelopeSchema = schema.GOBL.Add("envelope")

// Calculable defines the methods expected of a document payload that contains a `Calculate`
// method to be used to perform any additional calculations.
type Calculable interface {
	Calculate() error
}

// Validatable describes a document that can be validated.
type Validatable interface {
	Validate() error
}

// NewEnvelope builds a new envelope object ready for data to be inserted
// and signed. If you are loading data from json, you can safely use a regular
// `new(Envelope)` call directly.
func NewEnvelope() *Envelope {
	e := new(Envelope)
	e.Schema = EnvelopeSchema
	e.Head = NewHeader()
	e.Document = new(Document)
	e.Signatures = make([]*dsig.Signature, 0)
	return e
}

// Validate ensures that the envelope contains everything it should to be considered valid GoBL.
func (e *Envelope) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.Schema, validation.Required),
		validation.Field(&e.Head, validation.Required),
		validation.Field(&e.Document, validation.Required),
		validation.Field(&e.Signatures, validation.When(e.Head != nil && !e.Head.Draft, validation.Required)),
	)
}

// Verify ensures the digest headers still match the document contents.
func (e *Envelope) Verify() error {
	d1 := e.Head.Digest
	d2, err := e.Document.Digest()
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
func (e *Envelope) Insert(doc interface{}) error {
	if e.Head == nil {
		return ErrInternal.WithErrorf("missing head")
	}

	e.Document = new(Document)
	if err := e.Document.insert(doc); err != nil {
		return err
	}

	if err := e.complete(doc); err != nil {
		return err
	}

	return nil
}

// Complete is used to perform calculations and validations on the envelopes document
// if it responds to the Calculable and Validatable interfaces. Behind the scenes,
// this method will determine the document type, extract, calculate, validate,
// and then re-insert the potentially updated contents.
func (e *Envelope) Complete() error {
	if e.Document == nil {
		return ErrNoDocument
	}

	obj := e.Document.Instance()
	if obj == nil {
		return ErrUnknownSchema.WithErrorf("schema: %v", e.Document.Schema().String())
	}
	if err := e.complete(obj); err != nil {
		return err
	}

	return nil
}

func (e *Envelope) complete(doc interface{}) error {
	// arm doors and cross check
	if obj, ok := doc.(Calculable); ok {
		if err := obj.Calculate(); err != nil {
			return ErrCalculation.WithCause(err)
		}
	}
	if obj, ok := doc.(Validatable); ok {
		if err := obj.Validate(); err != nil {
			return ErrValidation.WithCause(err)
		}
	}

	var err error
	if e.Head == nil {
		e.Head = NewHeader()
	}
	e.Head.Digest, err = e.Document.Digest()
	if err != nil {
		return err
	}

	return nil
}

// Extract the contents of the envelope into the provided document type.
func (e *Envelope) Extract() interface{} {
	if e.Document == nil {
		return nil
	}
	return e.Document.Instance()
}
