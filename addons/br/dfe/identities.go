package dfe

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
)

// Identity keys
const (
	IdentityKeyMunicipalReg = "br-dfe-municipal-reg"
	IdentityKeyStateReg     = "br-dfe-state-reg"
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

	// migrate legacy addon keys
	for oldKey, newKey := range map[cbc.Key]cbc.Key{
		"br-nfse-municipal-reg": IdentityKeyMunicipalReg,
		"br-nfse-national-reg":  IdentityKeyStateReg,
	} {
		if i.Key == oldKey {
			i.Key = newKey
		}
	}
}
