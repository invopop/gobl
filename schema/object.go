package schema

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/invopop/gobl/internal"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/uuid"
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

// Object helps handle json objects that must contain a schema to correctly identify
// the contents and ensuring that a `$schema` property is added automatically when
// marshalling back into JSON.
type Object struct {
	Schema  ID `json:"$schema"`
	payload any
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
	CorrectionOptionsSchema() (any, error)
}

// Replicable defines the methods expected of a document payload that can be replicated.
type Replicable interface {
	Replicate() error
}

// Identifiable defines the methods expected of a document payload that contains a UUID.
// The `uuid` packages `Identify` struct can be embedded to satisfy this.
type Identifiable interface {
	GetUUID() uuid.UUID
	SetUUID(uuid.UUID)
}

// NewObject instantiates an Object wrapper around the provided payload.
func NewObject(payload interface{}) (*Object, error) {
	d := new(Object)
	return d, d.insert(payload)
}

// IsEmpty returns true if no payload has been set yet.
func (d *Object) IsEmpty() bool {
	return d.payload == nil
}

// Instance returns a prepared version of the document's content.
func (d *Object) Instance() interface{} {
	return d.payload
}

// Calculate will attempt to run the calculation method on the
// document payload. If the object implements the Identifiable
// interface, it will also ensure the UUID is set.
func (d *Object) Calculate() error {
	if ident, ok := d.payload.(Identifiable); ok {
		id := ident.GetUUID()
		if id.IsZero() {
			ident.SetUUID(uuid.V7())
		}
	}
	pl, ok := d.payload.(Calculable)
	if !ok {
		return nil
	}
	return pl.Calculate()
}

// Validate checks to ensure the document has everything it needs
// and will pass on the validation call to the payload.
func (d *Object) Validate() error {
	return d.ValidateWithContext(context.Background())
}

// ValidateWithContext checks to ensure the document has everything it needs
// and will pass on the validation call to the payload.
func (d *Object) ValidateWithContext(ctx context.Context) error {
	if ctx.Value(internal.KeyDraft) == nil {
		// if draft not set previously, assume true
		ctx = context.WithValue(ctx, internal.KeyDraft, true)
	}
	err := validation.ValidateStructWithContext(ctx, d,
		validation.Field(&d.Schema, validation.Required),
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
func (d *Object) Correct(opts ...Option) error {
	pl, ok := d.payload.(Correctable)
	if !ok {
		return errors.New("document cannot be corrected")
	}
	if err := pl.Correct(opts...); err != nil {
		return err
	}
	return nil
}

// CorrectionOptionsSchema provides a schema with the correction options available
// for the schema, if available.
func (d *Object) CorrectionOptionsSchema() (any, error) {
	pl, ok := d.payload.(Correctable)
	if !ok {
		return nil, nil
	}
	res, err := pl.CorrectionOptionsSchema()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Replicate will attempt to clone and run the Replicate method of the object
// if it has one.
func (d *Object) Replicate() error {
	obj, ok := d.payload.(Replicable)
	if ok {
		if err := obj.Replicate(); err != nil {
			return err
		}
	}
	if ident, ok := d.payload.(Identifiable); ok {
		id := ident.GetUUID()
		if id.IsZero() {
			ident.SetUUID(uuid.V7())
		}
	}
	return nil
}

// Insert places the provided object inside the document and looks up the schema
// information to ensure it is known.
func (d *Object) insert(payload interface{}) error {
	d.Schema = Lookup(payload)
	if d.Schema == UnknownID {
		return ErrUnknownSchema
	}
	d.payload = payload
	return nil
}

// UUID extracts the UUID from the payload using reflection. An empty
// id is returned if the payload does not have a UUID field.
func (d *Object) UUID() uuid.UUID {
	obj, ok := d.payload.(Identifiable)
	if !ok {
		return uuid.Empty
	}
	return obj.GetUUID()
}

// Clone makes a copy of the document by serializing and deserializing
// the contents into a new document instance.
func (d *Object) Clone() (*Object, error) {
	d2 := new(Object)
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
func (d *Object) UnmarshalJSON(data []byte) error {
	var err error
	if d.Schema, err = Extract(data); err != nil {
		return err
	}
	if d.Schema == UnknownID {
		return nil // return silently
	}

	// Map the schema to an instance of the payload, or fail if we don't know what it is
	d.payload = d.Schema.Interface()
	if d.payload == nil {
		return ErrUnknownSchema
	}
	if err := json.Unmarshal(data, d.payload); err != nil {
		return err
	}

	return nil
}

// MarshalJSON satisfies the json.Marshaler interface.
func (d *Object) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(d.payload)
	if err != nil {
		return nil, err
	}

	data, err = Insert(d.Schema, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// JSONSchema returns a jsonschema.Schema instance.
func (Object) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:  "object",
		Title: "Object",
		Description: here.Doc(`
			Data object whose type is determined from the <code>$schema</code> property.
		`),
	}
}
