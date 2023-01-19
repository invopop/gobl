package tax

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
)

// Regime defines the holding structure for the definitions of taxes inside a country
// or territory.
type Regime struct {
	// Name of the country
	Name i18n.String `json:"name" jsonschema:"title=Name"`

	// Country code for the region
	Country l10n.CountryCode `json:"country" jsonschema:"title=Code"`

	// Specific Locality, region, city, province, county, or similar code inside
	// the country, if needed.
	Zone l10n.Code `json:"zone,omitempty" jsonschema:"title=Zone"`

	// List of sub-zones inside a country.
	Zones []Zone `json:"zones,omitempty" jsonschema:"title=Zones"`

	// Currency used by the country.
	Currency currency.Code `json:"currency" jsonschema:"title=Currency"`

	// Set of specific scheme definitions inside the region.
	Schemes Schemes `json:"schemes,omitempty" jsonschema:"title=Schemes"`

	// List of tax categories.
	Categories []*Category `json:"categories" jsonschema:"title=Categories"`

	// Validator is a method to use to validate a document in a given region.
	Validator func(doc interface{}) error `json:"-"`

	// Calculator is used to performs regime specific calculations on data,
	// including any normalization that might need to take place such as
	// with tax codes and removing white-space.
	Calculator func(doc interface{}) error `json:"-"`
}

// Zone represents an area inside a country, like a province
// or a state, which shares the basic definitions of the country, but
// may vary in some validation rules.
type Zone struct {
	// Unique zone code.
	Code l10n.Code `json:"code" jsonschema:"title=Code"`
	// Name of the zone to be use if a locality or region is not applicable.
	Name i18n.String `json:"name,omitempty" jsonschema:"title=Name"`
	// Village, town, district, or city name which should coincide with
	// address data.
	Locality i18n.String `json:"locality,omitempty" jsonschema:"title=Locality"`
	// Province, county, or state which should match address data.
	Region i18n.String `json:"region,omitempty" jsonschema:"title=Region"`
	// Any additional information
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Category contains the definition of a general type of tax inside a region.
type Category struct {
	Code cbc.Code    `json:"code" jsonschema:"title=Code"`
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
	Key cbc.Key `json:"key" jsonschema:"title=Key"`

	Name i18n.String `json:"name" jsonschema:"title=Name"`
	Desc i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`

	// Values contains a list of Value objects that contain the
	// current and historical percentage values for the rate and
	// additional filters.
	// Order is important, newer values should come before
	// older values.
	Values []*RateValue `json:"values" jsonschema:"title=Values"`
}

// RateValue contains a percentage rate or fixed amount for a given date range.
// Fiscal policy changes mean that rates are not static so we need to
// be able to apply the correct rate for a given period.
type RateValue struct {
	// Only use this value if one of the zones matches.
	Zones []l10n.Code `json:"zones,omitempty" jsonschema:"title=Zones"`
	// Date from which this value should be applied.
	Since *cal.Date `json:"since,omitempty" jsonschema:"title=Since"`
	// Percent rate that should be applied
	Percent num.Percentage `json:"percent" jsonschema:"title=Percent"`
	// An additional surcharge to apply.
	Surcharge *num.Percentage `json:"surcharge,omitempty" jsonschema:"title=Surcharge"`
	// When true, this value should no longer be used.
	Disabled bool `json:"disabled,omitempty" jsonschema:"title=Disabled"`
}

// ValidateDocument performs validation on the provided document.
func (r *Regime) ValidateDocument(obj interface{}) error {
	if r.Validator != nil {
		return r.Validator(obj)
	}
	return nil
}

// CalculateDocument performs any region specific calculations on the provided
// object.
func (r *Regime) CalculateDocument(obj interface{}) error {
	if r.Calculator != nil {
		return r.Calculator(obj)
	}
	return nil
}

// CurrencyDef provides the currency definition object for the region.
func (r *Regime) CurrencyDef() *currency.Def {
	d, ok := currency.Get(r.Currency)
	if !ok {
		return nil
	}
	return &d
}

// Validate enures the region definition is valid, including all
// subsequent categories.
func (r *Regime) Validate() error {
	err := validation.ValidateStruct(r,
		validation.Field(&r.Country, validation.Required),
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Categories, validation.Required),
		validation.Field(&r.Zones),
	)
	return err
}

// Validate ensures that the zone looks correct.
func (z *Zone) Validate() error {
	err := validation.ValidateStruct(z,
		validation.Field(&z.Code, validation.Required),
		validation.Field(&z.Name),
		validation.Field(&z.Locality),
		validation.Field(&z.Region),
		validation.Field(&z.Meta),
	)
	return err
}

// Validate ensures the Category's contents are correct.
func (c *Category) Validate() error {
	err := validation.ValidateStruct(c,
		validation.Field(&c.Code, validation.Required),
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Rates),
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
func (rv *RateValue) Validate() error {
	return validation.ValidateStruct(rv,
		validation.Field(&rv.Percent, validation.Required),
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
func (r *Regime) Category(code cbc.Code) *Category {
	for _, c := range r.Categories {
		if c.Code == code {
			return c
		}
	}
	return nil
}

// Rate provides the rate definition with a matching key for
// the category.
func (c *Category) Rate(key cbc.Key) *Rate {
	for _, r := range c.Rates {
		if r.Key == key {
			return r
		}
	}
	return nil
}

// Value determines the tax rate value for the provided date and zone, if applicable.
func (r *Rate) Value(date cal.Date, zone l10n.Code) *RateValue {
	for _, rv := range r.Values {
		if len(rv.Zones) > 0 {
			if !rv.HasZone(zone) {
				continue
			}
		}
		if rv.Since == nil || !rv.Since.IsValid() || rv.Since.Before(date.Date) {
			return rv
		}
	}
	return nil
}

// HasZone returns true if the rate value has a zone that matches the one provided.
func (rv *RateValue) HasZone(zone l10n.Code) bool {
	for _, z := range rv.Zones {
		if z == zone {
			return true
		}
	}
	return false
}
