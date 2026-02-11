package dk

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

const (
	// IdentityTypeCVR represents the Danish "CVR-nummer" (Centrale Virksomhedsregister),
	// the Central Business Register number used to identify businesses in Denmark.
	IdentityTypeCVR cbc.Code = "CVR"
)

var identityTypeDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeCVR,
		Name: i18n.String{
			i18n.EN: "CVR Number",
			i18n.DA: "CVR-nummer",
		},
	},
}

// validateIdentity checks to ensure the CVR identity code is valid.
func validateIdentity(id *org.Identity) error {
	if id == nil || id.Type != IdentityTypeCVR {
		return nil
	}
	return validation.ValidateStruct(id,
		validation.Field(&id.Code, validation.By(validateTaxCode)),
	)
}
