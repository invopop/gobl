package l10n

import (
	"errors"
	"regexp"
)

// Code is used for identifiers for specific states, cities,
// or localities. They are limited to upper-case letters and numbers
// only, and should be validated against region specific data.
type Code string

var codeFormat = regexp.MustCompile(`\A[A-Z0-9]+\z`)

// Validate ensures the code is formatted correctly.
func (c Code) Validate() error {
	if string(c) == "" {
		return nil
	}
	if !codeFormat.MatchString(string(c)) {
		return errors.New("invalid code format")
	}
	return nil
}
