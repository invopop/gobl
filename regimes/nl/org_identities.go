package nl

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
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
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("OIN Stelsel - Samenstelling OIN"),
				URL:   "https://gitdocumentatie.logius.nl/publicatie/dk/oin/2.2.1/#samenstelling-oin",
			},
		},
	},
}

// oinPattern validates the OIN format: 6-zero prefix, 2-digit register code, 9-digit identifier, 3-zero suffix (20 digits).
var oinPattern = `^0{6}(0[1-9]|10|99)\d{9}0{3}$`

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				org.IdentityTypeIn(IdentityTypeKVK),
				rules.Field("code",
					rules.Assert("01", "identity code for type KVK must be valid",
						is.Length(8, 8),
					),
				),
			),
			rules.When(
				org.IdentityTypeIn(IdentityTypeOIN),
				rules.Field("code",
					rules.Assert("02", "identity code for type OIN must be valid",
						is.Matches(oinPattern),
					),
				),
			),
		),
	)
}
