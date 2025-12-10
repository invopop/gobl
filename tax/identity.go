package tax

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/jsonschema"

	"github.com/invopop/validation"
)

// Identity stores the details required to identify an entity for tax
// purposes in a specific country. Typically this would be a code related
// to a specific indirect tax scheme like VAT or GST. Some countries, such as the
// US, do not have an official tax scheme and should omit the `code` field.
//
// Other fiscal identities should be defined in a party identities array
// with their own validation rules and country specific handling.
type Identity struct {
	// Tax country code for Where the tax identity was issued.
	Country l10n.TaxCountryCode `json:"country" jsonschema:"title=Country Code"`

	// Normalized code shown on the original identity document.
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`

	// Scheme is an optional field that may be used to override the tax regime's
	// default tax scheme. Many electronic formats such as UBL or CII define an
	// equivalent field. Examples: `VAT`, `GST`, `ST`, etc.
	Scheme cbc.Code `json:"scheme,omitempty" jsonschema:"title=Scheme"`

	// Type is set according to the requirements of each regime, some have a single
	// tax document type code, others require a choice to be made.
	//
	// Deprecated: Tax Identities should only be used for VAT or similar codes
	// for companies. Use the identities array for other types of identification.
	Type cbc.Key `json:"type,omitempty" jsonschema:"title=Type"`

	// Zone identifies a sub-locality within a country.
	//
	// Deprecated: Removed 2024-03-14 in favour of using tax tags
	// and extensions with local data when required. Maintained here to support
	// data migration.
	Zone l10n.Code `json:"zone,omitempty" jsonschema:"title=Zone"`
}

var (
	// IdentityCodePattern is the regular expression pattern used to validate tax identity codes.
	IdentityCodePattern = `^[A-Z0-9]+$`

	// IdentityCodePatternRegexp is the regular expression used to validate tax identity codes.
	IdentityCodePatternRegexp = regexp.MustCompile(IdentityCodePattern)

	// ErrIdentityCodeInvalid is returned when the tax identity code is not valid.
	ErrIdentityCodeInvalid = errors.New("invalid tax identity code")

	// IdentityCodeBadCharsRegexp is used to remove any characters that are not valid in a tax code.
	IdentityCodeBadCharsRegexp = regexp.MustCompile(`[^A-Z0-9]+`)

	// IdentityCodeValidationIgnore is a list of countries that should not have their tax identity
	// codes validated due to local rules.
	IdentityCodeValidationIgnore = []l10n.TaxCountryCode{"MX"}
)

// RequireIdentityCode is an additional check to use alongside
// regular validation that will ensure the tax ID has a code
// value set.
var RequireIdentityCode = validateTaxID{requireCode: true}

type validateTaxID struct {
	requireCode bool
}

// ParseIdentity will attempt to parse a tax identity from a string making
// the assumption that the first two characters are the country code and
// the rest is the tax code. If the country code is identified by a
// tax regime, the code will be normalized and validated.
func ParseIdentity(tin string) (*Identity, error) {
	if len(tin) < 2 {
		return nil, ErrIdentityCodeInvalid
	}
	id := &Identity{
		Country: l10n.TaxCountryCode(tin[:2]),
		Code:    cbc.Code(tin[2:]),
	}
	id.Normalize()
	if err := id.Validate(); err != nil {
		return nil, err
	}
	return id, nil
}

// String provides a string representation of the tax identity.
func (id *Identity) String() string {
	return fmt.Sprintf("%s%s", id.Country, id.Code)
}

// Regime provides the regime object for this tax identity.
func (id *Identity) Regime() *RegimeDef {
	if id == nil {
		return nil
	}
	return regimes.For(id.Country.Code())
}

// Calculate will simply perform normalization.
func (id *Identity) Calculate() error {
	id.Normalize()
	return nil
}

// GetScheme can be used to determine the tax identities Scheme
// either from the value defined directly in the identity, or
// from the tax regime.
func (id *Identity) GetScheme() cbc.Code {
	if id.Scheme != cbc.CodeEmpty {
		return id.Scheme
	}
	if r := regimes.For(id.Country.Code()); r != nil {
		return r.TaxScheme
	}
	return cbc.CodeEmpty
}

// Normalize will attempt to perform a regional tax normalization
// on the tax identity. Identities are an exception to the normal
// normalization rules as they cannot be normalized using addons.
func (id *Identity) Normalize() {
	if r := id.Regime(); r != nil {
		r.NormalizeObject(id)
	} else {
		// Fallback to common normalization
		NormalizeIdentity(id)
	}
}

// Validate checks to ensure the tax ID contains all the required
// fields and performs any regime specific validation based on the ID's
// country and zone properties.
func (id *Identity) Validate() error {
	err := validation.ValidateStruct(id,
		validation.Field(&id.Country, validation.Required),
		validation.Field(&id.Code,
			validation.Skip.When(
				id.Country.In(IdentityCodeValidationIgnore...),
			),
			validation.Match(IdentityCodePatternRegexp),
		),
		validation.Field(&id.Scheme),
		validation.Field(&id.Zone, validation.Empty),
		validation.Field(&id.Type),
	)
	if err != nil {
		return err
	}
	r := regimes.For(id.Country.Code())
	if r != nil {
		return r.ValidateObject(id)
	}
	return nil
}

// InEU checks if the tax identity is from a country that is part of the EU on
// the given date.
func (id *Identity) InEU(date cal.Date) bool {
	return l10n.Union(l10n.EU).HasMemberOn(date, id.Country.Code())
}

func (v validateTaxID) Validate(value any) error {
	id, ok := value.(*Identity)
	if id == nil || !ok {
		return nil
	}
	rules := []*validation.FieldRules{}
	if v.requireCode {
		rules = append(rules,
			validation.Field(&id.Code,
				validation.Required,
				validation.Skip,
			),
		)
	}
	return validation.ValidateStruct(id, rules...)
}

// JSONSchemaExtend adds extra details to the schema.
func (Identity) JSONSchemaExtend(js *jsonschema.Schema) {
	if cp, ok := js.Properties.Get("code"); ok {
		cp.Pattern = IdentityCodePattern
	}
	js.Extras = map[string]any{
		schema.Recommended: []string{
			"code",
		},
	}
}

// NormalizeIdentity removes any whitespace or separation characters and ensures all letters are
// uppercase.
func NormalizeIdentity(tID *Identity, altCodes ...l10n.Code) {
	if tID == nil {
		return
	}
	code := strings.ToUpper(tID.Code.String())
	code = IdentityCodeBadCharsRegexp.ReplaceAllString(code, "")
	code = strings.TrimPrefix(code, string(tID.Country))
	for _, alt := range altCodes {
		code = strings.TrimPrefix(code, string(alt))
	}
	tID.Code = cbc.Code(code)
}
