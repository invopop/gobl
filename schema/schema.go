package schema

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	// VERSION for the current version of the schema
	VERSION    = "draft0"
	GOBL    ID = "https://gobl.org/" + VERSION
)

const (
	// UnknownID is provided when the schema has not been registered
	UnknownID ID = ""
)

const (
	// idURLTemplate defines the base for how GOBL Schema definitions should be named.
	idURLTemplate = "https://gobl.org/%s/%s"
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
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return ID(b.String() + path)
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
