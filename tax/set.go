package tax

import (
	"context"
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/validation"
)

// Set defines a list of tax categories and their rates to be used alongside taxable items.
type Set []*Combo

// NormalizeSet tries to normalize the tax set by normalizing combos
// and returning nil if empty.
func NormalizeSet(s Set) Set {
	if s == nil {
		return nil
	}
	ns := make(Set, 0)
	for _, c := range s {
		c = NormalizeCombo(c)
		if c == nil {
			continue
		}
		ns = append(ns, c)
	}
	if len(ns) == 0 {
		return nil
	}
	return ns
}

// ValidateWithContext ensures the set of tax combos looks correct
func (s Set) ValidateWithContext(ctx context.Context) error {
	combos := make(map[cbc.Code]cbc.Key)
	for i, c := range s {
		if _, ok := combos[c.Category]; ok {
			return validation.Errors{
				fmt.Sprintf("%d", i): fmt.Errorf("category %v is duplicated", c.Category),
			}
		}
		if err := c.ValidateWithContext(ctx); err != nil {
			return validation.Errors{
				fmt.Sprintf("%d", i): err,
			}
		}
		combos[c.Category] = c.Rate
	}
	return nil
}

// Equals returns true if the sets match, regardless of order.
func (s Set) Equals(s2 Set) bool {
	for _, a := range s {
		match := false
		for _, b := range s2 {
			if a.Category == b.Category && a.Rate == b.Rate && a.Country == b.Country {
				match = true
			}
		}
		if !match {
			// implies the code defined in the base, was not present in the second
			// array.
			return false
		}
	}
	return true
}

// Get the Rate key for the given category
func (s Set) Get(cat cbc.Code) *Combo {
	for _, c := range s {
		if c.Category == cat {
			return c
		}
	}
	return nil
}

// Rate returns the rate from the matching category, if set.
func (s Set) Rate(cat cbc.Code) cbc.Key {
	for _, c := range s {
		if c.Category == cat {
			return c.Rate
		}
	}
	return ""
}

type setValidation struct {
	categories []cbc.Code
}

// SetHasCategory validates that the set contains the given category.
func SetHasCategory(categories ...cbc.Code) validation.Rule {
	return &setValidation{categories: categories}
}

func (sv *setValidation) Validate(value interface{}) error {
	s, ok := value.(Set)
	if !ok {
		return nil
	}
	for _, c := range sv.categories {
		if s.Get(c) == nil {
			return fmt.Errorf("missing category %s", c.String())
		}
	}
	return nil
}
