package es

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
)

var identityKeyDefinitions = []*cbc.KeyDefinition{
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
			i18n.EN: "National ID Card or similar from a foreign country",
			i18n.ES: "Documento oficial de identificación expedido por el país o territorio de residencia",
		},
	},
	{
		Key: org.IdentityKeyResident,
		Name: i18n.String{
			i18n.EN: "Residential permit",
			i18n.ES: "Certificado de residencia",
		},
	},
	{
		Key: org.IdentityKeyOther,
		Name: i18n.String{
			i18n.EN: "An other type of source not listed",
			i18n.ES: "Otro documento probatorio",
		},
	},
}
