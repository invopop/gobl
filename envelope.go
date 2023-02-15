package gobl

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/uuid"
)

// Envelope wraps around a document adding headers and
// digital signatures. An Envelope is similar to a regular envelope
// in the physical world, it keeps the contents safe and helps
// get the document where its needed.
type Envelope struct {
	// Schema identifies the schema that should be used to understand this document
	Schema schema.ID `json:"$schema" jsonschema:"title=JSON Schema ID"`
	// Details on what the contents are
	Head *Header `json:"head" jsonschema:"title=Header"`
	// The data inside the envelope
	Document *Document `json:"doc" jsonschema:"title=Document"`
	// JSON Web Signatures of the header
	Signatures []*dsig.Signature `json:"sigs,omitempty" jsonschema:"title=Signatures"`
}

// EnvelopeSchema sets the general definition of the schema ID for this version of the
// envelope.
var EnvelopeSchema = schema.GOBL.Add("envelope")

// Calculable defines the methods expected of a document payload that contains a `Calculate`
// method to be used to perform any additional calculations.
type Calculable interface {
	Calculate() error
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

// Envelop is a convenience method that will build a new envelope and insert
// the contents document provided in a single swoop. The resulting envelope
// will still need to be signed afterwards.
func Envelop(doc interface{}) (*Envelope, error) {
	e := NewEnvelope()
	if err := e.Insert(doc); err != nil {
		return nil, err
	}
	return e, nil
}

// Validate ensures that the envelope contains everything it should to be considered valid GoBL.
func (e *Envelope) Validate() error {
	err := validation.ValidateStruct(e,
		validation.Field(&e.Schema, validation.Required),
		validation.Field(&e.Head, validation.Required),
		validation.Field(&e.Document, validation.Required), // this will also check payload
		validation.Field(&e.Signatures, validation.When(e.Head == nil || e.Head.Draft, validation.Empty)),
	)
	if err != nil {
		return err
	}
	return e.verifyDigest()
}

func (e *Envelope) verifyDigest() error {
	d1 := e.Head.Digest
	d2, err := e.Document.Digest()
	if err != nil {
		return err
	}
	if err := d1.Equals(d2); err != nil {
		return fmt.Errorf("document: %w", err)
	}
	return nil
}

// Sign uses the private key to sign the envelope headers. Validation
// is performed automatically, and an error will be raised if still
// in draft mode.
func (e *Envelope) Sign(key *dsig.PrivateKey) error {
	if e.Head == nil || e.Head.Draft {
		return ErrDraft.WithErrorf("cannot sign draft envelope")
	}
	if err := e.Validate(); err != nil {
		return ErrValidation.WithCause(err)
	}
	sig, err := key.Sign(e.Head)
	if err != nil {
		return ErrSignature.WithCause(err)
	}
	e.Signatures = append(e.Signatures, sig)
	return nil
}

// Insert takes the provided document and inserts it into this
// envelope. Calculate will be called automatically.
func (e *Envelope) Insert(doc interface{}) error {
	if e.Head == nil {
		return ErrInternal.WithErrorf("missing head")
	}
	if doc == nil {
		return ErrNoDocument
	}

	if d, ok := doc.(*Document); ok {
		e.Document = d
	} else {
		var err error
		e.Document, err = NewDocument(doc)
		if err != nil {
			return err
		}
	}

	if err := e.calculate(); err != nil {
		return err
	}

	return nil
}

// Calculate is used to perform calculations on the envelope's
// document contents to ensure everything looks correct.
// Headers will be refreshed to ensure they have the latest valid
// digest.
func (e *Envelope) Calculate() error {
	if e.Document == nil {
		return ErrNoDocument
	}
	if e.Document.IsEmpty() {
		return ErrNoDocument
	}

	return e.calculate()
}

func (e *Envelope) calculate() error {
	// Always set our schema version
	e.Schema = EnvelopeSchema

	doc := e.Document.Instance()
	if doc == nil {
		return ErrUnknownSchema.WithErrorf("schema: %v", e.Document.Schema().String())
	}

	// arm doors and cross check
	if obj, ok := doc.(Calculable); ok {
		if err := obj.Calculate(); err != nil {
			return ErrCalculation.WithCause(err)
		}
	}

	// Double check the header looks okay
	if e.Head == nil {
		e.Head = NewHeader()
	}
	if e.Head.UUID.IsZero() {
		e.Head.UUID = uuid.MakeV1()
	}
	var err error
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
