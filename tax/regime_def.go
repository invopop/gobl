package tax

import (
	"context"
	"errors"
	"fmt"
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
	keyRegime contextKey = "regime"
)

// RegimeDef defines the holding structure for the definitions of taxes inside a country
// or territory.
type RegimeDef struct {
	// Name of the tax regime.
	Name i18n.String `json:"name" jsonschema:"title=Name"`

	// Introductory details about the regime.
	Description i18n.String `json:"description,omitempty" jsonschema:"title=Description"`

	// Location name for the country's central time zone. Accepted
	// values from IANA Time Zone Database (https://iana.org/time-zones).
	TimeZone string `json:"time_zone" jsonschema:"title=Time Zone"`

	// Country code for the region
	Country l10n.TaxCountryCode `json:"country" jsonschema:"title=Code"`

	// Alternative localization codes that may be used to identify the tax regime
	// in specific circumstances.
	AltCountryCodes []l10n.Code `json:"alt_country_codes,omitempty" jsonschema:"title=Alternative Country Codes"`

	// Specific Locality, region, city, province, county, or similar code inside
	// the country, if needed.
	Zone l10n.Code `json:"zone,omitempty" jsonschema:"title=Zone"`

	// Currency used by the country.
	Currency currency.Code `json:"currency" jsonschema:"title=Currency"`

	// Rounding rule to use when calculating the tax totals, default is always
	// `sum-then-round`.
	CalculatorRoundingRule CalculatorRoundingRule `json:"calculator_rounding_rule,omitempty" jsonschema:"title=Calculator Rounding Rule"`

	// Tags that can be applied at the document level to identify additional
	// considerations.
	Tags []*TagSet `json:"tags,omitempty" jsonschema:"title=Tags"`

	// Extensions defines the keys that can be used for extended or extra data inside the regime that
	// is specific to the regime and cannot be easily determined from other GOBL structures.
	// Typically these are used to define local codes for suppliers, customers, products, or tax rates.
	Extensions []*cbc.Definition `json:"extensions,omitempty" jsonschema:"title=Extensions"`

	// Identities used in addition to regular tax identities and specific for the
	// regime that may be validated against.
	Identities []*cbc.Definition `json:"identities,omitempty" jsonschema:"title=Identities"`

	// PaymentMeansKeys specific for the regime that extend the original
	// base payment means keys.
	PaymentMeansKeys []*cbc.Definition `json:"payment_means_keys,omitempty" jsonschema:"title=Payment Means Keys"`

	// InboxKeys specific to the regime that can be used to identify where a document
	// should be forwarded to.
	InboxKeys []*cbc.Definition `json:"inbox_keys,omitempty" jsonschema:"title=Inbox Keys"`

	Scenarios []*ScenarioSet `json:"scenarios,omitempty" jsonschema:"title=Scenarios"`

	// Configuration details for corrections to be used with correction options.
	Corrections CorrectionSet `json:"corrections,omitempty" jsonschema:"title=Corrections"`

	// List of tax categories.
	Categories []*CategoryDef `json:"categories" jsonschema:"title=Categories"`

	// Validator is a method to use to validate a document in a given region.
	Validator Validator `json:"-"`

	// Normalizer is used to perform regime specific normalizations on data,
	// that might need to take place such as with tax codes and removing white-space.
	Normalizer Normalizer `json:"-"`
}

