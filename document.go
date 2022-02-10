package gobl

import (
	"bytes"
	"encoding/json"

	"github.com/alecthomas/jsonschema"
	"github.com/invopop/gobl/c14n"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/region"
	"github.com/invopop/gobl/schema"
)

// Document helps us handle the document's contents by essentially wrapping around
// the json RawMessage.
type Document struct {
	data json.RawMessage
}

// Calculable defines the methods expected of a document payload that contains a `Calculate`
// method to be used to perform any additional calculations.
type Calculable interface {
	Calculate(r region.Region) error
}

// Validatable describes a document that can be validated.
type Validatable interface {
	Validate(r region.Region) error
}

// Type provides the payload documents schema type
func (p *Document) Type() (schema.Type, error) {
	def, err := p.Def()
	if err != nil {
		return schema.UnknownType, err
	}
	return def.Schema.Type(), nil
}

// Def extracts the schema def from the document
func (p *Document) Def() (schema.Def, error) {
	def := schema.Def{}
	err := json.Unmarshal(p.data, &def)
	return def, err
}

func (p *Document) insert(doc interface{}) error {
	var err error

	p.data, err = json.Marshal(doc)
	if err != nil {
		return ErrMarshal.WithCause(err)
	}
	return nil
}

func (p *Document) extract(doc interface{}) error {
	return json.Unmarshal(p.data, doc)
}

func (p *Document) digest() (*dsig.Digest, error) {
	r := bytes.NewReader(p.data)
	cd, err := c14n.CanonicalJSON(r)
	if err != nil {
		return nil, ErrInternal.WithErrorf("canonical JSON error: %w", err)
	}
	return dsig.NewSHA256Digest(cd), nil
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
func (p *Document) UnmarshalJSON(data []byte) error {
	p.data = json.RawMessage(data)
	return nil
}

// MarshalJSON satisfies the json.Marshaler interface.
func (p *Document) MarshalJSON() ([]byte, error) {
	return p.data, nil
}

// JSONSchemaType returns a jsonschema.Type object.
func (Document) JSONSchemaType() *jsonschema.Type {
	return &jsonschema.Type{
		Type:        "object",
		Title:       "Document",
		Description: "Contents of the envelope",
	}
}
