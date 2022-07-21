package l10n

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Code is used for short identifiers like country or state codes.
// They are limited to upper-case letters and numbers
// only, and should be validated against region specific data.
type Code string

var (
	codeFormat = regexp.MustCompile(`\A[A-Z0-9]+\z`)
)

// CodeEmpty is used for matching empty codes.
const CodeEmpty Code = ""

// Validate ensures the code is formatted correctly.
func (c Code) Validate() error {
	return validation.Validate(string(c),
		validation.Match(codeFormat),
	)
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
