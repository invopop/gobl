package es

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// The tax identity type is required for TicketBAI documents
// in the Basque Country.
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
		Map: cbc.CodeMap{
			KeyTicketBAIIDType: "03",
		},
	},
	{
		Key: IdentityKeyForeignID,
		Name: i18n.String{
			i18n.EN: "National ID Card or similar from a foreign country",
			i18n.ES: "Documento oficial de identificación expedido por el país o territorio de residencia",
		},
		Map: cbc.CodeMap{
			KeyTicketBAIIDType: "04",
		},
	},
	{
		Key: IdentityKeyResident,
		Name: i18n.String{
			i18n.EN: "Residential permit",
			i18n.ES: "Certificado de residencia",
		},
		Map: cbc.CodeMap{
			KeyTicketBAIIDType: "05",
		},
	},
	{
		Key: IdentityKeyOther,
		Name: i18n.String{
			i18n.EN: "An other type of source not listed",
			i18n.ES: "Otro documento probatorio",
		},
		Map: cbc.CodeMap{
			KeyTicketBAIIDType: "06",
		},
	},
}

var identityKeys = cbc.DefinitionKeys(identityKeyDefinitions)
