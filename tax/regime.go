package tax

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
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

	// Introductory details about the regime.
	Description i18n.String `json:"description,omitempty" jsonschema:"title=Description"`

	// Location name for the country's central time zone. Accepted
	// values from IANA Time Zone Database (https://iana.org/time-zones).
	TimeZone string `json:"time_zone" jsonschema:"title=Time Zone"`

	// Country code for the region
	Country l10n.CountryCode `json:"country" jsonschema:"title=Code"`

	// Specific Locality, region, city, province, county, or similar code inside
	// the country, if needed.
	Zone l10n.Code `json:"zone,omitempty" jsonschema:"title=Zone"`

	// Currency used by the country.
	Currency currency.Code `json:"currency" jsonschema:"title=Currency"`

	// Tags that can be applied at the document level to identify additional
	// considerations.
	Tags []*cbc.KeyDefinition `json:"tags,omitempty" jsonschema:"title=Tags"`

	// Identity types specific for the regime and may be validated
	// against.
	IdentityTypeKeys []*cbc.KeyDefinition `json:"identity_types,omitempty" jsonschema:"title=Identity Types"`

	// Extensions defines the keys that can be used for extended or extra data inside the regime that
	// is specific to the regime and cannot be easily determined from other GOBL structures.
	// Typically these are used to define local codes for suppliers, customers, products, or tax rates.
	Extensions []*cbc.KeyDefinition `json:"extensions,omitempty" jsonschema:"title=Extensions"`

	// Charge types specific for the regime and may be validated or used in the UI as suggestions
	ChargeKeys []*cbc.KeyDefinition `json:"charge_types,omitempty" jsonschema:"title=Charge Types"`

	// PaymentMeansKeys specific for the regime that extend the original
	// base payment means keys.
	PaymentMeansKeys []*cbc.KeyDefinition `json:"payment_means,omitempty" jsonschema:"title=Payment Means"`

	// InboxKeys specific to the regime that can be used to identify where a document
	// should be forwarded to.
	InboxKeys []*cbc.KeyDefinition `json:"inbox_keys,omitempty" jsonschema:"title=Inbox Keys"`

	Scenarios []*ScenarioSet `json:"scenarios,omitempty" jsonschema:"title=Scenarios"`

	// Configuration details for corrections to be used with correction options.
	Corrections []*CorrectionDefinition `json:"corrections,omitempty" jsonschema:"title=Corrections"`

	// List of tax categories.
	Categories []*Category `json:"categories" jsonschema:"title=Categories"`

	// Validator is a method to use to validate a document in a given region.
	Validator func(doc interface{}) error `json:"-"`

	// Calculator is used to performs regime specific calculations on data,
	// including any normalization that might need to take place such as
	// with tax codes and removing white-space.
	Calculator func(doc interface{}) error `json:"-"`
}

