package gobl

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/alecthomas/jsonschema"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/c14n"
	"github.com/invopop/gobl/dsig"
)

// Envelope wraps around a gobl document and provides support for digest creation
// and digital signatures.
type Envelope struct {
	Head       *Header           `json:"head" jsonschema:"title=Header,description=Details on what the contents are"`
	Document   *Payload          `json:"doc" jsonschema:"title=Document,description=The data being enveloped"`
	Signatures []*dsig.Signature `json:"sigs" jsonschema:"title=Signatures,description=JSON Web Signatures of the header"`
}

// Document defines what we expect from a document to be able to be included in an envelope.
type Document interface {
	Type() string
}

// Validate ensures that the envelope contains everything it should to be considered valid GoBL.
func (e *Envelope) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.Head, validation.Required),
		validation.Field(&e.Document, validation.Required),
		validation.Field(&e.Signatures, validation.Required),
	)
}

// Verify ensures the digest headers still match the document contents.
func (e *Envelope) Verify() error {
	d1 := e.Head.Digest
	d2, err := e.Document.digest()
	if err != nil {
		return err
	}
	return d1.Equals(d2)
}

// Sign uses the private key to the envelope headers.
func (e *Envelope) Sign(key *dsig.PrivateKey) error {
	sig, err := key.Sign(e.Head)
	if err != nil {
		return err
	}
	e.Signatures = append(e.Signatures, sig)
	return nil
}

// Insert takes the provided document an serializes it ready for use.
func (e *Envelope) Insert(doc Document) error {
	if e.Document == nil {
		e.Document = new(Payload)
	}
	err := e.Document.insert(doc)
	if err != nil {
		return err
	}

	if e.Head == nil {
		e.Head = new(Header)
	}
	e.Head.Type = doc.Type()

	e.Head.Digest, err = e.Document.digest()
	if err != nil {
		return err
	}

	return nil
}

// Extract the contents of the envelope into the provided document type.
func (e *Envelope) Extract(doc Document) error {
	return e.Document.extract(doc)
}

// Payload helps us handle the document's contents by essentially wrapping around
// the json RawMessage.
type Payload struct {
	data json.RawMessage
}

func (p *Payload) insert(doc Document) error {
	var err error
	p.data, err = json.Marshal(doc)
	return err
}

func (p *Payload) extract(doc Document) error {
	return json.Unmarshal(p.data, doc)
}

func (p *Payload) digest() (*dsig.Digest, error) {
	r := bytes.NewReader(p.data)
	cd, err := c14n.CanonicalJSON(r)
	if err != nil {
		return nil, fmt.Errorf("canonical JSON error: %w", err)
	}
	return dsig.NewSHA256Digest(cd), nil
}

func (p *Payload) UnmarshalJSON(data []byte) error {
	p.data = json.RawMessage(data)
	return nil
}

func (p *Payload) MarshalJSON() ([]byte, error) {
	return p.data, nil
}

func (Payload) JSONSchemaType() *jsonschema.Type {
	return &jsonschema.Type{
		Type:        "object",
		Title:       "Payload",
		Description: "Contents of the envelope",
	}
}
