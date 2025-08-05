package tax

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/validation"
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

	// TaxScheme defines the principal scheme of consumption tax that should be
	// applied to the regime and associated with Tax IDs in some export formats
	// such as UBL or CII. Some regimes may not have a Tax Scheme and as a
	// consequence will not use tax identities, like the US.
	TaxScheme cbc.Code `json:"tax_scheme,omitempty" jsonschema:"title=Tax Scheme"`

	// Rounding rule to use when calculating the tax totals, default is always
	// `sum-then-round`.
	CalculatorRoundingRule cbc.Key `json:"calculator_rounding_rule,omitempty" jsonschema:"title=Calculator Rounding Rule"`

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

// WithContext adds this regime to the given context, alongside
// its validator and tags in the contexts.
func (r *RegimeDef) WithContext(ctx context.Context) context.Context {
	if r == nil {
		return ctx
	}
	ctx = context.WithValue(ctx, keyRegime, r)
	ctx = ContextWithValidator(ctx, r.Validator)
	return ctx
}

// Normalizers returns the normalizers for this regime, if any,
// handling any potential for nil pointers.
func (r *RegimeDef) Normalizers() Normalizers {
	if r == nil || r.Normalizer == nil {
		return nil
	}
	return Normalizers{r.Normalizer}
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

// GetCurrency is a convenience method that will always return
// a code, even if the RegimeDef is nil.
func (r *RegimeDef) GetCurrency() currency.Code {
	if r != nil {
		return r.Currency
	}
	return currency.CodeEmpty
}

// GetCountry provides the country code for the regime, or an empty string.
func (r *RegimeDef) GetCountry() l10n.TaxCountryCode {
	if r != nil {
		return r.Country
	}
	return ""
}

// GetRoundingRule provides the regime's rounding rule, or the default.
func (r *RegimeDef) GetRoundingRule() cbc.Key {
	if r != nil && r.CalculatorRoundingRule != "" {
		return r.CalculatorRoundingRule
	}
	return RoundingRulePrecise
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
		validation.Field(&r.TaxScheme),
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

type inCategoryRule struct {
	cat   cbc.Code
	key   cbc.Key
	rates []cbc.Key
}

func (r *inCategoryRule) Validate(value any) error {
	rate, ok := value.(cbc.Key)
	if !ok || rate == cbc.KeyEmpty {
		return nil
	}
	if len(r.rates) == 0 {
		return nil
	}
	for _, k := range r.rates {
		if rate.Has(k) {
			return nil
		}
	}
	return fmt.Errorf("'%v' not defined in '%v' category for key '%s'", rate, r.cat, r.key)
}

// InCategoryRates is used to provide a validation rule that will
// ensure a rate is represented inside a category and has a key.
func (r *RegimeDef) InCategoryRates(cat cbc.Code, key cbc.Key) validation.Rule {
	if r == nil {
		return validation.Empty.Error("must be blank when regime is undefined")
	}
	cd := r.CategoryDef(cat)
	if cd == nil {
		return validation.Empty.Error(fmt.Sprintf("must be blank for undefined category '%s'", cat.String()))
	}

	rates := make([]cbc.Key, 0)
	for _, rd := range cd.Rates {
		if rd.HasKey(key) {
			rates = append(rates, rd.Rate)
		}
	}
	if len(rates) == 0 {
		if key == cbc.KeyEmpty {
			return validation.Empty.Error(fmt.Sprintf("must be blank for category '%s' with no key", cat))
		}
		return validation.Empty.Error(fmt.Sprintf("must be blank for category '%s' and key '%s'", cat, key))
	}
	return &inCategoryRule{cat: cat, key: key, rates: rates}
}

// InCategoryKeys returns a validation rule to ensure the key is inside the
// list of known keys for the category.
func (r *RegimeDef) InCategoryKeys(cat cbc.Code) validation.Rule {
	cd := r.CategoryDef(cat)
	if cd == nil {
		return validation.Empty.Error("must be blank when category is undefined")
	}

	keys := make([]cbc.Key, len(cd.Keys))
	for i, k := range cd.Keys {
		keys[i] = k.Key
	}

	return validation.In(keys...)
}

// RequiresPercent returns a validation rule that will ensure the percent
// is either nil or a valid percent for the given category and key.
func (r *RegimeDef) RequiresPercent(cat cbc.Code, key cbc.Key) validation.Rule {
	return validation.By(func(value any) error {
		p, ok := value.(*num.Percentage)
		if !ok {
			return nil
		}

		cd := r.CategoryDef(cat)
		if cd == nil {
			return nil // nothing to lookup
		}
		kd := cd.KeyDef(key)
		if kd == nil {
			return nil // no key definition, no percent required
		}
		if kd.NoPercent {
			if p != nil {
				return fmt.Errorf("must be nil for '%s' in '%s'", key.String(), cat.String())
			}
		} else {
			if p == nil {
				return fmt.Errorf("required for '%s' in '%s'", key.String(), cat.String())
			}
		}
		return nil
	})
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

// CategoryDef provides the requested category definition by its code.
func (r *RegimeDef) CategoryDef(code cbc.Code) *CategoryDef {
	if r == nil {
		// Lookup a global category definition
		return Category(code)
	}
	for _, c := range r.Categories {
		if c.Code == code {
			return c
		}
	}
	return nil
}
