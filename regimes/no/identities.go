package no

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

const (
	// IdentityTypeORG represents the Norwegian "Organisasjonsnummer"
	// (Organization Number) used to identify businesses in Norway.
	IdentityTypeORG cbc.Code = "ORG"
)

var identityTypeDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeORG,
		Name: i18n.String{
			i18n.EN: "Organization Number",
			i18n.NB: "Organisasjonsnummer",
		},
	},
}

func validateIdentity(id *org.Identity) error {
	if id == nil || id.Type != IdentityTypeORG {
		return nil
	}
	return validation.ValidateStruct(id,
		validation.Field(&id.Code, validation.By(validateTaxCode)),
	)
}
