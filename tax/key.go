package tax

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Key is used to define an ID or code that more closely represents
// a human name.
type Key string

var (
	keyValidationRegexp = regexp.MustCompile(`^[a-z][a-z0-9-]+$`)
)

// Validate ensures the key complies with the basic syntax
// requirements.
func (k Key) Validate() error {
	return validation.Validate(string(k),
		validation.Match(keyValidationRegexp),
	)
}

// String provides string representation of key
func (k Key) String() string {
	return string(k)
}
