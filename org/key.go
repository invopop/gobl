package org

import (
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/jsonschema"
)

// Key is used to define an ID or code that more closely represents
// a human name. The objective is to make it easier to define constants
// that can be re-used more easily.
type Key string

var (
	// KeyPattern describes what should keys look like
	KeyPattern = `^[a-z0-9][a-z0-9-+]*[a-z0-9]$`
	// KeyValidationRegexp is used for key validation
	KeyValidationRegexp = regexp.MustCompile(KeyPattern)
)

// KeyEmpty is used when no key is available.
const KeyEmpty Key = ""

// Validate ensures the key complies with the basic syntax
// requirements.
func (k Key) Validate() error {
	return validation.Validate(string(k),
		validation.Match(KeyValidationRegexp),
		validation.Length(2, 64),
	)
}

// String provides string representation of key
func (k Key) String() string {
	return string(k)
}

// With provides a new key that combines another joining them together
// with a `+` symbol.
func (k Key) With(ke Key) Key {
	return Key(fmt.Sprintf("%s+%s", k, ke))
}

// In returns true if the key's value matches one of those
// in the provided list.
func (k Key) In(set ...Key) bool {
	for _, v := range set {
		if v == k {
			return true
		}
	}
	return false
}

// JSONSchema provides a representation of the struct for usage in Schema.
func (Key) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Pattern:     KeyPattern,
		Title:       "Key",
		MinLength:   2,
		MaxLength:   64,
		Description: "Text identifier to be used instead of a code for a more verbose but readable identifier.",
	}
}
