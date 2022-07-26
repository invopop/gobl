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

// Validate checks to ensure the tax ID contains all the required
// fields. The check the value itself is in the expected format according
// to the country, you'll need to use the region packages directly. See also
// the region `ValidateTaxID` method.
func (id *TaxIdentity) Validate() error {
	return validation.ValidateStruct(id,
		validation.Field(&id.UUID),
		validation.Field(&id.Country, validation.Required),
		validation.Field(&id.Locality),
		validation.Field(&id.Source, validation.In(validSourceKeys()...)),
		validation.Field(&id.Meta),
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
