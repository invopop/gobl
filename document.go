package gobl

import (
	"encoding/json"
	"fmt"
)

// Header defines the meta data of the body.
type Header struct {
	Type   BodyType `json:"type"`
	Digest *Digest  `json:"digest"`
}

// Digest of the body
type Digest struct {
	Algorithm string `json:"alg" jsonschema:"title=Algorithm"`
	Value     string `json:"val" jsonschema:"title=Value"`
}

// Document defines the base GoBL structure used for serializing data
// for persistence and sharing.
type Document struct {
	Head *Header         `json:"head" jsonschema:"title=Header contents"`
	Body json.RawMessage `json:"body" jsonschema:"title=Raw document payload"`
	Sigs []*Signature    `json:"sigs" jsonschema:"title=Signatures"`

	body Body // instance of body
}

// SetBody set's the documents body content, marhsals to JSON,
// and updates the header digest.
func (d *Document) SetBody(b Body) error {
	var err error
	d.Body, err = json.Marshal(b)
	if err != nil {
		return fmt.Errorf("setting document body: %w", err)
	}

	// Determine the digest
	// TODO!

	d.body = b
	return nil
}
