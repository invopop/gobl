package schema

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/invopop/gobl/internal"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Error is used to define schema errors
type Error string

// Error provides the error code
func (e Error) Error() string {
	return string(e)
}

const (
	// ErrUnknownSchema is returned when the schema has not been registered.
	ErrUnknownSchema Error = "unknown-schema"
)

// Document helps us handle the document's contents by essentially wrapping around
// the contents and ensuring that a `$schema` property is added automatically when
// marshalling into JSON.
type Document struct {
	schema  ID
	payload interface{}
}

// Calculable defines the methods expected of a document payload that contains a `Calculate`
// method to be used to perform any additional calculations.
type Calculable interface {
	Calculate() error
}

// Correctable defines the expected interface of a document that can be
// corrected.
type Correctable interface {
	Correct(...Option) error
}

// NewDocument instantiates a Document wrapper around the provided object.
func NewDocument(payload interface{}) (*Document, error) {
	d := new(Document)
	return d, d.insert(payload)
}

// IsEmpty returns true if no payload has been set yet.
func (d *Document) IsEmpty() bool {
	return d.payload == nil
}

// Schema provides the document's schema.
func (d *Document) Schema() ID {
	return d.schema
}

// Instance returns a prepared version of the document's content.
func (d *Document) Instance() interface{} {
	return d.payload
}

// Calculate will attempt to run the calculation method on the
// document payload.
func (d *Document) Calculate() error {
	pl, ok := d.payload.(Calculable)
	if !ok {
		return nil
	}
	return pl.Calculate()
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
func (d *Document) Correct(opts ...Option) error {
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
	d.schema = Lookup(payload)
	if d.schema == UnknownID {
		return ErrUnknownSchema
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
	if d.schema, err = Extract(data); err != nil {
		return fmt.Errorf("%w: %s", ErrUnknownSchema, err.Error())
	}

	// Map the schema to an instance of the payload, or fail if we don't know what it is
	d.payload = d.schema.Interface()
	if d.payload == nil {
		return err
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
		return nil, err
	}

	data, err = Insert(d.schema, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// JSONSchema returns a jsonschema.Schema instance.
func (Document) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:  "object",
		Title: "Document",
		Description: here.Doc(`
			Data object whose type is determined from the <code>$schema</code> property.
		`),
	}
}
