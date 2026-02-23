package no

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

const (
	// IdentityTypeOrgNr represents the Norwegian "organisasjonsnummer",
	// the 9-digit number assigned by Brønnøysundregistrene to identify
	// businesses in Norway.
	IdentityTypeOrgNr cbc.Code = "ON"
)

var identityTypeDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeOrgNr,
		Name: i18n.String{
			i18n.EN: "Organization Number",
			i18n.NB: "Organisasjonsnummer",
		},
		Desc: i18n.String{
			i18n.EN: "Norwegian organization number assigned by Brønnøysundregistrene.",
			i18n.NB: "Norsk organisasjonsnummer tildelt av Brønnøysundregistrene.",
		},
	},
}

// normalizeOrgIdentity strips non-numeric characters from the organization number.
func normalizeOrgIdentity(id *org.Identity) {
	if id == nil || id.Type != IdentityTypeOrgNr {
		return
	}
	id.Code = cbc.NormalizeNumericalCode(id.Code)
}

// validateOrgIdentity checks the Norwegian organisasjonsnummer using the same
// mod-11 algorithm as the tax identity. The organisation number is the base
// for the VAT number (org.nr + "MVA").
func validateOrgIdentity(id *org.Identity) error {
	if id == nil || id.Type != IdentityTypeOrgNr {
		return nil
	}
	return validation.ValidateStruct(id,
		validation.Field(&id.Code, validation.By(validateTaxCode)),
	)
}
