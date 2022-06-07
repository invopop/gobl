package tax

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Code represents a string used to uniquely identify the data we're looking
// at. We use "code" instead of "id", to reenforce the fact that codes should
// be more easily set and used by humans within definitions than IDs or UUIDs.
// Tax codes are standardised so that when validated they must contain between
// 2 and 6 inclusive upper-case letters or numbers.
type Code string

var (
	codeValidationRegexp = regexp.MustCompile(`^[A-Z0-9]{2,6}$`)
)

// Validate ensures that the code complies with the expected rules.
func (c Code) Validate() error {
	return validation.Validate(string(c),
		validation.Length(2, 6),
		validation.Match(codeValidationRegexp),
	)
}

// IsEmpty returns true if no code is specified.
func (c Code) IsEmpty() bool {
	return c == ""
}

// String returns string representation of code.
func (c Code) String() string {
	return string(c)
}

// In returns true if the code's value matches one of those
// in the provided array.
func (c Code) In(ary []Code) bool {
	for _, v := range ary {
		if v == c {
			return true
		}
	}
	return false
}
