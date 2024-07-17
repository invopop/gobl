package cbc

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Key is used to define an ID or code that more closely represents
// a human name. The objective is to make it easier to define constants
// that can be re-used more easily.
type Key string

var (
	// KeyPattern describes what should keys look like
	KeyPattern = `^(?:[a-z]|[a-z0-9][a-z0-9-+]*[a-z0-9])$`
	// KeyValidationRegexp is used for key validation
	KeyValidationRegexp = regexp.MustCompile(KeyPattern)
	// KeySeparator is used to separate keys join using the "With"
	// method.
	KeySeparator = "+"
)

var (
	// KeyMinLength defines the minimum key length
	KeyMinLength uint64 = 1
	// KeyMaxLength defines the maximum key length
	KeyMaxLength uint64 = 64
)

// KeyEmpty is used when no key is available.
const KeyEmpty Key = ""

// Validate ensures the key complies with the basic syntax
// requirements.
func (k Key) Validate() error {
	return validation.Validate(string(k),
		validation.Match(KeyValidationRegexp),
		validation.Length(int(KeyMinLength), int(KeyMaxLength)),
	)
}

// String provides string representation of key
func (k Key) String() string {
	return string(k)
}

// KeyStrings is a convenience method to convert a list of keys
// into a list of strings.
func KeyStrings(keys []Key) []string {
	l := make([]string, len(keys))
	for i, v := range keys {
		l[i] = v.String()
	}
	return l
}

// With provides a new key that combines another joining them together
// with a `+` symbol.
func (k Key) With(ke Key) Key {
	return Key(fmt.Sprintf("%s%s%s", k, KeySeparator, ke))
}

// Has returns true if the key contains the provided key.
func (k Key) Has(ke Key) bool {
	for _, v := range strings.Split(k.String(), KeySeparator) {
		if Key(v) == ke {
			return true
		}
	}
	return false
}

// HasPrefix checks to see if the key starts with the provided key.
// As per `Has`, only the complete key between `+` symbols are
// matched.
func (k Key) HasPrefix(ke Key) bool {
	ks := strings.SplitN(k.String(), KeySeparator, 2)
	if len(ks) == 0 {
		return false
	}
	return ks[0] == ke.String()
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

// IsEmpty returns true if the key has no value.
func (k Key) IsEmpty() bool {
	return k == KeyEmpty
}

// HasValidKeyIn provides a validator to check the Key's
// value is within the provided known set.
func HasValidKeyIn(keys ...Key) validation.Rule {
	return hasKeyRule{elements: keys}
}

type hasKeyRule struct {
	elements []Key
}

func (r hasKeyRule) Validate(v interface{}) error {
	mk, ok := v.(Key)
	if !ok || mk == KeyEmpty {
		return nil
	}
	for _, e := range r.elements {
		if mk.HasPrefix(e) {
			return nil
		}
	}
	return errors.New("must be or start with a valid key")
}

// JSONSchema provides a representation of the struct for usage in Schema.
func (Key) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Pattern:     KeyPattern,
		Title:       "Key",
		MinLength:   &KeyMinLength,
		MaxLength:   &KeyMaxLength,
		Description: "Text identifier to be used instead of a code for a more verbose but readable identifier.",
	}
}
