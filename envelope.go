package gobl

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/invopop/validation"

	"github.com/invopop/gobl/c14n"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/head"
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
	Head *head.Header `json:"head" jsonschema:"title=Header"`
	// The data inside the envelope
	Document *schema.Object `json:"doc" jsonschema:"title=Document"`
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
	e.Head = head.NewHeader()
	e.Document = new(schema.Object)
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

// Verify checks the envelope's signatures to ensure the headers they contain
// still matches with the current headers. If a list of public keys are provided,
// they will be used to ensure that the signatures we're signed by at least
// one of them. If no keys are provided, only the contents will be checked.
func (e *Envelope) Verify(keys ...*dsig.PublicKey) error {
	if e.Head == nil || e.Head.Draft {
		return errors.New("cannot verify draft document")
	}
	if len(e.Signatures) == 0 {
		return errors.New("no signatures to verify")
	}

	ve := make(validation.Errors)
	for i, s := range e.Signatures {
		if err := e.verifySignature(s, keys...); err != nil {
			ve[strconv.Itoa(i)] = err
		}
	}
	if len(ve) > 0 {
		return ErrValidation.WithCause(validation.Errors{
			"signatures": ve,
		})
	}

	return nil
}

// VerifySignature checks a specific signature with the envelope to see if its
// contents are still valid.
// If a list of public keys are provided, they will be used to ensure that the
// signature was signed by at least one of them. If no keys are provided, only
// the contents will be checked.
func (e *Envelope) VerifySignature(sig *dsig.Signature, keys ...*dsig.PublicKey) error {
	return wrapError(e.verifySignature(sig, keys...))
}

func (e *Envelope) verifySignature(sig *dsig.Signature, keys ...*dsig.PublicKey) error {
	if len(keys) == 0 {
		// no keys provided, only check the contents
		h := new(head.Header)
		if err := sig.UnsafePayload(h); err != nil {
			return errors.New("invalid signature payload")
		}
		if !e.Head.Contains(h) {
			return errors.New("header mismatch")
		}
		return nil
	}
	for _, k := range keys {
		h := new(head.Header)
		if err := sig.VerifyPayload(k, h); err != nil {
			continue
		}
		if e.Head.Contains(h) {
			return nil
		}
		return errors.New("header mismatch")
	}
	return errors.New("no key match found")
}

// ValidateWithContext ensures that the envelope contains everything it should to be considered valid GoBL.
func (e *Envelope) ValidateWithContext(ctx context.Context) error {
	ctx = context.WithValue(ctx, internal.KeyDraft, e.Head != nil && e.Head.Draft)
	err := validation.ValidateStructWithContext(ctx, e,
		validation.Field(&e.Schema, validation.Required),
		validation.Field(&e.Head, validation.Required),
		validation.Field(&e.Document, validation.Required), // this will also check payload
		validation.Field(&e.Signatures,
			validation.When(
				e.Head == nil || e.Head.Draft,
				validation.Empty,
			),
		),
	)
	if err != nil {
		return wrapError(err)
	}
	return wrapError(e.verifyDigest())
}

func (e *Envelope) verifyDigest() error {
	d1 := e.Head.Digest
	d2, err := e.Digest()
	if err != nil {
		return err
	}
	if err := d1.Equals(d2); err != nil {
		return ErrDigest.WithCause(err)
	}
	return nil
}

// Sign uses the private key to sign the envelope headers. The header
// draft flag will be set to false and validation is performed so that
// only valid non-draft documents will be signed.
func (e *Envelope) Sign(key *dsig.PrivateKey) error {
	if e.Head == nil {
		return ErrValidation.WithReason("header required")
	}
	e.Head.Draft = false
	if err := e.Validate(); err != nil {
		return err
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
		return ErrInternal.WithReason("missing head")
	}
	if doc == nil {
		return ErrNoDocument
	}

	if d, ok := doc.(*schema.Object); ok {
		e.Document = d
	} else {
		var err error
		e.Document, err = schema.NewObject(doc)
		if err != nil {
			return wrapError(err)
		}
	}

	if err := e.calculate(); err != nil {
		return wrapError(err)
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

	// arm doors and cross check
	if err := e.Document.Calculate(); err != nil {
		return ErrCalculation.WithCause(err)
	}

	// Double check the header looks okay
	if e.Head == nil {
		e.Head = head.NewHeader()
	}
	if e.Head.UUID.IsZero() {
		e.Head.UUID = uuid.MakeV1()
	}
	var err error
	e.Head.Digest, err = e.Digest()
	if err != nil {
		return err
	}

	return nil
}

// Digest calculates a digital digest using the canonical JSON of the document.
func (e *Envelope) Digest() (*dsig.Digest, error) {
	data, err := json.Marshal(e.Document)
	if err != nil {
		return nil, ErrMarshal.WithCause(err)
	}
	r := bytes.NewReader(data)
	cd, err := c14n.CanonicalJSON(r)
	if err != nil {
		return nil, ErrInternal.WithReason("canonical JSON error: %s", err.Error())
	}
	return dsig.NewSHA256Digest(cd), nil
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
func (e *Envelope) Correct(opts ...schema.Option) (*Envelope, error) {
	if e.Head != nil && len(e.Head.Stamps) > 0 {
		opts = append(opts, head.WithHead(e.Head))
	}

	nd, err := e.Document.Clone()
	if err != nil {
		return nil, wrapError(err)
	}
	if err := nd.Correct(opts...); err != nil {
		return nil, wrapError(err)
	}

	// Create a completely new envelope with a new set of data.
	return Envelop(nd)
}

// CorrectionOptionsSchema will attempt to provide a corrective options JSON Schema
// that can be used to generate a JSON object to send when correcting a document.
// If none are available, the result will be nil.
func (e *Envelope) CorrectionOptionsSchema() (interface{}, error) {
	if e.Document == nil {
		return nil, ErrNoDocument
	}
	opts, err := e.Document.CorrectionOptionsSchema()
	if err != nil {
		return nil, wrapError(err)
	}
	return opts, nil
}
