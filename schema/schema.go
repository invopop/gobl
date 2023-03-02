// Package schema provides a simple way to register and lookup schemas.
package schema

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

const (
	// VERSION for the current version of the schema
	VERSION = "draft-0"
	// GOBL stores the base schema ID for GOBL, including current schema version.
	GOBL ID = "https://gobl.org/" + VERSION
)

const (
	// UnknownID is provided when the schema has not been registered
	UnknownID ID = ""
)

func init() {
	schemas = newRegistry()
}

// ID contains the official schema URL.
type ID string

type document struct {
	Schema ID `json:"$schema,omitempty"`
}

// Extract attempts to Unmarshal the provided JSON document in order to extract
// the payload's Schema ID.
func Extract(data []byte) (ID, error) {
	def := new(document)
	if err := json.Unmarshal(data, def); err != nil {
		return UnknownID, err
	}
	return def.Schema, nil
}

// Insert adds the provided schema ID to the JSON data provided.
func Insert(id ID, data []byte) ([]byte, error) {
	doc := &document{Schema: id}
	sdata, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	// Combine the base data with the JSON schema information.
	// We manually create and add the JSON as this is just simply the quickest
	// way to do it.
	data = bytes.TrimLeft(data, "{")
	sdata = append(bytes.TrimRight(sdata, "}"), byte(','))
	data = append(sdata, data...)

	return data, nil
}

// Validate ensures the schema ID looks good.
func (id ID) Validate() error {
	return validation.Validate(string(id), is.URL)
}

// Anchor either adds or replaces the anchor part of the schema URI.
func (id ID) Anchor(name string) ID {
	b := id.Base()
	return ID(b.String() + "#" + name)
}

// Add appends the provided path to the id, and removes any
// anchor data that might be there.
func (id ID) Add(path string) ID {
	b := id.Base()
	path = strings.TrimLeft(path, "/")
	return ID(b.String() + "/" + ToSnakeCase(path))
}

// Base removes any anchor information from the schema
func (id ID) Base() ID {
	s := id.String()
	i := strings.LastIndex(s, "#")
	if i != -1 {
		s = s[0:i]
	}
	s = strings.TrimRight(s, "/")
	return ID(s)
}

// String provides string version of ID
func (id ID) String() string {
	return string(id)
}

// Interface attempts to determine the type by looking up the ID in the
// registered list of schemas, and providing an empty instance.
func (id ID) Interface() interface{} {
	typ := Type(id)
	if typ == nil {
		return nil
	}
	return reflect.New(typ).Interface()
}
