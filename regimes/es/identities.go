package es

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// Spain recognises a certain set of identity types specifically that may be mapped
// to local keys.
const (
	IdentityKeyPassport  cbc.Key = "es-passport"
	IdentityKeyForeignID cbc.Key = "es-foreign-id"
	IdentityKeyResident  cbc.Key = "es-resident"
	IdentityKeyOther     cbc.Key = "es-other"
)

var identityKeyDefinitions = []*cbc.KeyDefinition{
	{
		Key: IdentityKeyPassport,
		Name: i18n.String{
			i18n.EN: "Passport",
			i18n.ES: "Pasaporte",
		},
	},
	{
		Key: IdentityKeyForeignID,
		Name: i18n.String{
			i18n.EN: "National ID Card or similar from a foreign country",
			i18n.ES: "Documento oficial de identificación expedido por el país o territorio de residencia",
		},
	},
	{
		Key: IdentityKeyResident,
		Name: i18n.String{
			i18n.EN: "Residential permit",
			i18n.ES: "Certificado de residencia",
		},
	},
	{
		Key: IdentityKeyOther,
		Name: i18n.String{
			i18n.EN: "An other type of source not listed",
			i18n.ES: "Otro documento probatorio",
		},
	},
}
