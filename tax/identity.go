package tax

import (
	"context"
	"errors"
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/uuid"

	"github.com/invopop/validation"
)

// Identity stores the details required to identify an entity for tax
// purposes. There are two levels of accuracy that may be used to
// describe where an entity is located: Country and Locality.
// Country is a required field, but locality is optional according to
// rules of a given tax jurisdiction.
type Identity struct {
	// Unique universal identity code for this tax identity.
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`

	// ISO country code for Where the tax identity was issued.
	Country l10n.CountryCode `json:"country" jsonschema:"title=Country Code"`

	// Where inside the country the tax identity holder is based for tax purposes
	// like a village, town, district, city, county, state or province. For some
	// areas, this could be a regular post or zip code. See the regime packages
	// for specific validation rules.
	Zone l10n.Code `json:"zone,omitempty" jsonschema:"title=Zone Code"`

	// Type is set according the requirements of each regime, some have a single
	// tax document type code, others require a choice to be made.
	Type cbc.Key `json:"type,omitempty" jsonschema:"title=Type"`

	// Normalized code shown on the original identity document.
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`

	// Additional details that may be required.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Standard error responses to be used by regimes.
var (
	ErrIdentityCodeInvalid = errors.New("invalid tax identity code")
)

// RequireIdentityCode is an additional check to use alongside
// regular validation that will ensure the tax ID has a code
// value set.
var RequireIdentityCode = validateTaxID{requireCode: true}

// RequireIdentityType ensures that the identity type is set.
var RequireIdentityType = validateTaxID{requireType: true}

// IdentityTypeIn checks that the identity code is within one of the
// acceptable keys.
var IdentityTypeIn = func(keys ...cbc.Key) validation.Rule {
	out := make([]interface{}, len(keys))
	for i, l := range keys {
		out[i] = l
	}
	return validateTaxID{typeIn: out}
}

type validateTaxID struct {
	requireCode bool
	requireType bool
	typeIn      []interface{}
}

// String provides a string representation of the tax identity.
func (id *Identity) String() string {
	return fmt.Sprintf("%s%s", id.Country, id.Code)
}

// Regime provides the regime object for this tax identity.
func (id *Identity) Regime() *Regime {
	return regimes.For(id.Country, id.Zone)
}

// Calculate will attempt to perform a regional tax normalization
// on the tax identity.
func (id *Identity) Calculate(ctx context.Context) error {
	r := id.Regime()
	if r != nil {
		return r.CalculateObject(ctx, id)
	}
	return nil
}

// Validate checks to ensure the tax ID contains all the required
// fields and performs any regime specific validation based on the ID's
// country and zone properties.
func (id *Identity) Validate() error {
	err := validation.ValidateStruct(id,
		validation.Field(&id.UUID),
		validation.Field(&id.Country, validation.Required),
		validation.Field(&id.Zone),
		validation.Field(&id.Type),
		validation.Field(&id.Code),
		validation.Field(&id.Meta),
	)
	if err != nil {
		return err
	}
	r := regimes.For(id.Country, id.Zone)
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
	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.When(v.requireCode, validation.Required),
		),
		validation.Field(&id.Type,
			validation.When(v.requireType, validation.Required),
			validation.When(len(v.typeIn) > 0, validation.In(v.typeIn...)),
		),
	)
}
