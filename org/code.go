package org

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/jsonschema"
)

// Code represents a string used to uniquely identify the data we're looking
// at. We use "code" instead of "id", to reenforce the fact that codes should
// be more easily set and used by humans within definitions than IDs or UUIDs.
// Codes are standardised so that when validated they must contain between
// 2 and 6 inclusive upper-case letters or numbers.
type Code string

var (
	codePattern          = `^[A-Z0-9]{1,6}$`
	codeValidationRegexp = regexp.MustCompile(codePattern)
)

// CodeEmpty is used when no code is defined.
const CodeEmpty Code = ""

// Validate ensures that the code complies with the expected rules.
func (c Code) Validate() error {
	return validation.Validate(string(c),
		validation.Length(1, 6),
		validation.Match(codeValidationRegexp),
	)
}

// IsEmpty returns true if no code is specified.
func (c Code) IsEmpty() bool {
	return c == CodeEmpty
}

// String returns string representation of code.
func (c Code) String() string {
	return string(c)
}

// In returns true if the code's value matches one of those
// in the provided list.
func (c Code) In(ary ...Code) bool {
	for _, v := range ary {
		if v == c {
			return true
		}
	}
	return false
}

// JSONSchema provides a representation of the struct for usage in Schema.
func (Code) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Pattern:     codePattern,
		Title:       "Code",
		MinLength:   1,
		MaxLength:   6,
		Description: "Short upper-case identifier.",
	}
}
