package tax

import (
	"context"
	"errors"
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/validation"
)

// Set defines a list of tax categories and their rates to be used alongside taxable items.
type Set []*Combo

// Combo represents the tax combination of a category code and rate key. The percent
// and retained attributes will be determined automatically from the Rate key if set
// during calculation.
type Combo struct {
	// Tax category code from those available inside a region.
	Category cbc.Code `json:"cat" jsonschema:"title=Category"`
	// Rate within a category to apply.
	Rate cbc.Key `json:"rate,omitempty" jsonschema:"title=Rate"`
	// Percent defines the percentage set manually or determined from the rate key (calculated if rate present).
	Percent num.Percentage `json:"percent" jsonschema:"title=Percent" jsonschema_extras:"calculated=true"`
	// Some countries require an additional surcharge (calculated if rate present).
	Surcharge *num.Percentage `json:"surcharge,omitempty" jsonschema:"title=Surcharge" jsonschema_extras:"calculated=true"`
	// Additional data may be required in some regimes, the tags
	// property helps reference them.
	Tags []cbc.Key `json:"tags,omitempty" jsonschema:"title=Tags"`
	// Internal link back to the category object
	category *Category
}

// ValidateWithContext ensures the Combo has the correct details.
func (c *Combo) ValidateWithContext(ctx context.Context) error {
	r, _ := ctx.Value(KeyRegime).(*Regime)
	if r == nil {
		return errors.New("tax regime not found in context")
	}
	return validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.Category, validation.Required, r.InCategories()),
		validation.Field(&c.Rate), // optional, but should be checked if present
		validation.Field(&c.Percent, validation.Required),
		validation.Field(&c.Surcharge), // not required, but should be valid number
		validation.Field(&c.Tags, validation.Each(r.InCategoryTags(c.Category))),
	)
}

// ContainsTag returns true if the tax combo contains the given tag.
func (c *Combo) ContainsTag(key cbc.Key) bool {
	if c == nil {
		return false
	}
	return key.In(c.Tags...)
}

// prepare updates the Combo object's Percent and Retained properties using the base totals
// as a source of additional data for making decisions.
func (c *Combo) prepare(tc *TotalCalculator) error {
	c.category = tc.Regime.Category(c.Category)
	if c.category == nil {
		return ErrInvalidCategory.WithMessage("'%s' not defined in regime", c.Category.String())
	}

	if c.Rate != cbc.KeyEmpty {
		rate := c.category.Rate(c.Rate)
		if rate == nil {
			return ErrInvalidRate.WithMessage("'%s' rate not defined in category '%s'", c.Rate.String(), c.Category.String())
		}
		value := rate.Value(tc.Date, tc.Zone)
		if value == nil {
			return ErrInvalidDate.WithMessage("rate value unavailable for '%s' in '%s' on '%s'", c.Rate.String(), c.Category.String(), tc.Date.String())
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

// ValidateWithContext ensures the set of tax combos looks correct
func (s Set) ValidateWithContext(ctx context.Context) error {
	combos := make(map[cbc.Code]cbc.Key)
	for i, c := range s {
		if _, ok := combos[c.Category]; ok {
			return fmt.Errorf("%d: category %v is duplicated", i, c.Category)
		}
		if err := c.ValidateWithContext(ctx); err != nil {
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
