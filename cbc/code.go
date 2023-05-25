package cbc

import (
	"regexp"

	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Code represents a string used to uniquely identify the data we're looking
// at. We use "code" instead of "id", to reenforce the fact that codes should
// be more easily set and used by humans within definitions than IDs or UUIDs.
// Codes are standardised so that when validated they must contain between
// 1 and 24 inclusive upper-case letters or numbers with optional periods
// to separate blocks.
type Code string

// CodeSet is a map of keys to specific codes, useful to determine regime specific
// codes from their key counterparts.
type CodeSet map[Key]Code

// Basic code constants.
const (
	CodePattern   = `^[A-Z0-9]+(\.?[A-Z0-9]+)*$`
	CodeMinLength = 1
	CodeMaxLength = 24
)

var (
	codeValidationRegexp = regexp.MustCompile(CodePattern)
)

// CodeEmpty is used when no code is defined.
const CodeEmpty Code = ""

// Validate ensures that the code complies with the expected rules.
func (c Code) Validate() error {
	return validation.Validate(string(c),
		validation.Length(1, CodeMaxLength),
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
		Pattern:     CodePattern,
		Title:       "Code",
		MinLength:   CodeMinLength,
		MaxLength:   CodeMaxLength,
		Description: "Alphanumerical text identifier with upper-case letters, no whitespace, nor symbols.",
	}
}

// Validate ensures the code set data looks correct.
func (cs CodeSet) Validate() error {
	err := make(validation.Errors)
	// values are already tested
	for k := range cs {
		if e := k.Validate(); e != nil {
			err[k.String()] = e
		}
	}
	if len(err) == 0 {
		return nil
	}
	return err
}
