package tax

import (
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
)

// Schemes defines an array of scheme objects with helper functions.
type Schemes []*Scheme

// SchemeKeys stores a list of keys that makes it easier to perform matches.
type SchemeKeys []org.Key

// Scheme contains the definition of a scheme that belongs to a region and can be used
// to simplify validation processes for document contents.
type Scheme struct {
	// Key used to identify this scheme
	Key org.Key `json:"key" jsonschema:"title=Key"`

	// Name of this scheme.
	Name i18n.String `json:"name" jsonschema:"title=Name"`
	// Human details describing what this scheme is used for.
	Description i18n.String `json:"description,omitempty" jsonschema:"title=Description"`

	// List of tax category codes that can be used when this scheme is
	// applied.
	Categories []org.Code `json:"categories,omitempty" jsonschema:"title=Category Codes"`

	// Note defines a message that should be added to a document
	// when this scheme is used.
	Note *org.Note `json:"note,omitempty" jsonschema:"title=Note"`
}

// ForKey finds the scheme with a matching key.
func (ss Schemes) ForKey(key org.Key) *Scheme {
	for _, s := range ss {
		if s.Key == key {
			return s
		}
	}
	return nil
}

// Contains returns true if the list of keys contains a match for the provided
// key.
func (sk SchemeKeys) Contains(key org.Key) bool {
	for _, v := range sk {
		if key == v {
			return true
		}
	}
	return false
}
