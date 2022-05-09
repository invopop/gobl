package tax

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Set defines a list of tax categories and their rates to be used alongside taxable items.
type Set []*Combo

// Combo represents the tax combination of a category code and rate key.
type Combo struct {
	// Tax category code from those available inside a region.
	Category Code `json:"cat" jsonschema:"title=Category"`
	// Rate within a category to apply.
	Rate Key `json:"rate" jsonschema:"title=Rate"`

	// Objects used internally for making calculations on specific dates
	// see the Region#prepareCombo method for usage.
	category *Category
	rate     *Rate
	value    *RateValue
}

// Validate ensures the Combo contains all the details required.
func (c *Combo) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Category, validation.Required),
		validation.Field(&c.Rate, validation.Required),
	)
}

// Validate ensures the set of tax combos looks correct
func (s Set) Validate() error {
	combos := make(map[Code]Key)
	for i, c := range s {
		if _, ok := combos[c.Category]; ok {
			return fmt.Errorf("%d: category %v is duplicated", i, c.Category)
		}
		if err := c.Validate(); err != nil {
			return fmt.Errorf("%d: %w", i, err)
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
			if a.Category == b.Category && a.Rate == b.Rate {
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
func (s Set) Get(cat Code) *Combo {
	for _, c := range s {
		if c.Category == cat {
			return c
		}
	}
	return nil
}

// Rate returns the rate from the matching category, if set.
func (s Set) Rate(cat Code) Key {
	for _, c := range s {
		if c.Category == cat {
			return c.Rate
		}
	}
	return ""
}
