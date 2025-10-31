package ar

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
)

// identityDefinitions provides the list of identity types supported in Argentina
// beyond the tax identification (CUIT/CUIL/CDI).
func identityDefinitions() []*cbc.Definition {
	return []*cbc.Definition{
		{
			Key: org.IdentityKeyPassport,
			Name: i18n.String{
				i18n.EN: "Passport",
				i18n.ES: "Pasaporte",
			},
		},
		{
			Key: org.IdentityKeyForeign,
			Name: i18n.String{
				i18n.EN: "Foreign National ID",
				i18n.ES: "Documento de Identidad Extranjero",
			},
		},
		{
			Key: org.IdentityKeyResident,
			Name: i18n.String{
				i18n.EN: "Residential Permit",
				i18n.ES: "Certificado de Residencia",
			},
		},
		{
			Key: "dni",
			Name: i18n.String{
				i18n.EN: "National Identity Document",
				i18n.ES: "Documento Nacional de Identidad (DNI)",
			},
		},
		{
			Key: org.IdentityKeyOther,
			Name: i18n.String{
				i18n.EN: "Other identification document",
				i18n.ES: "Otro documento de identificación",
			},
		},
	}
}
