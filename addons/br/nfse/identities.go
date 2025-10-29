package nfse

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// Identity keys
const (
	IdentityKeyMunicipalReg = "br-nfse-municipal-reg"
)

var identities = []*cbc.Definition{
	{
		Key: IdentityKeyMunicipalReg,
		Name: i18n.String{
			i18n.EN: "Company Municipal Registration",
			i18n.PT: "Inscrição Municipal da Empresa",
		},
	},
}
