package tax

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"

	validation "github.com/go-ozzo/ozzo-validation/v4"
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

	// What is the source document of the tax identity.
	Source SourceKey `json:"source,omitempty" jsonschema:"title=Source Key"`

	// Normalized code shown on the original identity document.
	Code string `json:"code,omitempty" jsonschema:"title=Code"`

	// Additional details that may be required.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// SourceKey is used to identify different sources of tax
// identities that may be required by some regions.
type SourceKey cbc.Key

// DefSourceKey defines the details we have regarding a document
// source key.
type DefSourceKey struct {
	Key         SourceKey `json:"key" jsonschema:"title=Key"`
	Description string    `json:"description" jsonschema:"title=Description"`
}

// RequireIdentityCode is an additional check to use alongside
// regular validation that will ensure the tax ID has a code
// value set.
var RequireIdentityCode = validateTaxID{requireCode: true}

type validateTaxID struct {
	requireCode bool
}

// Main Source Key definitions.
const (
	// Directly from tax Agency
	SourceKeyTaxAgency SourceKey = "tax-agency"
	// A passport document
	SourceKeyPassport SourceKey = "passport"
	// National ID Card or similar
	SourceKeyNational SourceKey = "national"
	// Residential permit
	SourceKeyPermit SourceKey = "permit"
	// Something else
	SourceKeyOther SourceKey = "other"
)

// SourceKeyDefinitions lists all the keys with their descriptions
var SourceKeyDefinitions = []DefSourceKey{
	{
		Key:         SourceKeyTaxAgency,
		Description: "Sourced directly from a tax agency",
	},
	{
		Key:         SourceKeyPassport,
		Description: "A passport document",
	},
	{
		Key:         SourceKeyNational,
		Description: "National ID Card or similar",
	},
	{
		Key:         SourceKeyPermit,
		Description: "Residential or similar permit",
	},
	{
		Key:         SourceKeyOther,
		Description: "An other type of source not listed",
	},
}

// Regime provides the regime object for this tax identity.
func (id *Identity) Regime() *Regime {
	return regimes.For(id.Country, id.Zone)
}

// Calculate will attempt to perform a regional tax normalization
// on the tax identity.
func (id *Identity) Calculate() error {
	r := id.Regime()
	if r != nil {
		return r.CalculateDocument(id)
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
		validation.Field(&id.Source, validation.In(validSourceKeys...)),
		validation.Field(&id.Meta),
	)
	if err != nil {
		return err
	}
	r := regimes.For(id.Country, id.Zone)
	if r != nil {
		return r.ValidateDocument(id)
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
	)
}

var validSourceKeys = generateValidSourceKeys()

func generateValidSourceKeys() []interface{} {
	ks := make([]interface{}, len(SourceKeyDefinitions))
	for i, v := range SourceKeyDefinitions {
		ks[i] = v.Key
	}
	return ks
}

// JSONSchema provides a representation of the struct for usage in Schema.
func (k SourceKey) JSONSchema() *jsonschema.Schema {
	s := &jsonschema.Schema{
		Title:       "Source Key",
		Type:        "string",
		OneOf:       make([]*jsonschema.Schema, len(SourceKeyDefinitions)),
		Description: "SourceKey identifies the source of a tax identity",
	}
	for i, v := range SourceKeyDefinitions {
		s.OneOf[i] = &jsonschema.Schema{
			Const:       cbc.Key(v.Key).String(),
			Description: v.Description,
		}
	}
	return s
}
