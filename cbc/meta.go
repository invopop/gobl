package cbc

import (
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Meta defines a structure for data about the data being defined.
// Typically would be used for adding additional IDs or specifications
// not already defined or required by the base structure.
//
// GOBL is focussed on ensuring the recipient has everything they need,
// as such, meta should only be used for data that may be used by intermediary
// conversion processes that should not be needed by the end-user.
//
// We need to always use strings for values so that meta-data is easy to convert
// into other formats, such as protobuf which has strict type requirements.
type Meta map[Key]string

// Validate ensures the meta data looks correct.
func (m Meta) Validate() error {
	// we're only interested in checking the keys
	err := make(validation.Errors)
	for k := range m {
		if e := k.Validate(); e != nil {
			err[k.String()] = e
		}
	}
	if len(err) == 0 {
		return nil
	}
	return err
}

// JSONSchemaExtend ensures the meta keys are valid.
func (Meta) JSONSchemaExtend(schema *jsonschema.Schema) {
	prop := schema.PatternProperties[".*"] // get default
	delete(schema.PatternProperties, ".*") // remove default
	schema.PatternProperties[KeyPattern] = prop
}
