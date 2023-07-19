package tax

import (
	"context"
	"errors"
	"strings"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/validation"
)

const (
	// KeyRegime is used in the context to store the tax regime during validation.
	KeyRegime cbc.Key = "tax-regime"
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

	// Tags that can be applied at the document level to identify additional
	// considerations.
	Tags []*KeyDefinition `json:"tags,omitempty" jsonschema:"title=Tags"`

	// Identity types specific for the regime and may be validated
	// against.
	IdentityTypeKeys []*KeyDefinition `json:"identity_types,omitempty" jsonschema:"title=Identity Types"`

	// Charge types specific for the regime and may be validated or used in the UI as suggestions
	ChargeKeys []*KeyDefinition `json:"charge_types,omitempty" jsonschema:"title=Charge Types"`

	// PaymentMeansKeys specific for the regime that extend the original
	// base payment means keys.
	PaymentMeansKeys []*KeyDefinition `json:"payment_means,omitempty" jsonschema:"title=Payment Means"`

	// ItemKeys specific for the regime that need to be added to `org.Item` data
	// in line rows.
	ItemKeys []*KeyDefinition `json:"item_keys,omitempty" jsonschema:"title=Item Keys"`

	// Sets of scenario definitions for the regime.
	Scenarios []*ScenarioSet `json:"scenarios,omitempty" jsonschema:"title=Scenarios"`

	// Configuration details for preceding options.
	Preceding *PrecedingDefinitions `json:"preceding,omitempty" jsonschema:"title=Preceding"`

	// List of tax categories.
	Categories []*Category `json:"categories" jsonschema:"title=Categories"`

	// Validator is a method to use to validate a document in a given region.
	Validator func(doc interface{}) error `json:"-"`

	// Calculator is used to performs regime specific calculations on data,
	// including any normalization that might need to take place such as
	// with tax codes and removing white-space.
	Calculator func(ctx context.Context, doc interface{}) error `json:"-"`
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
	// Codes defines a set of regime specific code mappings.
	Codes cbc.CodeSet `json:"codes,omitempty" jsonschema:"title=Codes"`
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

	// RateRequired when true implies that when a tax combo is defined using
	// this category that one of the rates must be defined.
	RateRequired bool `json:"rate_required,omitempty" jsonschema:"title=Rate Required"`

	// Specific tax definitions inside this category.
	Rates []*Rate `json:"rates,omitempty" jsonschema:"title=Rates"`

	// Codes defines a set of regime specific code mappings.
	Codes cbc.CodeSet `json:"codes,omitempty" jsonschema:"title=Codes"`

	// Meta contains additional information about the category that is relevant
	// for local frequently used formats.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Rate defines a single rate inside a category
type Rate struct {
	// Key identifies this rate within the system
	Key cbc.Key `json:"key" jsonschema:"title=Key"`

	// Human name of the rate
	Name i18n.String `json:"name" jsonschema:"title=Name"`
	// Useful description of the rate.
	Desc i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`

	// Exempt when true implies that the rate when used in a tax Combo should
	// not define a percent value.
	Exempt bool `json:"exempt,omitempty" jsonschema:"title=Exempt"`

	// Values contains a list of Value objects that contain the
	// current and historical percentage values for the rate and
	// additional filters.
	// Order is important, newer values should come before
	// older values.
	Values []*RateValue `json:"values,omitempty" jsonschema:"title=Values"`

	// Codes defines a set of regime specific code mappings.
	Codes cbc.CodeSet `json:"codes,omitempty" jsonschema:"title=Codes"`

	// Meta contains additional information about the rate that is relevant
	// for local frequently used implementations.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
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

// PrecedingDefinitions contains details about what can be defined in Invoice
// preceding document data.
type PrecedingDefinitions struct {
	// The types of sub-documents supported by the regime
	Types []cbc.Key `json:"types,omitempty" jsonschema:"title=Types"`
	// Stamps that must be copied from the preceding document.
	Stamps []cbc.Key `json:"stamps,omitempty" jsonschema:"title=Stamps"`
	// Corrections contains a list of all the keys that can be used to identify a correction.
	Corrections []*KeyDefinition `json:"corrections,omitempty" jsonschema:"title=Corrections"`
	// CorrectionMethods describe the methods used to correct an invoice.
	CorrectionMethods []*KeyDefinition `json:"correction_methods,omitempty" jsonschema:"title=Correction Methods"`
}

// KeyDefinition defines properties of a key that is specific for a regime.
type KeyDefinition struct {
	// Actual key value.
	Key cbc.Key `json:"key" jsonschema:"title=Key"`
	// Short name for the key, if relevant.
	Name i18n.String `json:"name,omitempty" jsonschema:"title=Name"`
	// Description offering more details about when the key should be used.
	Desc i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`
	// Codes defines a set of codes that may be used within a given regime or format
	// that will help identify which key definition to use.
	Codes cbc.CodeSet `json:"codes,omitempty" jsonschema:"title=Codes"`
	// Any additional data that might be relevant in some regimes.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// ValidateObject performs validation on the provided object in the context
// of the regime.
func (r *Regime) ValidateObject(obj interface{}) error {
	if r.Validator != nil {
		return r.Validator(obj)
	}
	return nil
}

// CalculateObject performs any regime specific calculations on the provided
// object.
func (r *Regime) CalculateObject(ctx context.Context, obj interface{}) error {
	if r.Calculator != nil {
		return r.Calculator(ctx, obj)
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

// ScenarioSet returns a single scenario set instance for the provided
// document schema.
func (r *Regime) ScenarioSet(schema string) *ScenarioSet {
	for _, s := range r.Scenarios {
		if strings.HasSuffix(schema, s.Schema) {
			return s
		}
	}
	return nil
}

// Validate enures the region definition is valid, including all
// subsequent categories.
func (r *Regime) Validate() error {
	err := validation.ValidateStruct(r,
		validation.Field(&r.Country, validation.Required),
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Scenarios),
		validation.Field(&r.Categories, validation.Required),
		validation.Field(&r.Zones),
	)
	return err
}

// InTags returns a validation rule to ensure the tag key
// is inside the list of known tags.
func (r *Regime) InTags() validation.Rule {
	if r == nil {
		return validation.In()
	}
	tags := make([]interface{}, len(r.Tags))
	for i, t := range r.Tags {
		tags[i] = t.Key
	}
	return validation.In(tags...)
}

// InCategoryRates is used to provide a validation rule that will
// ensure a rate key is defined inside a category.
func (r *Regime) InCategoryRates(cat cbc.Code) validation.Rule {
	if r == nil {
		return validation.In()
	}
	c := r.Category(cat)
	if c == nil {
		return validation.In()
	}
	keys := make([]interface{}, len(c.Rates))
	for i, x := range c.Rates {
		keys[i] = x.Key
	}
	return validation.In(keys...)
}

// InCategories returns a validation rule to ensure the category code
// is inside the list of known codes.
func (r *Regime) InCategories() validation.Rule {
	if r == nil {
		return validation.In()
	}
	cats := make([]interface{}, len(r.Categories))
	for i, c := range r.Categories {
		cats[i] = c.Code
	}
	return validation.In(cats...)
}

// WithContext adds this regime to the given context.
func (r *Regime) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, KeyRegime, r)
}

// RegimeFromContext returns the regime from the given context, or nil.
func RegimeFromContext(ctx context.Context) *Regime {
	r, ok := ctx.Value(KeyRegime).(*Regime)
	if !ok {
		return nil
	}
	return r
}

// ValidateInRegime ensures that the object is valid in the context of the
// regime.
func ValidateInRegime(ctx context.Context, obj interface{}) error {
	r := RegimeFromContext(ctx)
	if r == nil {
		return nil
	}
	return r.ValidateObject(obj)
}

// ValidateStructWithRegime wraps around the standard validation.ValidateStructWithContext
// method to add an additional check for the tax regime.
func ValidateStructWithRegime(ctx context.Context, obj interface{}, fields ...*validation.FieldRules) error {
	// First run regular validation
	if err := validation.ValidateStructWithContext(ctx, obj, fields...); err != nil {
		return err
	}
	return ValidateInRegime(ctx, obj)
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
		validation.Field(&r.Values,
			validation.When(r.Exempt, validation.Nil),
			validation.By(checkRateValuesOrder),
		),
		validation.Field(&r.Meta),
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

// Tag returns the KeyDefinition for the provided tag key
func (r *Regime) Tag(key cbc.Key) *KeyDefinition {
	for _, t := range r.Tags {
		if t.Key == key {
			return t
		}
	}
	return nil
}

// Rate provides the rate definition for the provided category code
// and rate key.
func (r *Regime) Rate(cat cbc.Code, key cbc.Key) *Rate {
	c := r.Category(cat)
	if c == nil {
		return nil
	}
	return c.Rate(key)
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

// HasType returns true if the preceding definitions has a type that matches the one provided.
func (pd *PrecedingDefinitions) HasType(t cbc.Key) bool {
	if pd == nil {
		return false // no preceding definitions
	}
	return t.In(pd.Types...)
}

// Validate ensures the key definition looks correct in the context of the regime.
func (kd *KeyDefinition) Validate() error {
	err := validation.ValidateStruct(kd,
		validation.Field(&kd.Key, validation.Required),
		validation.Field(&kd.Name, validation.Required),
		validation.Field(&kd.Desc),
		validation.Field(&kd.Codes),
		validation.Field(&kd.Meta),
	)
	return err
}
