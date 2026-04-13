package tax

import (
	"errors"
	"strings"
	"time"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// RegimeDef defines the holding structure for the definitions of taxes inside a country
// or territory.
type RegimeDef struct {
	// Name of the tax regime.
	Name i18n.String `json:"name" jsonschema:"title=Name"`

	// Introductory details about the regime.
	Description i18n.String `json:"description,omitempty" jsonschema:"title=Description"`

	// Sources is a list of references to tax authority or other relevant documentation that
	// can be used to validate the regime's data and be used in the future to check for updates
	// and changes.
	Sources []*cbc.Source `json:"sources,omitempty" jsonschema:"title=Sources"`

	// Location name for the country's central time zone. Accepted
	// values from IANA Time Zone Database (https://iana.org/time-zones). If a country has multiple
	// time zones, the most common or central one should be used.
	TimeZone string `json:"time_zone" jsonschema:"title=Time Zone"`

	// Country code for tax purposes which usually coincides with the ISO 3166-1 alpha-2 code, but not
	// always.
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

	// Rounding rule to use when calculating the tax totals. See the RoundingRule
	// constants for more details. If not provided, the default is RoundingRulePrecise.
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

	// Scenarios are used to describe a specific set of conditions and rules that apply to a specific
	// document schema. These provide a more generic solution for normalization and validation of documents
	// in regimes with specific requirements.
	Scenarios []*ScenarioSet `json:"scenarios,omitempty" jsonschema:"title=Scenarios"`

	// Configuration details for corrections to be used with correction options.
	Corrections CorrectionSet `json:"corrections,omitempty" jsonschema:"title=Corrections"`

	// List of tax categories.
	Categories []*CategoryDef `json:"categories" jsonschema:"title=Categories"`

	// Normalizer is used to perform regime specific normalizations on data,
	// that might need to take place such as with tax codes and removing white-space.
	Normalizer Normalizer `json:"-"`
}

// Normalizers returns the normalizers for this regime, if any,
// handling any potential for nil pointers.
func (r *RegimeDef) Normalizers() Normalizers {
	if r == nil || r.Normalizer == nil {
		return nil
	}
	return Normalizers{r.Normalizer}
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

func regimeDefRules() *rules.Set {
	return rules.For(new(RegimeDef),
		rules.Field("country",
			rules.Assert("01", "country code is required", is.Present),
		),
		rules.Field("name",
			rules.Assert("02", "name is required", is.Present),
		),
		rules.Field("time_zone",
			rules.Assert("03", "time zone is required", is.Present),
			rules.AssertIfPresent("04", "time zone must be valid",
				is.FuncError("valid time zone", validateTimeZone),
			),
		),
		rules.Field("categories",
			rules.Assert("05", "at least one category is required",
				is.Length(1, 0),
			),
		),
	)
}

func validateTimeZone(value any) error {
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
