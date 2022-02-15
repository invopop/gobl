package schema

import (
	"reflect"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
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