// CategoryDef contains the definition of a general type of tax inside a region.
type CategoryDef struct {
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

	// Specific tax definitions inside this category. Order is important.
	Rates []*RateDef `json:"rates,omitempty" jsonschema:"title=Rates"`

	// Extensions defines a list of extension keys that may be used or required
	// as an alternative or alongside choosing a rate for the tax category.
	// Every key must be defined in the Regime's extensions table.
	Extensions []cbc.Key `json:"extensions,omitempty" jsonschema:"title=Extensions"`

	// Map defines a set of regime specific code mappings.
	Map cbc.CodeMap `json:"map,omitempty" jsonschema:"title=Map"`

	// List of sources for the information contained in this category.
	Sources []*Source `json:"sources,omitempty" jsonschema:"title=Sources"`

	// Extensions key-value pairs that will be copied to the tax combo if this
	// category is used.
	Ext Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`

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

// RateDef defines a single rate inside a category
type RateDef struct {
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
	Values []*RateValueDef `json:"values,omitempty" jsonschema:"title=Values"`

	// Extensions key-value pair that will be copied to the tax combo if this
	// rate is used.
	Ext Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`

	// Map is used to associate specific codes with the chosen rate.
	// Map cbc.CodeMap `json:"map,omitempty" jsonschema:"title=Map"`

	// Meta contains additional information about the rate that is relevant
	// for local frequently used implementations.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// RateValueDef contains a percentage rate or fixed amount for a given date range.
// Fiscal policy changes mean that rates are not static so we need to
// be able to apply the correct rate for a given period.
type RateValueDef struct {
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

// WithContext adds this regime to the given context, alongside
// its validator and tags in the contexts.
func (r *RegimeDef) WithContext(ctx context.Context) context.Context {
	if r == nil {
		return ctx
	}
	ctx = context.WithValue(ctx, keyRegime, r)
	ctx = contextWithValidator(ctx, r.Validator)
	return ctx
}

// RegimeDefFromContext returns the regime from the given context, or nil.
func RegimeDefFromContext(ctx context.Context) *RegimeDef {
	r, ok := ctx.Value(keyRegime).(*RegimeDef)
	if !ok {
		return nil
	}
	return r
}

// Code provides a unique code for this tax regime based on the country.
func (r *RegimeDef) Code() cbc.Code {
	return cbc.Code(r.Country)
}

// ValidateObject performs validation on the provided object in the context
// of the regime.
func (r *RegimeDef) ValidateObject(value interface{}) error {
	if r == nil {
		return nil
	}
	if r.Validator != nil {
		return r.Validator(value)
	}
	return nil
}

// NormalizeObject performs any regime specific normalizations on the provided
// object.
func (r *RegimeDef) NormalizeObject(obj interface{}) {
	if r == nil {
		return
	}
	if r.Normalizer != nil {
		r.Normalizer(obj)
	}
}

// CurrencyDef provides the currency definition object for the region.
func (r *RegimeDef) CurrencyDef() *currency.Def {
	return currency.Get(r.Currency)
}

// ScenarioSet returns a single scenario set instance for the provided
// document schema.
func (r *RegimeDef) ScenarioSet(schema string) *ScenarioSet {
	for _, s := range r.Scenarios {
		if strings.HasSuffix(schema, s.Schema) {
			return s
		}
	}
	return nil
}

// Validate enures the region definition is valid, including all
// subsequent categories.
func (r *RegimeDef) Validate() error {
	return r.ValidateWithContext(context.Background())
}

// ValidateWithContext enures the region definition is valid, including all
// subsequent categories, and passes through the context.
func (r *RegimeDef) ValidateWithContext(ctx context.Context) error {
	ctx = r.WithContext(ctx)
	err := validation.ValidateStructWithContext(ctx, r,
		validation.Field(&r.Country, validation.Required),
		validation.Field(&r.AltCountryCodes),
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Description),
		validation.Field(&r.TimeZone, validation.Required, validation.By(validateTimeZone)),
		validation.Field(&r.Country),
		validation.Field(&r.Zone),
		validation.Field(&r.Currency),
		validation.Field(&r.Tags),
		validation.Field(&r.Identities),
		validation.Field(&r.Extensions),
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
func (r *RegimeDef) TimeLocation() *time.Location {
	if r == nil {
		return time.UTC
	}
	loc, err := time.LoadLocation(r.TimeZone)
	if err != nil {
		return time.UTC
	}
	return loc
}

type inCategoryRatesRule struct {
	cat  cbc.Code
	keys []cbc.Key
}

func (r *inCategoryRatesRule) Validate(value any) error {
	key, ok := value.(cbc.Key)
	if !ok || key == cbc.KeyEmpty {
		return nil
	}
	for _, k := range r.keys {
		if key.Has(k) {
			return nil
		}
	}
	return fmt.Errorf("'%v' not defined in '%v' category", key, r.cat)
}

// InCategoryRates is used to provide a validation rule that will
// ensure a rate key is represented inside a category.
func (r *RegimeDef) InCategoryRates(cat cbc.Code) validation.Rule {
	if r == nil {
		return validation.Empty.Error("must be blank when regime is undefined")
	}
	c := r.CategoryDef(cat)
	if c == nil {
		return validation.Empty.Error("must be blank when category is undefined")
	}
	keys := make([]cbc.Key, len(c.Rates))
	for i, x := range c.Rates {
		keys[i] = x.Key
	}
	return &inCategoryRatesRule{cat: cat, keys: keys}
}

// InCategories returns a validation rule to ensure the category code
// is inside the list of known codes.
func (r *RegimeDef) InCategories() validation.Rule {
	if r == nil {
		return validation.Skip
	}
	cats := make([]cbc.Code, len(r.Categories))
	for i, c := range r.Categories {
		cats[i] = c.Code
	}
	return validation.In(cats...)
}

// ValidateWithContext ensures the Category's contents are correct.
func (c *CategoryDef) ValidateWithContext(ctx context.Context) error {
	r := RegimeDefFromContext(ctx)
	err := validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.Code, validation.Required),
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Title, validation.Required),
		validation.Field(&c.Description),
		validation.Field(&c.Sources),
		validation.Field(&c.Rates),
		validation.Field(&c.Extensions,
			validation.Each(cbc.InKeyDefs(r.Extensions)),
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
func (r *RateDef) ValidateWithContext(ctx context.Context) error {
	err := validation.ValidateStructWithContext(ctx, r,
		validation.Field(&r.Key, validation.Required),
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Values,
			validation.When(r.Exempt, validation.Nil),
			validation.By(checkRateValuesOrder),
		),
		validation.Field(&r.Ext),
		validation.Field(&r.Meta),
	)
	return err
}

// Validate ensures the tax rate contains all the required fields.
func (rv *RateValueDef) Validate() error {
	return validation.ValidateStruct(rv,
		validation.Field(&rv.Percent, validation.Required),
	)
}

func checkRateValuesOrder(list interface{}) error {
	values, ok := list.([]*RateValueDef)
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

// CategoryDef provides the requested category definition by its code.
func (r *RegimeDef) CategoryDef(code cbc.Code) *CategoryDef {
	if r == nil {
		return nil
	}
	for _, c := range r.Categories {
		if c.Code == code {
			return c
		}
	}
	return nil
}

// RateDef provides the rate definition for the provided category code
// and rate key.
func (r *RegimeDef) RateDef(cat cbc.Code, key cbc.Key) *RateDef {
	c := r.CategoryDef(cat)
	if c == nil {
		return nil
	}
	return c.RateDef(key)
}

// ExtensionDef provides the extension definition with a matching key.
func (r *RegimeDef) ExtensionDef(key cbc.Key) *cbc.Definition {
	for _, e := range r.Extensions {
		if e.Key == key {
			return e
		}
	}
	return nil
}

// RateDef provides the rate definition with a matching key for
// the category. Key comparison is made using two loops. The first
// will find an exact match, while the second will see if the provided
// key has the rate key as a prefix.
func (c *CategoryDef) RateDef(key cbc.Key) *RateDef {
	for _, r := range c.Rates {
		if r.Key == key {
			return r
		}
	}
	for _, r := range c.Rates {
		if key.Has(r.Key) {
			return r
		}
	}
	return nil
}

// Value determines the tax rate value for the provided date and zone, if applicable.
func (r *RateDef) Value(date cal.Date, tags []cbc.Key, ext Extensions) *RateValueDef {
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
func (rv *RateValueDef) HasATag(tags []cbc.Key) bool {
	for _, t := range rv.Tags {
		for _, tag := range tags {
			if t == tag {
				return true
			}
		}
	}
	return false
}
