package org

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// SourceKey is used to identify different sources of tax
// identities that may be required by some regions.
type SourceKey Key

// DefSourceKey defines the details we have regarding a document
// source key.
type DefSourceKey struct {
	Key         SourceKey `json:"key" jsonschema:"title=Key"`
	Description string    `json:"description" jsonschema:"title=Description"`
}

// RequireTaxIdentityCode is an additional check to use alongside
// regular validation that will ensure the tax ID has a code
// value set.
var RequireTaxIdentityCode = validateTaxID{requireCode: true}

type validateTaxID struct {
	requireCode bool
}

var (
	regionTaxIDValidation func(tID *TaxIdentity) error
	regionTaxIDNormalizer func(tID *TaxIdentity) error
)

// SetTaxIdentityValidation will prepare a reference to the tax ID regional
// validator. This is an internal method and will panic if called more than once.
func SetTaxIdentityValidation(cb func(tID *TaxIdentity) error) {
	if regionTaxIDValidation != nil {
		panic("tax identity regional validation function already set")
	}
	regionTaxIDValidation = cb
}

// SetTaxIdentityNormalizer will prepare a reference to the tax ID regional
// cleaner. This is an internal method and will panic if called more than once.
func SetTaxIdentityNormalizer(cb func(tID *TaxIdentity) error) {
	if regionTaxIDNormalizer != nil {
		panic("tax identity regional cleaner function already set")
	}
	regionTaxIDNormalizer = cb
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

// TaxIdentity stores the details required to identify an entity for tax
// purposes.
type TaxIdentity struct {
	// Unique universal identity code.
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`

	// ISO country code for Where the tax identity was issued.
	Country l10n.CountryCode `json:"country" jsonschema:"title=Country Code"`

	// Where inside a country the Tax ID was issued, if required.
	Locality l10n.Code `json:"locality,omitempty" jsonschema:"title=Locality Code"`

	// What is the source document of this tax identity.
	Source SourceKey `json:"source,omitempty" jsonschema:"title=Source Key"`

	// Tax identity Code
	Code string `json:"code,omitempty" jsonschema:"title=Code"`

	// Additional details that may be required.
	Meta Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Calculate will attempt to perform a regional tax normalization
// on the tax identity.
func (id *TaxIdentity) Calculate() error {
	return regionTaxIDNormalizer(id)
}

// Validate checks to ensure the tax ID contains all the required
// fields. The check the value itself is in the expected format according
// to the country, you'll need to use the region packages directly. See also
// the region `ValidateTaxID` method.
func (id *TaxIdentity) Validate() error {
	err := validation.ValidateStruct(id,
		validation.Field(&id.UUID),
		validation.Field(&id.Country, validation.Required),
		validation.Field(&id.Locality),
		validation.Field(&id.Source, validation.In(validSourceKeys()...)),
		validation.Field(&id.Meta),
	)
	if err != nil {
		return err
	}
	if regionTaxIDValidation != nil {
		return regionTaxIDValidation(id)
	}
	return nil
}

func (v validateTaxID) Validate(value interface{}) error {
	id, ok := value.(*TaxIdentity)
	if !ok {
		return nil
	}
	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.When(v.requireCode, validation.Required),
		),
	)
}

func validSourceKeys() []interface{} {
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
			Const:       Key(v.Key).String(),
			Description: v.Description,
		}
	}
	return s
}
