package tax

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
)

// Set defines a list of tax categories and their rates to be used alongside taxable items.
type Set []*Combo

// Combo represents the tax combination of a category code and rate key. The percent
// and retained attributes will be determined automatically from the Rate key if set
// during calculation.
type Combo struct {
	// Tax category code from those available inside a region.
	Category org.Code `json:"cat" jsonschema:"title=Category"`
	// Rate within a category to apply.
	Rate org.Key `json:"rate,omitempty" jsonschema:"title=Rate"`
	// Percent defines the percentage set manually or determined from the rate key (calculated if rate present).
	Percent num.Percentage `json:"percent" jsonschema:"title=Percent" jsonschema_extras:"calculated=true"`
	// Some countries require an additional surcharge (calculated if rate present).
	Surcharge *num.Percentage `json:"surcharge,omitempty" jsonschema:"title=Surcharge" jsonschema_extras:"calculated=true"`
	// Internal link back to the category object
	category *Category
}

// Validate ensures the Combo has the correct details.
func (c *Combo) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Category, validation.Required),
		validation.Field(&c.Rate), // optional, but should be checked if present
		validation.Field(&c.Percent, validation.Required),
		validation.Field(&c.Surcharge), // not required, but should be valid number
	)
}

// prepare updates the Combo object's Percent and Retained properties according
// to the region and date provided.
func (c *Combo) prepare(r *Region, date cal.Date) error {
	c.category = r.Category(c.Category)
	if c.category == nil {
		return ErrInvalidCategory.WithMessage("'%s' not in region", c.Category.String())
	}

	if c.Rate != org.KeyEmpty {
		rate := c.category.Rate(c.Rate)
		if rate == nil {
			return ErrInvalidRate.WithMessage("'%s' not in category '%s'", c.Rate.String(), c.Category.String())
		}
		value := rate.On(date)
		if value == nil {
			return ErrInvalidDate.WithMessage("data unavailable for '%s' in '%s' on '%s'", c.Rate.String(), c.Category.String(), date.String())
		}
		c.Percent = value.Percent
		if value.Surcharge != nil {
			s := *value.Surcharge // copy
			c.Surcharge = &s
		} else {
			c.Surcharge = nil
		}
	}

	return nil
}

// Validate ensures the set of tax combos looks correct
func (s Set) Validate() error {
	combos := make(map[org.Code]org.Key)
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
func (s Set) Get(cat org.Code) *Combo {
	for _, c := range s {
		if c.Category == cat {
			return c
		}
	}
	return nil
}

// Rate returns the rate from the matching category, if set.
func (s Set) Rate(cat org.Code) org.Key {
	for _, c := range s {
		if c.Category == cat {
			return c.Rate
		}
	}
	return ""
}
