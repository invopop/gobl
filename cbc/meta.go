package cbc

import (
	"github.com/invopop/jsonschema"
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

// Equals checks if the meta data is the same.
func (m Meta) Equals(m2 Meta) bool {
	if len(m) != len(m2) {
		return false
	}
	for k, v := range m {
		if m2[k] != v {
			return false
		}
	}
	return true
}

// JSONSchemaExtend ensures the meta keys are valid.
func (Meta) JSONSchemaExtend(schema *jsonschema.Schema) {
	prop := schema.AdditionalProperties
	schema.AdditionalProperties = nil
	schema.PatternProperties = map[string]*jsonschema.Schema{
		KeyPattern: prop,
	}
}
