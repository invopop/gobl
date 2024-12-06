package dian

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// Known base tax identity types for Colombia
const (
	IdentityKeyCivilRegister   cbc.Key = "co-civil-register"
	IdentityKeyIDCard          cbc.Key = "co-id-card"
	IdentityKeyCitizenID       cbc.Key = "co-citizen-id"
	IdentityKeyForeignerIDCard cbc.Key = "co-foreigner-id-card"
	IdentityKeyForeignerID     cbc.Key = "co-foreigner-id"
	IdentityKeyPassport        cbc.Key = "co-passport"
	IdentityKeyForeignID       cbc.Key = "co-foreign-id"
	IdentityKeyPEP             cbc.Key = "co-pep"
	IdentityKeyNUIP            cbc.Key = "co-nuip"
)

// Special keys to use in meta data.
const (
	KeyCompanyID cbc.Key = "dian-company-id"
)

var identities = []*cbc.Definition{
	{
		Key: IdentityKeyCivilRegister,
		Name: i18n.String{
			i18n.ES: "Registro Civil",
			i18n.EN: "Civil Registry",
		},
		Map: cbc.CodeMap{
			KeyCompanyID: "11",
		},
	},
	{
		Key: IdentityKeyIDCard,
		Name: i18n.String{
			i18n.EN: "Identity Card",
			i18n.ES: "Tarjeta de Identidad",
		},
		Map: cbc.CodeMap{
			KeyCompanyID: "12",
		},
	},
	{
		Key: IdentityKeyCitizenID,
		Name: i18n.String{
			i18n.EN: "Citizen Identity Card",
			i18n.ES: "Cédula de ciudadanía",
		},
		Map: cbc.CodeMap{
			KeyCompanyID: "13",
		},
	},
	{
		Key: IdentityKeyForeignerIDCard,
		Name: i18n.String{
			i18n.EN: "Foreigner Identity Card",
			i18n.ES: "Tarjeta de Extranjería",
		},
		Map: cbc.CodeMap{
			KeyCompanyID: "21",
		},
	},
	{
		Key: IdentityKeyForeignerID,
		Name: i18n.String{
			i18n.EN: "Foreigner Citizen Identity",
			i18n.ES: "Cédula de extranjería",
		},
		Map: cbc.CodeMap{
			KeyCompanyID: "22",
		},
	},
	{
		Key: IdentityKeyPassport,
		Name: i18n.String{
			i18n.EN: "Passport",
			i18n.ES: "Pasaporte",
		},
		Map: cbc.CodeMap{
			KeyCompanyID: "41",
		},
	},
	{
		Key: IdentityKeyForeignID,
		Name: i18n.String{
			i18n.EN: "Foreign Document",
			i18n.ES: "Documento de identificación extranjero",
		},
		Map: cbc.CodeMap{
			KeyCompanyID: "42",
		},
	},
	{
		Key: IdentityKeyPEP,
		Name: i18n.String{
			i18n.EN: "PEP - Special Permit to Stay",
			i18n.ES: "PEP - Permiso Especial de Permanencia",
		},
		Map: cbc.CodeMap{
			KeyCompanyID: "47",
		},
	},
	{
		Key: IdentityKeyNUIP,
		Name: i18n.String{
			i18n.EN: "NUIP - National Unique Personal Identification Number",
			i18n.ES: "NUIP - Número Único de Identificación Personal",
		},
		Map: cbc.CodeMap{
			KeyCompanyID: "91",
		},
	},
}

var identityKeys = cbc.DefinitionKeys(identities)
