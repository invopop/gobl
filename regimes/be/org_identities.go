package be

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

const (
	// IdentityTypeBCE represents the Belgian "Banque-Carrefour des Entreprises"
	// (Kruispuntbank van Ondernemingen, KBO) enterprise number used to identify
	// businesses in Belgium.
	IdentityTypeBCE cbc.Code = "BCE"
)

var identityDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeBCE,
		Name: i18n.String{
			i18n.EN: "BCE/KBO Number",
			i18n.FR: "Numéro BCE",
			i18n.NL: "KBO-nummer",
		},
	},
}
