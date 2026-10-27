package hr

import (
	"strings"

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

// normalizeOrgIdentity cleans an OIB the way the tax ID is cleaned: it drops
// separators and the HR prefix of the VAT number, which shares the OIB's
// digits, e.g. "HR 1234-5678-903" becomes "12345678903". Any other letter is
// kept, so that the rules flag the code rather than a typo becoming valid.
func normalizeOrgIdentity(id *org.Identity) {
	if id.Type != IdentityTypeOIB {
		return
	}
	code := cbc.NormalizeAlphanumericalCode(id.Code)
	id.Code = cbc.Code(strings.TrimPrefix(code.String(), CountryCode))
}
