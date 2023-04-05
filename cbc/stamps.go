package cbc

import (
	"github.com/invopop/gobl/num"
	"github.com/invopop/validation"
)

// Stamp defines an official seal of approval from a third party like a governmental agency
// or intermediary and should thus be included in any official envelopes.
type Stamp struct {
	// Identity of the agency used to create the stamp usually defined by each region.
	Provider Key `json:"prv" jsonschema:"title=Provider"`
	// The serialized stamp value generated for or by the external agency
	Value string `json:"val" jsonschema:"title=Value"`
	// The fee incurred by the stamp
	Fee num.Amount `json:"fee,omitempty" jsonschema:"title=Fee"`
	// Any additional semi-structured information
	Meta Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate checks that the header contains the basic information we need to function.
func (s *Stamp) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Provider, validation.Required),
		validation.Field(&s.Value, validation.Required),
	)
}
