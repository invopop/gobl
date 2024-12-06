package nfse

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// Brazilian identity keys required to issue NFS-e documents. In an initial
// assessment, these identities do not seem to apply to documents other than
// NFS-e. However, if when implementing other Fiscal Notes it is found that some
// of these extensions are common, they can be moved to the regime or to a
// shared addon.
const (
	IdentityKeyMunicipalReg = "br-nfse-municipal-reg"
	IdentityKeyNationalReg  = "br-nfse-national-reg"
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
		Key: IdentityKeyNationalReg,
		Name: i18n.String{
			i18n.EN: "Company National Registration",
			i18n.PT: "Inscrição Nacional da Empresa",
		},
	},
}
