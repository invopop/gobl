package cbc

import (
	"errors"

	"github.com/invopop/jsonschema"
)

// KeyOrCode is a special type that can be either a Key or a Code. This is meant
// for situations when the value to be used has to be flexible enough
// to either a key defined by GOBL, or a code usually defined by an external
// entity.
type KeyOrCode string

// String provides the string representation.
func (kc KeyOrCode) String() string {
	return string(kc)
}

// Key returns the key value or empty if the value is a Code.
func (kc KeyOrCode) Key() Key {
	k := Key(kc)
	if err := k.Validate(); err == nil {
		return k
	}
	return KeyEmpty
}

// Code returns the code value or empty if the value is a Key.
func (kc KeyOrCode) Code() Code {
	c := Code(kc)
	if err := c.Validate(); err == nil {
		return c
	}
	return CodeEmpty
}

// Validate ensures the value is either a key or a code.
func (kc KeyOrCode) Validate() error {
	if err := Key(kc).Validate(); err == nil {
		return nil
	}
	if err := Code(kc).Validate(); err == nil {
		return nil
	}
	return errors.New("value is not a key or code")
}

func (KeyOrCode) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.OneOf = []*jsonschema.Schema{
		KeyEmpty.JSONSchema(),
		CodeEmpty.JSONSchema(),
	}
}
