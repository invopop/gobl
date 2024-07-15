package tax

import (
	"errors"
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"

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
	// ISO country code for Where the tax identity was issued.
	Country l10n.CountryCode `json:"country" jsonschema:"title=Country Code"`

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

// String provides a string representation of the tax identity.
func (id *Identity) String() string {
	return fmt.Sprintf("%s%s", id.Country, id.Code)
}

// Regime provides the regime object for this tax identity.
func (id *Identity) Regime() *Regime {
	return regimes.For(id.Country)
}

// Calculate will attempt to perform a regional tax normalization
// on the tax identity.
func (id *Identity) Calculate() error {
	r := id.Regime()
	if r != nil {
		return r.CalculateObject(id)
	}
	return nil
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
	r := regimes.For(id.Country)
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
