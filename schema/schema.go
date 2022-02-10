package schema

import (
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	// VERSION for the current version of the schema
	VERSION = "draft0"
)

const (
	// UnknownType is provided when the type cannot be extracted
	UnknownType Type = ""
)

const (
	// idURLTemplate defines the base for how GOBL Schema definitions should be named.
	idURLTemplate = "https://gobl.org/%s/%s"
)

// Type represents an internal name
type Type string

// Def defines the base to use inside documents that require a schema definition.
type Def struct {
	// The GOBL Schema and version used to generate the envelope.
	Schema ID `json:"$schema" jsonschema:"title=Schema"`
}

// ID contains the official schema URL.
type ID string

// Validate ensures the schema ID looks good.
func (id ID) Validate() error {
	return validation.Validate(string(id), is.URL)
}

// IsGOBL returns true, if the schema ID looks like it belongs to
// GOBL.
func (id ID) IsGOBL() bool {
	return strings.HasPrefix(string(id), "https://gobl.org/")
}

// Version provides the GOBL version or empty string it isn't GOBL
func (id ID) Version() string {
	if !id.IsGOBL() {
		return ""
	}
	s := id.split()
	return s[2]
}

// Type extracts the schema type from the ID.
func (id ID) Type() Type {
	if !id.IsGOBL() {
		return UnknownType
	}
	s := id.split()
	return Type(s[4])
}

func (id ID) split() []string {
	return strings.SplitN(string(id), "/", 5)
}

func For(typ Type) ID {
	s := fmt.Sprintf(idURLTemplate, VERSION, string(typ))
	return ID(s)
}

// ID returns the complete schema ID for the provided type
func (t Type) ID() ID {
	return For(t)
}

// Def returns a schema definition model for the type
func (t Type) Def() Def {
	return Def{
		Schema: t.ID(),
	}
}
