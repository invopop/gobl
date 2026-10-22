package hr

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

const (
	// IdentityTypeOIB represents the Croatian "Osobni identifikacijski broj" (OIB),
	// the Personal Identification Number assigned in Croatia.
	IdentityTypeOIB cbc.Code = "OIB"
)

var identityTypeDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeOIB,
		Name: i18n.String{
			i18n.EN: "Personal Identification Number",
			i18n.HR: "Osobni identifikacijski broj",
		},
	},
}

func validateIdentity(id *org.Identity) error {
	if id == nil || id.Type != IdentityTypeOIB {
		return nil
	}
	return validation.ValidateStruct(id,
		validation.Field(&id.Code, validation.By(validateTaxCode)),
	)
}
