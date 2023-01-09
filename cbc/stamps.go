package cbc

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Stamp defines an official seal of approval from a third party like a governmental agency
// or intermediary and should thus be included in any official envelopes.
type Stamp struct {
	// Identity of the agency used to create the stamp usually defined by each region.
	Provider Key `json:"prv" jsonschema:"title=Provider"`
	// The serialized stamp value generated for or by the external agency
	Value string `json:"val" jsonschema:"title=Value"`
}

// Validate checks that the header contains the basic information we need to function.
func (s *Stamp) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Provider, validation.Required),
		validation.Field(&s.Value, validation.Required),
	)
}
