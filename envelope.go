package gobl

import (
	"context"
	"errors"
	"fmt"

	"github.com/invopop/validation"

	"github.com/invopop/gobl/base"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/internal"
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
	Head *base.Header `json:"head" jsonschema:"title=Header"`
	// The data inside the envelope
	Document *base.Document `json:"doc" jsonschema:"title=Document"`
	// JSON Web Signatures of the header
	Signatures []*dsig.Signature `json:"sigs,omitempty" jsonschema:"title=Signatures"`
}

// EnvelopeSchema sets the general definition of the schema ID for this version of the
// envelope.
var EnvelopeSchema = schema.GOBL.Add("envelope")

// NewEnvelope builds a new envelope object ready for data to be inserted
// and signed. If you are loading data from json, you can safely use a regular
// `new(Envelope)` call directly.
func NewEnvelope() *Envelope {
	e := new(Envelope)
	e.Schema = EnvelopeSchema
	e.Head = base.NewHeader()
	e.Document = new(base.Document)
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
	return e.ValidateWithContext(context.Background())
}

// ValidateWithContext ensures that the envelope contains everything it should to be considered valid GoBL.
func (e *Envelope) ValidateWithContext(ctx context.Context) error {
	ctx = context.WithValue(ctx, internal.KeyDraft, e.Head != nil && e.Head.Draft)
	err := validation.ValidateStructWithContext(ctx, e,
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

// Sign uses the private key to sign the envelope headers. The header
// draft flag will be set to false and validation is performed so that
// only valid non-draft documents will be signed.
func (e *Envelope) Sign(key *dsig.PrivateKey) error {
	if e.Head == nil {
		return base.ErrValidation.WithCause(errors.New("missing header"))
	}
	e.Head.Draft = false
	if err := e.Validate(); err != nil {
		return base.ErrValidation.WithCause(err)
	}
	sig, err := key.Sign(e.Head)
	if err != nil {
		return base.ErrSignature.WithCause(err)
	}
	e.Signatures = append(e.Signatures, sig)
	return nil
}

// Insert takes the provided document and inserts it into this
// envelope. Calculate will be called automatically.
func (e *Envelope) Insert(doc interface{}) error {
	if e.Head == nil {
		return base.ErrInternal.WithErrorf("missing head")
	}
	if doc == nil {
		return base.ErrNoDocument
	}

	if d, ok := doc.(*base.Document); ok {
		e.Document = d
	} else {
		var err error
		e.Document, err = base.NewDocument(doc)
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
		return base.ErrNoDocument
	}
	if e.Document.IsEmpty() {
		return base.ErrNoDocument
	}

	return e.calculate()
}

func (e *Envelope) calculate() error {
	// Always set our schema version
	e.Schema = EnvelopeSchema

	// arm doors and cross check
	if err := e.Document.Calculate(); err != nil {
		return base.ErrCalculation.WithCause(err)
	}

	// Double check the header looks okay
	if e.Head == nil {
		e.Head = base.NewHeader()
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

// Correct will attempt to build a new envelope as a correction of the
// current envelope contents, if possible.
func (e *Envelope) Correct(opts ...base.Option) (*Envelope, error) {
	if e.Head != nil && len(e.Head.Stamps) > 0 {
		opts = append(opts, base.WithHead(e.Head))
	}

	nd, err := e.Document.Clone()
	if err != nil {
		return nil, err
	}
	if err := nd.Correct(opts...); err != nil {
		return nil, err
	}

	// Create a completely new envelope with a new set of data.
	return Envelop(nd)
}
