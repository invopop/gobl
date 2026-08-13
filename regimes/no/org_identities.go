package no

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// IdentityTypeOrgNr represents the Norwegian "organisasjonsnummer",
	// the 9-digit number assigned by Brønnøysundregistrene to identify
	// businesses in Norway.
	IdentityTypeOrgNr cbc.Code = "ON"
)

var identityTypeDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeOrgNr,
		Name: i18n.String{
			i18n.EN: "Organization Number",
			i18n.NB: "Organisasjonsnummer",
		},
		Desc: i18n.String{
			i18n.EN: "Norwegian organization number assigned by Brønnøysundregistrene.",
			i18n.NB: "Norsk organisasjonsnummer tildelt av Brønnøysundregistrene.",
		},
	},
}

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				org.IdentityTypeIn(IdentityTypeOrgNr),
				rules.Field("code",
					rules.Assert("01", "invalid organisasjonsnummer",
						is.Func("valid mod-11 org number", isValidOrgNumber),
					),
				),
			),
		),
	)
}

// normalizeOrgIdentity strips non-numeric characters from the organization number.
func normalizeOrgIdentity(id *org.Identity) {
	if id.Type != IdentityTypeOrgNr {
		return
	}
	id.Code = cbc.NormalizeNumericalCode(id.Code)
}
