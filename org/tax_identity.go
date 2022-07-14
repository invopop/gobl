package org

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/uuid"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Main DocumentCode definitions.
const (
	DocumentKeyCompany  Key = ""         // Tax Authority
	DocumentKeyPassport Key = "passport" // A passport document
	DocumentKeyNational Key = "national" // National ID Card or similar
	DocumentKeyPermit   Key = "permit"   // Residential permit
	DocumentKeyOther    Key = "other"    // Something else
)

// TaxIdentity stores the details required to identify an entity for tax
// purposes.
type TaxIdentity struct {
	// Unique universal identity code.
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`

	// ISO country code for Where the tax identity was issued.
	Country l10n.Code `json:"country" jsonschema:"title=Country Code"`

	// Where inside a country the Tax ID was issued, if required.
	Locality l10n.Code `json:"locality,omitempty" jsonschema:"title=Locality Code"`

	// What is the source document of this tax identity.
	Document Key `json:"document,omitempty" jsonschema:"title=Document Type Key"`

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
		validation.Field(&id.Country, validation.Required, l10n.IsCountry),
		validation.Field(&id.Locality),
		validation.Field(&id.Document),
		validation.Field(&id.Meta),
	)
}
