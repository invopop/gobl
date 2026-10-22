package hr

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// IdentityTypeOIB represents the Croatian "Osobni identifikacijski broj" (OIB),
	// the Personal Identification Number assigned in Croatia.
	IdentityTypeOIB cbc.Code = "OIB"
)

var identityTypeDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeOIB,
		Name: i18n.String{
			i18n.EN: "Personal Identification Number",
			i18n.HR: "Osobni identifikacijski broj",
		},
	},
}

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				org.IdentityTypeIn(IdentityTypeOIB),
				rules.Field("code",
					rules.AssertIfPresent("01", "Croatian OIB identity code must be 11 digits",
						is.MatchesRegexp(oibRegexp),
					),
					rules.AssertIfPresent("02", "Croatian OIB identity code failed the ISO 7064 MOD 11,10 checksum",
						is.StringFunc("checksum", validOIBChecksum),
					),
				),
			),
		),
	)
}
