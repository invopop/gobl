package tax

import (
	"errors"
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/jsonschema"

	"github.com/invopop/validation"
)

// Identity stores the details required to identify an entity for tax
// purposes in a specific country. Typically this would be a code related
// to a specific indirect tax like VAT or GST. Some countries, such as the
// US, do not have a VAT system so will not have a code here.
//
// Other fiscal identities should be defined in a parties identities array
// with their own validation rules and country specific handling.
type Identity struct {
	// Tax country code for Where the tax identity was issued.
	Country l10n.TaxCountryCode `json:"country" jsonschema:"title=Country Code"`

	// Normalized code shown on the original identity document.
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`

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

// Standard error responses to be used by regimes.
var (
	ErrIdentityCodeInvalid = errors.New("invalid tax identity code")
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
	if err := id.Normalize(); err != nil {
		return nil, err
	}
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
func (id *Identity) Regime() *Regime {
	if id == nil {
		return nil
	}
	return regimes.For(id.Country.Code())
}

// Normalize will attempt to perform a regional tax normalization
// on the tax identity.
func (id *Identity) Normalize() error {
	r := id.Regime()
	if r != nil {
		return r.CalculateObject(id)
	}
	return nil
}

// Calculate is an alias for Normalize and will perform normalization
// on the tax identity code.
func (id *Identity) Calculate() error {
	return id.Normalize()
}

// Validate checks to ensure the tax ID contains all the required
// fields and performs any regime specific validation based on the ID's
// country and zone properties.
func (id *Identity) Validate() error {
	err := validation.ValidateStruct(id,
		validation.Field(&id.Country, validation.Required),
		validation.Field(&id.Code),
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

func (v validateTaxID) Validate(value interface{}) error {
	id, ok := value.(*Identity)
	if !ok {
		return nil
	}
	rules := []*validation.FieldRules{
		validation.Field(&id.Code,
			validation.When(v.requireCode, validation.Required),
		),
	}
	return validation.ValidateStruct(id, rules...)
}

// JSONSchemaExtend adds extra details to the schema.
func (Identity) JSONSchemaExtend(js *jsonschema.Schema) {
	js.Extras = map[string]any{
		schema.Recommended: []string{
			"code",
		},
	}
}
