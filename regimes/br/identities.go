package br

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
)

// Identity keys
const (
	IdentityKeyMunicipalReg = "br-municipal-reg"
	IdentityKeyStateReg     = "br-state-reg"
)

var identities = []*cbc.Definition{
	{
		Key: IdentityKeyMunicipalReg,
		Name: i18n.String{
			i18n.EN: "Company Municipal Registration",
			i18n.PT: "Inscrição Municipal da Empresa",
		},
	},
	{
		Key: IdentityKeyStateReg,
		Name: i18n.String{
			i18n.EN: "Company State Registration",
			i18n.PT: "Inscrição Estadual da Empresa",
		},
	},
}

func normalizeIdentity(i *org.Identity) {
	if i == nil {
		return
	}

	// migrate old addon keys to the regime
	switch i.Key {
	case "br-nfse-municipal-reg":
		i.Key = IdentityKeyMunicipalReg
	case "br-nfse-national-reg":
		i.Key = IdentityKeyStateReg
	}
}
