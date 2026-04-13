package head

import (
	"errors"
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// Stamp defines an official seal of approval from a third party like a governmental agency
// or intermediary and should thus be included in any official envelopes.
type Stamp struct {
	// Identity of the agency used to create the stamp usually defined by each region.
	Provider cbc.Key `json:"prv" jsonschema:"title=Provider"`
	// The serialized stamp value generated for or by the external agency
	Value string `json:"val" jsonschema:"title=Value"`
}

func stampRules() *rules.Set {
	return rules.For(new(Stamp),
		rules.Field("prv",
			rules.Assert("01", "stamp provider is required", is.Present),
		),
		rules.Field("val",
			rules.Assert("02", "stamp value is required", is.Present),
		),
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
var DetectDuplicateStamps = is.FuncError("no duplicate stamps", detectDuplicateStamps)

func detectDuplicateStamps(list any) error {
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

// GetStamp provides the matching stamp for the given provider.
func GetStamp(in []*Stamp, provider cbc.Key) *Stamp {
	if in == nil {
		return nil
	}
	for _, v := range in {
		if v.Provider == provider {
			return v
		}
	}
	return nil
}

// NormalizeStamps will try to clean the stamps by removing rows with empty
// providers or values. If empty, the function will return nil.
func NormalizeStamps(in []*Stamp) []*Stamp {
	if in == nil {
		return nil
	}
	out := make([]*Stamp, 0)
	for _, v := range in {
		if v.Value == "" || v.Provider == "" {
			continue
		}
		out = append(out, v)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// StampsHasRule is a combined rules.Test and validation.Rule that checks if
// a specific stamp provider is present in a list of stamps.
type StampsHasRule struct {
	desc     string
	provider cbc.Key
}

// StampsHas provides a validation rule that checks if a list of stamps includes
// one for the given provider.
func StampsHas(provider cbc.Key) *StampsHasRule {
	return &StampsHasRule{desc: fmt.Sprintf("stamps have %s", provider), provider: provider}
}

func (r *StampsHasRule) String() string {
	return r.desc
}

// Check returns true if the stamp with the given provider is present.
func (r *StampsHasRule) Check(value any) bool {
	in, ok := value.([]*Stamp)
	if !ok {
		return false // invalid type
	}
	return GetStamp(in, r.provider) != nil
}

// Validate returns an error if the stamp with the given provider is absent.
func (r *StampsHasRule) Validate(value any) error {
	if !r.Check(value) {
		return fmt.Errorf("missing %s stamp", r.provider)
	}
	return nil
}
