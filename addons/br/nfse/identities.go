package nfse

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// Brazilian identity keys required to issue NFS-e documents.
const (
	IdentityKeyMunicipalReg = "br-nfse-municipal-reg"
	IdentityKeyNationalReg  = "br-nfse-national-reg"
)

var identities = []*cbc.KeyDefinition{
	{
		Key: IdentityKeyMunicipalReg,
		Name: i18n.String{
			i18n.EN: "Company Municipal Registration",
			i18n.PT: "Inscrição Municipal da Empresa",
		},
	},
	{
		Key: IdentityKeyNationalReg,
		Name: i18n.String{
			i18n.EN: "Company National Registration",
			i18n.PT: "Inscrição Nacional da Empresa",
		},
	},
}