// Category contains the definition of a general type of tax inside a region.
type Category struct {
	// Code to be used in documents
	Code cbc.Code `json:"code" jsonschema:"title=Code"`

	// Short name of the category to be used instead of code in output
	Name i18n.String `json:"name" jsonschema:"title=Name"`

	// Human name for the code to use for titles
	Title i18n.String `json:"title,omitempty" jsonschema:"title=Title"`

	// Useful description of the category.
	Description *i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`

	// Retained when true implies that the tax amount will be retained
	// by the buyer on behalf of the supplier, and thus subtracted from
	// the invoice taxable base total. Typically used for taxes related to
	// income.
	Retained bool `json:"retained,omitempty" jsonschema:"title=Retained"`

	// RateRequired when true implies that when a tax combo is defined using
	// this category that one of the rate's keys must be defined. This is
	// normally needed for regimes that categorize taxes in local document
	// formats as opposed to grouping by percentage values.
	RateRequired bool `json:"rate_required,omitempty" jsonschema:"title=Rate Required"`

	// Specific tax definitions inside this category.
	Rates []*Rate `json:"rates,omitempty" jsonschema:"title=Rates"`

	// Extensions defines a list of extension keys that may be used or required
	// as an alternative or alongside choosing a rate for the tax category.
	// Every key must be defined in the Regime's extensions table.
	Extensions []cbc.Key `json:"extensions,omitempty" jsonschema:"title=Extensions"`

	// Map defines a set of regime specific code mappings.
	Map cbc.CodeMap `json:"map,omitempty" jsonschema:"title=Map"`

	// List of sources for the information contained in this category.
	Sources []*Source `json:"sources,omitempty" jsonschema:"title=Sources"`

	// Meta contains additional information about the category that is relevant
	// for local frequently used formats.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Source describes where the information for the taxes comes from.
type Source struct {
	// Title of the linked source to help distinguish between this and other links.
	Title i18n.String `json:"title,omitempty" jsonschema:"title=Title"`
	// URL for the website.
	URL string `json:"url" jsonschema:"title=URL,format=uri"`
}

// Rate defines a single rate inside a category
type Rate struct {
	// Key identifies this rate within the system
	Key cbc.Key `json:"key" jsonschema:"title=Key"`

	// Human name of the rate
	Name i18n.String `json:"name" jsonschema:"title=Name"`
	// Useful description of the rate.
	Description i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`

	// Exempt when true implies that the rate when used in a tax Combo should
	// not define a percent value.
	Exempt bool `json:"exempt,omitempty" jsonschema:"title=Exempt"`

	// Values contains a list of Value objects that contain the
	// current and historical percentage values for the rate and
	// additional filters.
	// Order is important, newer values should come before
	// older values.
	Values []*RateValue `json:"values,omitempty" jsonschema:"title=Values"`

	// Extensions defines a list of keys for codes that can or must be associated with the
	// tax rate for it to be validated. Every key must be defined in the Regime's
	// extensions.
	Extensions []cbc.Key `json:"extensions,omitempty" jsonschema:"title=Extensions"`

	// Map is used to associate specific codes with the chosen rate.
	Map cbc.CodeMap `json:"map,omitempty" jsonschema:"title=Map"`

	// Meta contains additional information about the rate that is relevant
	// for local frequently used implementations.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// RateValue contains a percentage rate or fixed amount for a given date range.
// Fiscal policy changes mean that rates are not static so we need to
// be able to apply the correct rate for a given period.
type RateValue struct {
	// Only apply this rate if one of the tags is present in the invoice.
	Tags []cbc.Key `json:"tags,omitempty" jsonschema:"title=Tags"`
	// Ext map of keys that can be used to filter to determine if the rate applies.
	Ext Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`
	// Date from which this value should be applied.
	Since *cal.Date `json:"since,omitempty" jsonschema:"title=Since"`
	// Percent rate that should be applied
	Percent num.Percentage `json:"percent" jsonschema:"title=Percent"`
	// An additional surcharge to apply.
	Surcharge *num.Percentage `json:"surcharge,omitempty" jsonschema:"title=Surcharge"`
	// When true, this value should no longer be used.
	Disabled bool `json:"disabled,omitempty" jsonschema:"title=Disabled"`
}

// CorrectionDefinition contains details about what can be defined in .
type CorrectionDefinition struct {
	// Partial or complete schema URL for the document type supported by correction.
	Schema string `json:"schema" jsonschema:"title=Schema"`
	// The types of sub-documents supported by the regime
	Types []cbc.Key `json:"types,omitempty" jsonschema:"title=Types"`
	// Extension keys that can be included
	Extensions []cbc.Key `json:"extensions,omitempty" jsonschema:"title=Extensions"`
	// ReasonRequired when true implies that a reason must be provided
	ReasonRequired bool `json:"reason_required,omitempty" jsonschema:"title=Reason Required"`
	// Stamps that must be copied from the preceding document.
	Stamps []cbc.Key `json:"stamps,omitempty" jsonschema:"title=Stamps"`
}

// Code provides a unique code for this tax regime based on the country.
func (r *Regime) Code() cbc.Code {
	return cbc.Code(r.Country)
}

// ValidateObject performs validation on the provided object in the context
// of the regime.
func (r *Regime) ValidateObject(value interface{}) error {
	if r.Validator != nil {
		return r.Validator(value)
	}
	return nil
}

// CalculateObject performs any regime specific calculations on the provided
// object.
func (r *Regime) CalculateObject(obj interface{}) error {
	if r.Calculator != nil {
		return r.Calculator(obj)
	}
	return nil
}

// CurrencyDef provides the currency definition object for the region.
func (r *Regime) CurrencyDef() *currency.Def {
	return currency.Get(r.Currency)
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

// CorrectionDefinitionFor provides the correction definition for the matching schema.
func (r *Regime) CorrectionDefinitionFor(schema string) *CorrectionDefinition {
	for _, c := range r.Corrections {
		if strings.HasSuffix(schema, c.Schema) {
			return c
		}
	}
	return nil
}

// Validate enures the region definition is valid, including all
// subsequent categories.
func (r *Regime) Validate() error {
	return r.ValidateWithContext(context.Background())
}

// ValidateWithContext enures the region definition is valid, including all
// subsequent categories, and passes through the context.
func (r *Regime) ValidateWithContext(ctx context.Context) error {
	ctx = context.WithValue(ctx, KeyRegime, r)
	err := validation.ValidateStructWithContext(ctx, r,
		validation.Field(&r.Country, validation.Required),
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Description),
		validation.Field(&r.TimeZone, validation.Required, validation.By(validateTimeZone)),
		validation.Field(&r.Country),
		validation.Field(&r.Zone),
		validation.Field(&r.Currency),
		validation.Field(&r.Tags),
		validation.Field(&r.IdentityTypeKeys),
		validation.Field(&r.Extensions),
		validation.Field(&r.ChargeKeys),
		validation.Field(&r.PaymentMeansKeys),
		validation.Field(&r.InboxKeys),
		validation.Field(&r.Scenarios),
		validation.Field(&r.Corrections),
		validation.Field(&r.Categories, validation.Required),
	)
	return err
}

func validateTimeZone(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return errors.New("invalid time zone")
	}
	_, err := time.LoadLocation(s)
	return err
}

// TimeLocation returns the time.Location for the regime.
func (r *Regime) TimeLocation() *time.Location {
	if r == nil {
		return nil
	}
	loc, err := time.LoadLocation(r.TimeZone)
	if err != nil {
		return time.UTC
	}
	return loc
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

// ValidateWithContext ensures the Category's contents are correct.
func (c *Category) ValidateWithContext(ctx context.Context) error {
	reg := ctx.Value(KeyRegime).(*Regime)
	err := validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.Code, validation.Required),
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Title, validation.Required),
		validation.Field(&c.Description),
		validation.Field(&c.Sources),
		validation.Field(&c.Rates),
		validation.Field(&c.Extensions,
			validation.Each(cbc.InKeyDefs(reg.Extensions)),
		),
		validation.Field(&c.Map),
	)
	return err
}

// Validate ensures the Source's contents are correct.
func (s *Source) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Title),
		validation.Field(&s.URL, validation.Required, is.URL),
	)
}

// ValidateWithContext checks that our tax definition is valid. This is only really
// meant to be used when testing new regional tax definitions.
func (r *Rate) ValidateWithContext(ctx context.Context) error {
	reg := ctx.Value(KeyRegime).(*Regime)
	err := validation.ValidateStructWithContext(ctx, r,
		validation.Field(&r.Key, validation.Required),
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Values,
			validation.When(r.Exempt, validation.Nil),
			validation.By(checkRateValuesOrder),
		),
		validation.Field(&r.Extensions,
			validation.Each(cbc.InKeyDefs(reg.Extensions)),
		),
		validation.Field(&r.Map),
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
		if len(v.Tags) > 0 || len(v.Ext) > 0 {
			// TODO: check tags and extensions order also
			// Not too important at the moment.
			continue
		}
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
func (r *Regime) Tag(key cbc.Key) *cbc.KeyDefinition {
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

// ExtensionDef provides the extension definition with a matching key.
func (r *Regime) ExtensionDef(key cbc.Key) *cbc.KeyDefinition {
	for _, e := range r.Extensions {
		if e.Key == key {
			return e
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
func (r *Rate) Value(date cal.Date, tags []cbc.Key, ext Extensions) *RateValue {
	for _, rv := range r.Values {
		if len(rv.Tags) > 0 {
			if !rv.HasATag(tags) {
				continue
			}
		}
		if len(rv.Ext) > 0 {
			if !ext.Contains(rv.Ext) {
				continue
			}
		}
		if rv.Since == nil || !rv.Since.IsValid() || rv.Since.Before(date.Date) {
			return rv
		}
	}
	return nil
}

// HasATag returns true if the rate value has a tag that matches
// one of those provided.
func (rv *RateValue) HasATag(tags []cbc.Key) bool {
	for _, t := range rv.Tags {
		for _, tag := range tags {
			if t == tag {
				return true
			}
		}
	}
	return false
}

// HasType returns true if the correction definition has a type that matches the one provided.
func (cd *CorrectionDefinition) HasType(t cbc.Key) bool {
	if cd == nil {
		return false // no preceding definitions
	}
	return t.In(cd.Types...)
}

// HasExtension returns true if the correction definition has the change key provided.
func (cd *CorrectionDefinition) HasExtension(key cbc.Key) bool {
	if cd == nil {
		return false // no correction definitions
	}
	return key.In(cd.Extensions...)
}

// Validate ensures the key definition looks correct in the context of the regime.
func (cd *CorrectionDefinition) Validate() error {
	err := validation.ValidateStruct(cd,
		validation.Field(&cd.Schema, validation.Required),
		validation.Field(&cd.Types),
		validation.Field(&cd.Stamps),
		validation.Field(&cd.Extensions),
	)
	return err
}
