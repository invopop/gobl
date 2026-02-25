package nl

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

const (
	// IdentityTypeKVK represents the Dutch "Kamer van Koophandel" (Chamber of Commerce)
	// registration number used to identify businesses in the Netherlands.
	IdentityTypeKVK cbc.Code = "KVK"

	// IdentityTypeOIN represents the Dutch "Organisatie Identificatie Nummer" used
	// to identify government organizations in the Netherlands.
	IdentityTypeOIN cbc.Code = "OIN"
)

var identityDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeKVK,
		Name: i18n.String{
			i18n.EN: "KVK Number",
			i18n.NL: "KVK-nummer",
		},
	},
	{
		Code: IdentityTypeOIN,
		Name: i18n.String{
			i18n.EN: "OIN Number",
			i18n.NL: "OIN-nummer",
		},
	},
}
