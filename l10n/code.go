package l10n

import (
	"regexp"

	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Code is used for short identifiers like country or state codes.
// They are limited to upper-case letters and numbers
// only, and should be validated against region specific data.
type Code string

var (
	codePattern = `^[A-Z0-9]+$`
	codeFormat  = regexp.MustCompile(codePattern)
)

// CodeEmpty is used for matching empty codes.
const CodeEmpty Code = ""

// Validate ensures the code is formatted correctly.
func (c Code) Validate() error {
	return validation.Match(codeFormat).Validate(string(c))
}

// In returns true if the code is contained inside the provided set
func (c Code) In(set ...Code) bool {
	for _, x := range set {
		if c == x {
			return true
		}
	}
	return false
}

// String provides string representation of code
func (c Code) String() string {
	return string(c)
}

// ISO is a convenience method to provide the ISO country code of
// a code.
func (c Code) ISO() ISOCountryCode {
	return ISOCountryCode(c)
}

// Tax is a convenience method to provide the tax country code of a
// code.
func (c Code) Tax() TaxCountryCode {
	return TaxCountryCode(c)
}

// JSONSchema provides a representation of the struct for usage in Schema.
func (Code) JSONSchema() *jsonschema.Schema {
	s := &jsonschema.Schema{
		Title:       "Code",
		Type:        "string",
		Pattern:     codePattern,
		Description: "Code is used for short identifies like country or state codes.",
	}
	return s
}
