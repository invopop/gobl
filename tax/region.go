package tax

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
)

// Region defines the holding structure for a regions categories and subsequent
// Rates and Values.
type Region struct {
	// Name of the region
	Name i18n.String `json:"name" jsonschema:"title=Name"`

	// Country code for the region
	Country l10n.Code `json:"country" jsonschema:"title=Code"`
	// Locality, city, region, or similar code inside the country, if needed.
	Locality l10n.Code `json:"locality,omitempty" jsonschema:"title=Locality"`

	// Currency used by the region for tax purposes.
	Currency currency.Code `json:"currency" jsonschema:"title=Currency"`

	// Set of specific scheme definitions inside the region.
	Schemes Schemes `json:"schemes,omitempty" jsonschema:"title=Schemes"`

	// List of tax categories.
	Categories []*Category `json:"categories" jsonschema:"title=Categories"`

	// ValidateDocument is a method to use to validate a document in a given region.
	ValidateDocument func(doc interface{}) error `json:"-"`
}

// Category contains the definition of a general type of tax inside a region.
type Category struct {
	Code Code        `json:"code" jsonschema:"title=Code"`
	Name i18n.String `json:"name" jsonschema:"title=Name"`
	Desc i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`

	// Retained when true implies that the tax amount will be retained
	// by the buyer on behalf of the supplier, and thus subtracted from
	// the invoice taxable base total. Typically used for taxes related to
	// income.
	Retained bool `json:"retained,omitempty" jsonschema:"title=Retained"`

	// Specific tax definitions inside this category.
	Rates []*Rate `json:"rates" jsonschema:"title=Rates"`
}

// Rate defines a single rate inside a category
type Rate struct {
	// Key identifies this rate within the system
	Key Key `json:"key" jsonschema:"title=Key"`

	Name i18n.String `json:"name" jsonschema:"title=Name"`
	Desc i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`

	// Values contains a list of Value objects that contain the
	// current and historical percentage values for the rate;
	// order is important, newer values should come before
	// older values.
	Values []*RateValue `json:"values" jsonschema:"title=Values"`
}

// RateValue contains a percentage rate or fixed amount for a given date range.
// Fiscal policy changes mean that rates are not static so we need to
// be able to apply the correct rate for a given period.
type RateValue struct {
	// Date from which this value should be applied.
	Since *cal.Date `json:"since,omitempty" jsonschema:"title=Since"`
	// Rate that should be applied
	Percent num.Percentage `json:"percent" jsonschema:"title=Percent"`
	// When true, this value should no longer be used.
	Disabled bool `json:"disabled,omitempty" jsonschema:"title=Disabled"`
}

// CurrencyDef provides the currency definition object for the region.
func (r *Region) CurrencyDef() *currency.Def {
	d, ok := currency.Get(r.Currency)
	if !ok {
		return nil
	}
	return &d
}

// Validate enures the region definition is valid, including all
// subsequent categories.
func (r *Region) Validate() error {
	err := validation.ValidateStruct(r,
		validation.Field(&r.Country, validation.Required),
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Categories, validation.Required),
	)
	return err
}

// Validate ensures the Category's contents are correct.
func (c *Category) Validate() error {
	err := validation.ValidateStruct(c,
		validation.Field(&c.Code, validation.Required),
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Rates, validation.Required),
	)
	return err
}

// Validate checks that our tax definition is valid. This is only really
// meant to be used when testing new regional tax definitions.
func (r *Rate) Validate() error {
	err := validation.ValidateStruct(r,
		validation.Field(&r.Key, validation.Required),
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Values, validation.Required, validation.By(checkRateValuesOrder)),
	)
	return err
}

// Validate ensures the tax rate contains all the required fields.
func (v *RateValue) Validate() error {
	return validation.ValidateStruct(v,
		validation.Field(&v.Percent, validation.Required),
	)
}

func checkRateValuesOrder(list interface{}) error {
	values, ok := list.([]*RateValue)
	if !ok {
		return errors.New("must be a tax rate value array")
	}
	var date *cal.Date
	// loop through and check order of Since value
	for i := range values {
		v := values[i]
		if date != nil && date.IsValid() {
			if v.Since.IsValid() && !v.Since.Before(date.Date) {
				return errors.New("invalid date order")
			}
		}
		date = v.Since
	}
	return nil
}

// Category provides the requested category by its code.
func (r *Region) Category(code Code) *Category {
	for _, c := range r.Categories {
		if c.Code == code {
			return c
		}
	}
	return nil
}

// Rate provides the rate definition with a matching key for
// the category.
func (c *Category) Rate(key Key) *Rate {
	for _, r := range c.Rates {
		if r.Key == key {
			return r
		}
	}
	return nil
}

// On determines the tax rate value for the provided date.
func (r *Rate) On(date cal.Date) *RateValue {
	for _, v := range r.Values {
		if v.Since == nil || !v.Since.IsValid() || v.Since.Before(date.Date) {
			return v
		}
	}
	return nil
}

// prepareCombo updates the Combo object's internal properties to include the objects
// for the region on a given date.
func (r *Region) prepareCombo(c *Combo, date cal.Date) error {
	c.category = r.Category(c.Category)
	if c.category == nil {
		return ErrInvalidCategory.WithMessage("'%s'", c.Category.String())
	}
	c.rate = c.category.Rate(c.Rate)
	if c.rate == nil {
		return ErrInvalidRate.WithMessage("'%s' in category '%s'", c.Rate.String(), c.Category.String())
	}
	c.value = c.rate.On(date)
	if c.value == nil {
		return ErrInvalidDate.WithMessage("data unavailable for '%s' in '%s' on '%s'", c.Rate.String(), c.Category.String(), date.String())
	}
	return nil
}
