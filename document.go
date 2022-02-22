package gobl

import (
	"bytes"
	"encoding/json"

	"github.com/invopop/gobl/c14n"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/jsonschema"
)

// Document helps us handle the document's contents by essentially wrapping around
// the contents and ensuring that a `$schema` property is added automatically when
// marshalling into JSON.
type Document struct {
	schema schema.ID
	obj    interface{}
}

type schemaDoc struct {
	Schema schema.ID `json:"$schema,omitempty"`
}

// NewDocument instantiates a Document wrapper around the provided object.
func NewDocument(obj interface{}) (*Document, error) {
	d := new(Document)
	return d, d.insert(obj)
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

// Schema provides the document's schema.
func (d *Document) Schema() schema.ID {
	return d.schema
}

// Instance returns a prepared version of the document's content.
func (d *Document) Instance() interface{} {
	return d.obj
}

// Insert places the provided object inside the document and looks up the schema
// information to ensure it is known.
func (d *Document) insert(doc interface{}) error {
	d.schema = schema.Lookup(doc)
	if d.schema == schema.UnknownID {
		return ErrMarshal.WithErrorf("unregistered schema")
	}
	d.obj = doc
	return nil
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
func (d *Document) UnmarshalJSON(data []byte) error {
	def := new(schemaDoc)
	if err := json.Unmarshal(data, def); err != nil {
		return err
	}
	d.schema = def.Schema

	// Map the schema to an instance of the object, or fail if we don't know what it is
	d.obj = d.schema.Interface()
	if d.obj == nil {
		return ErrMarshal.WithErrorf("unregistered schema: %v", d.schema.String())
	}
	if err := json.Unmarshal(data, d.obj); err != nil {
		return err
	}

	return nil
}

// MarshalJSON satisfies the json.Marshaler interface.
func (d *Document) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(d.obj)
	if err != nil {
		return nil, ErrMarshal.WithCause(err)
	}

	sdata, err := json.Marshal(d.schemaDoc())
	if err != nil {
		return nil, ErrMarshal.WithCause(err)
	}

	// Combine the base data with the JSON schema information.
	// We manually create and add the JSON as this is just simply the quickest
	// way to do it.
	data = bytes.TrimLeft(data, "{")
	sdata = append(bytes.TrimRight(sdata, "}"), byte(','))
	data = append(sdata, data...)

	return data, nil
}

func (d *Document) schemaDoc() *schemaDoc {
	return &schemaDoc{
		Schema: d.schema,
	}
}

// JSONSchema returns a jsonschema.Schema instance.
func (Document) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "object",
		Title:       "Document",
		Description: "Contents of the envelope",
	}
}
