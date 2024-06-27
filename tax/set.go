package tax

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/invopop/gobl/cal"
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
	// Percent defines the percentage set manually or determined from the rate
	// key (calculated if rate present). A nil percent implies that this tax combo
	// is **exempt** from tax.
	Percent *num.Percentage `json:"percent,omitempty" jsonschema:"title=Percent" jsonschema_extras:"calculated=true"`
	// Some countries require an additional surcharge (calculated if rate present).
	Surcharge *num.Percentage `json:"surcharge,omitempty" jsonschema:"title=Surcharge" jsonschema_extras:"calculated=true"`
	// Local codes that apply for a given rate or percentage that need to be identified and validated.
	Ext Extensions `json:"ext,omitempty" jsonschema:"title=Ext"`

	// Internal link back to the category object
	category *Category
}

// ValidateWithContext ensures the Combo has the correct details.
func (c *Combo) ValidateWithContext(ctx context.Context) error {
	r, _ := ctx.Value(KeyRegime).(*Regime)
	if r == nil {
		return errors.New("tax regime not found in context")
	}
	cat := r.Category(c.Category)
	rate := r.Rate(c.Category, c.Rate)
	err := validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.Category,
			validation.Required,
			r.InCategories(),
		),
		validation.Field(&c.Rate,
			validation.When(
				(cat != nil && cat.RateRequired),
				validation.Required,
			),
			r.InCategoryRates(c.Category),
		),
		validation.Field(&c.Ext,
			ExtensionsHas(combineExtKeys(cat, rate)...),
			validation.When(
				(cat != nil && len(cat.Extensions) == 0) &&
					(rate != nil && len(rate.Extensions) == 0),
				validation.Empty,
				validation.Skip,
			),
		),
		validation.Field(&c.Percent),
		validation.Field(&c.Surcharge,
			validation.When(
				c.Percent == nil,
				validation.Nil.Error("required with percent"),
			),
		),
	)
	if err != nil {
		return err
	}
	return r.ValidateObject(c)
}

// NormalizeCombo tries to normalize the data inside the tax combo.
func NormalizeCombo(c *Combo) *Combo {
	if c == nil {
		return nil
	}
	c.Ext = NormalizeExtensions(c.Ext)
	return c
}

func combineExtKeys(cat *Category, rate *Rate) []cbc.Key {
	keys := make([]cbc.Key, 0)
	if cat != nil {
		keys = append(keys, cat.Extensions...)
	}
	if rate != nil {
		keys = append(keys, rate.Extensions...)
	}
	return keys
}

// prepare updates the Combo object's Percent and Retained properties using the base totals
// as a source of additional data for making decisions.
func (c *Combo) prepare(r *Regime, tags []cbc.Key, date cal.Date) error {
	c.category = r.Category(c.Category)
	if c.category == nil {
		return ErrInvalidCategory.WithMessage("'%s' not defined in regime", c.Category.String())
	}

	// If there is no rate for the combo, there isn't much else we can do.
	if c.Rate == cbc.KeyEmpty {
		return nil
	}

	rate := c.category.Rate(c.Rate)
	if rate == nil {
		return ErrInvalidRate.WithMessage("'%s' rate not defined in category '%s'", c.Rate.String(), c.Category.String())
	}
	if rate.Exempt {
		c.Percent = nil
		c.Surcharge = nil
		return nil
	}

	// if there are no rate values, don't attempt to prepare anything else.
	if len(rate.Values) == 0 {
		return nil
	}

	value := rate.Value(date, tags, c.Ext)
	if value == nil {
		return ErrInvalidDate.WithMessage("rate value unavailable for '%s' in '%s' on '%s'", c.Rate.String(), c.Category.String(), date.String())
	}

	// 2024-03-14: only update the percentage if none previously set.
	// This means that custom percentages can be used even if the
	// rate classification is required by a regime (like PT).
	if c.Percent == nil {
		p := value.Percent // copy
		c.Percent = &p
	}
	if c.Surcharge == nil {
		if value.Surcharge != nil {
			s := *value.Surcharge // copy
			c.Surcharge = &s
		} else {
			c.Surcharge = nil
		}
	}

	// Run the regime's calculations and normalisations
	return r.CalculateObject(c)
}

// UnmarshalJSON is a temporary migration helper that will move the
// first of the "tags" array used in earlier versions of GOBL into
// the rate field.
func (c *Combo) UnmarshalJSON(data []byte) error {
	type Alias Combo
	aux := struct {
		*Alias
		Tags []cbc.Key `json:"tags"`
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if len(aux.Tags) > 0 && c.Rate == cbc.KeyEmpty {
		c.Rate = aux.Tags[0]
	}
	return nil
}

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
