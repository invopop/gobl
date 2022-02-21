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
// the json RawMessage.
type Document struct {
	schema schema.ID
	obj    interface{}
}

type schemaDoc struct {
	Schema schema.ID `json:"$schema,omitempty"`
}

// Digest calculates a digital digest using the canonical JSON of the document.
func (p *Document) Digest() (*dsig.Digest, error) {
	data, err := json.Marshal(p)
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
func (p *Document) Schema() schema.ID {
	return p.schema
}

// Instance returns a prepared version of the document's content.
func (p *Document) Instance() interface{} {
	return p.obj
}

// Insert places the provided object inside the document and looks up the schema
// information to ensure it is known.
func (p *Document) insert(doc interface{}) error {
	p.schema = schema.Lookup(doc)
	if p.schema == schema.UnknownID {
		return ErrMarshal.WithErrorf("unregistered schema")
	}
	p.obj = doc
	return nil
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
func (p *Document) UnmarshalJSON(data []byte) error {
	def := new(schemaDoc)
	if err := json.Unmarshal(data, def); err != nil {
		return err
	}
	p.schema = def.Schema

	// Map the schema to an instance of the object, or fail if we don't know what it is
	p.obj = p.schema.Interface()
	if p.obj == nil {
		return ErrMarshal.WithErrorf("unregistered schema: %v", p.schema.String())
	}
	if err := json.Unmarshal(data, p.obj); err != nil {
		return err
	}

	return nil
}

// MarshalJSON satisfies the json.Marshaler interface.
func (p *Document) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(p.obj)
	if err != nil {
		return nil, ErrMarshal.WithCause(err)
	}

	// Combine the base data with the JSON schema information.
	// We manually create and add the JSON as this is just simply the quickest
	// way to do it.
	buf := bytes.NewBufferString(`{"$schema":"` + p.schema.String() + `",`)
	_, _ = buf.Write(bytes.TrimLeft(data, "{")) //nolint:errcheck

	return buf.Bytes(), nil
}

// JSONSchema returns a jsonschema.Schema instance.
func (Document) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "object",
		Title:       "Document",
		Description: "Contents of the envelope",
	}
}
