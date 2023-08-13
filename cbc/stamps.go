package cbc

import (
	"errors"
	"fmt"

	"github.com/invopop/validation"
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

// In checks if the stamp is in the list of stamps.
func (s *Stamp) In(ss []*Stamp) bool {
	for _, r := range ss {
		if s.Provider == r.Provider {
			return true
		}
	}
	return false
}

// DetectDuplicateStamps checks if the list of stamps contains duplicate
// provider keys.
var DetectDuplicateStamps = validation.By(duplicateDuplicateStamps)

func duplicateDuplicateStamps(list interface{}) error {
	values, ok := list.([]*Stamp)
	if !ok {
		return errors.New("must be a stamp array")
	}
	set := []*Stamp{}
	// loop through and check order of Since value
	for _, v := range values {
		if v.In(set) {
			return fmt.Errorf("duplicate stamp '%v'", v.Provider)
		}
		set = append(set, v)
	}
	return nil
}

// AddStamp makes it easier to add a new Stamp by replacing a previous
// entry with a matching Key.
func AddStamp(in []*Stamp, s *Stamp) []*Stamp {
	if in == nil {
		return []*Stamp{s}
	}
	for _, v := range in {
		if v.Provider == s.Provider {
			*v = *s // copy in place
			return in
		}
	}
	return append(in, s)
}
