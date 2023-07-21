package gobl

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/invopop/gobl/c14n"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/internal"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Document helps us handle the document's contents by essentially wrapping around
// the contents and ensuring that a `$schema` property is added automatically when
// marshalling into JSON.
type Document struct {
	schema  schema.ID
	payload interface{}
}

// Normalizable defines the methods expected of a document payload that contains a `Normalize`
// method to be used to perform any additional data normalization and calculations.
type Normalizable interface {
	Normalize() error
}

// Correctable defines the expected interface of a document that can be
// corrected.
type Correctable interface {
	Correct(...cbc.Option) error
}

// NewDocument instantiates a Document wrapper around the provided object.
func NewDocument(payload interface{}) (*Document, error) {
	d := new(Document)
	return d, d.insert(payload)
}

// Digest calculates a digital digest using the canonical JSON of the document.
func (d *Document) Digest() (*dsig.Digest, error) {
	data, err := json.Marshal(d)
	if err != nil {
		return nil, ErrMarshal.WithCause(err)
	}
	r := bytes.NewReader(data)
	cd, err := c14n.CanonicalJSON(r)
	if err != nil {
		return nil, ErrInternal.WithErrorf("canonical JSON error: %w", err)
	}
	return dsig.NewSHA256Digest(cd), nil
}

// IsEmpty returns true if no payload has been set yet.
func (d *Document) IsEmpty() bool {
	return d.payload == nil
}

// Schema provides the document's schema.
func (d *Document) Schema() schema.ID {
	return d.schema
}

// Instance returns a prepared version of the document's content.
func (d *Document) Instance() interface{} {
	return d.payload
}

// Normalize will attempt to run the normalization method on the
// document payload.
func (d *Document) Normalize() error {
	pl, ok := d.payload.(Normalizable)
	if !ok {
		return nil
	}
	return pl.Normalize()
}

// Calculate wraps around the Normalize method.
//
// Deprecated: use the Normalize method instead.
func (d *Document) Calculate() error {
	return d.Normalize()
}

// Validate checks to ensure the document has everything it needs
// and will pass on the validation call to the payload.
func (d *Document) Validate() error {
	return d.ValidateWithContext(context.Background())
}

// ValidateWithContext checks to ensure the document has everything it needs
// and will pass on the validation call to the payload.
func (d *Document) ValidateWithContext(ctx context.Context) error {
	if ctx.Value(internal.KeyDraft) == nil {
		// if draft not set previously, assume true
		ctx = context.WithValue(ctx, internal.KeyDraft, true)
	}
	err := validation.ValidateStructWithContext(ctx, d,
		validation.Field(&d.schema, validation.Required),
	)
	if err != nil {
		return err
	}
	// return any errors from the payload as if they were for the document
	// itself.
	return validation.ValidateWithContext(ctx, d.payload)
}

// Correct will attempt to run the correction method on the document
// using some of the provided options.
func (d *Document) Correct(opts ...cbc.Option) error {
	pl, ok := d.payload.(Correctable)
	if !ok {
		return errors.New("document cannot be corrected")
	}
	if err := pl.Correct(opts...); err != nil {
		return err
	}
	return nil
}

// Insert places the provided object inside the document and looks up the schema
// information to ensure it is known.
func (d *Document) insert(payload interface{}) error {
	d.schema = schema.Lookup(payload)
	if d.schema == schema.UnknownID {
		return ErrMarshal.WithErrorf("unregistered or invalid schema")
	}
	d.payload = payload
	return nil
}

// Clone makes a copy of the document by serializing and deserializing it.
// the contents into a new document instance.
func (d *Document) Clone() (*Document, error) {
	d2 := new(Document)
	data, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, d2); err != nil {
		return nil, err
	}
	return d2, nil
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
func (d *Document) UnmarshalJSON(data []byte) error {
	var err error
	if d.schema, err = schema.Extract(data); err != nil {
		return ErrUnknownSchema.WithCause(err)
	}

	// Map the schema to an instance of the payload, or fail if we don't know what it is
	d.payload = d.schema.Interface()
	if d.payload == nil {
		return ErrMarshal.WithErrorf("unregistered or invalid schema")
	}
	if err := json.Unmarshal(data, d.payload); err != nil {
		return err
	}

	return nil
}

// MarshalJSON satisfies the json.Marshaler interface.
func (d *Document) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(d.payload)
	if err != nil {
		return nil, ErrMarshal.WithCause(err)
	}

	data, err = schema.Insert(d.schema, data)
	if err != nil {
		return nil, ErrMarshal.WithCause(err)
	}

	return data, nil
}

// JSONSchema returns a jsonschema.Schema instance.
func (Document) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:  "object",
		Title: "Document",
		Description: here.Doc(`
			Contains the document payload to be included inside an Envelope.
			
			The document must contain a ` + "`" + `$schema` + "`" + ` property that identifies
			the data's structure otherwise it will be rejected.`),
	}
}
