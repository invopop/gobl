package org

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/uuid"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// TaxID represents a party's tax identify number for a given
// country.
type TaxID struct {
	// Unique universal identity code
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// ISO country code for Where the tax identity was issued
	Country l10n.Country `json:"country" jsonschema:"title=Country Code"`
	// Identity Code
	Code string `json:"code" jsonschema:"title=Code"`
	// Additional details.
	Meta Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate checks to ensure the tax ID contains all the required
// fields. The check the value itself is in the expected format according
// to the country, you'll need to use the region packages directly. See
// the region `ValidateTaxID` method.
func (tid *TaxID) Validate() error {
	return validation.ValidateStruct(tid,
		validation.Field(&tid.Country, validation.Required),
		validation.Field(&tid.Code, validation.Required),
	)
}
