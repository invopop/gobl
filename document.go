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
	Schema schema.ID
	data   json.RawMessage
}

type schemaDoc struct {
	Schema schema.ID `json:"$schema,omitempty"`
}

// Insert places the provided object inside the document and looks up the schema
// information to ensure it is known.
func (p *Document) Insert(doc interface{}) error {
	p.Schema = schema.Lookup(doc)
	if p.Schema == schema.UnknownID {
		return ErrMarshal.WithErrorf("unregistered schema")
	}

	data, err := json.Marshal(doc)
	if err != nil {
		return ErrMarshal.WithCause(err)
	}

	// Combine the base data with the JSON schema information.
	// We manually create and add the JSON as this is just simply the quickest
	// way to do it.
	buf := bytes.NewBufferString(`{"$schema":"` + p.Schema.String() + `",`)
	_, _ = buf.Write(bytes.TrimLeft(data, "{")) //nolint:errcheck
	p.data = buf.Bytes()

	return nil
}

// Extract will unmarshal the documents contents into the provide object. You'll
// need have checked the type proviously to ensure this works.
func (p *Document) extract(doc interface{}) error {
	if err := json.Unmarshal(p.data, doc); err != nil {
		return ErrMarshal.WithCause(err)
	}
	return nil
}

// Digest calculates a digital digest using the canonical JSON of the embedded
// data.
func (p *Document) Digest() (*dsig.Digest, error) {
	r := bytes.NewReader(p.data)
	cd, err := c14n.CanonicalJSON(r)
	if err != nil {
		return nil, ErrInternal.WithErrorf("canonical JSON error: %w", err)
	}
	return dsig.NewSHA256Digest(cd), nil
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
func (p *Document) UnmarshalJSON(data []byte) error {
	def := new(schemaDoc)
	if err := json.Unmarshal(data, def); err != nil {
		return err
	}
	p.Schema = def.Schema
	p.data = json.RawMessage(data)
	return nil
}

// MarshalJSON satisfies the json.Marshaler interface.
func (p *Document) MarshalJSON() ([]byte, error) {
	return p.data, nil
}

// JSONSchema returns a jsonschema.Schema instance.
func (Document) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "object",
		Title:       "Document",
		Description: "Contents of the envelope",
	}
}
