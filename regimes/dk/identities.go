package dk

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

const (
	// IdentityKeyCVR is the type of identity that represents the Danish
	// "CVR-nummer" (Centrale Virksomhedsregister), the Central Business Register
	// number used to identify businesses in Denmark.
	IdentityKeyCVR cbc.Key = "cvr"
)

var identityKeyDefinitions = []*cbc.Definition{
	{
		Key: IdentityKeyCVR,
		Name: i18n.String{
			i18n.EN: "CVR Number",
			i18n.DA: "CVR-nummer",
		},
	},
}

// validateIdentity checks to ensure the CVR identity code is valid.
func validateIdentity(id *org.Identity) error {
	if id == nil || id.Key != IdentityKeyCVR {
		return nil
	}
	return validation.ValidateStruct(id,
		validation.Field(&id.Code, validation.By(validateTaxCode)),
	)
}
