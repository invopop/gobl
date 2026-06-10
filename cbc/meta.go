package cbc

import (
	"bytes"
	"encoding/json"
	"iter"
	"sort"

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

// Keys returns all the keys of the Meta sorted alphabetically.
func (m Meta) Keys() []Key {
	if len(m) == 0 {
		return nil
	}
	keys := make([]Key, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

// Values returns all the values of the Meta in the order of their keys'
// alphabetical sorting.
func (m Meta) Values() []string {
	if len(m) == 0 {
		return nil
	}
	keys := m.Keys()
	values := make([]string, len(keys))
	for i, k := range keys {
		values[i] = m[k]
	}
	return values
}

// All returns an iterator over the Meta entries in alphabetical order of
// the keys. Intended for use with Go 1.23+ range-over-func:
//
//	for k, v := range m.All() {
//	    // ...
//	}
func (m Meta) All() iter.Seq2[Key, string] {
	return func(yield func(Key, string) bool) {
		for _, k := range m.Keys() {
			if !yield(k, m[k]) {
				return
			}
		}
	}
}

// MarshalJSON emits the Meta as a JSON object with keys sorted
// alphabetically for deterministic output. An empty or nil Meta marshals
// to "null", which combines with `json:"meta,omitempty"` on map fields to
// omit the field from the encoded parent object.
func (m Meta) MarshalJSON() ([]byte, error) {
	if len(m) == 0 {
		return []byte("null"), nil
	}
	keys := m.Keys()
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		kb, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}
		buf.Write(kb)
		buf.WriteByte(':')
		vb, err := json.Marshal(m[k])
		if err != nil {
			return nil, err
		}
		buf.Write(vb)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// JSONSchemaExtend ensures the meta keys are valid.
func (Meta) JSONSchemaExtend(schema *jsonschema.Schema) {
	prop := schema.AdditionalProperties
	schema.AdditionalProperties = nil
	schema.PatternProperties = map[string]*jsonschema.Schema{
		KeyPattern: prop,
	}
}
