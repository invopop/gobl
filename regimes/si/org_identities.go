package si

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// IdentityTypeMaticna represents the Slovenian registration number
	// (matična številka) assigned by AJPES to every entity entered in
	// the Slovenian Business Register (e.g. "5043611000").
	IdentityTypeMaticna cbc.Code = "MŠ"
)

// registrationNumberPattern matches a registration number (matična
// številka): six base digits, a modulo-11 check digit, and a three-digit
// suffix identifying the business unit ("000" for the head office).
const registrationNumberPattern = `^\d{10}$`

// registrationNumberMultipliers are the weights applied to the first six
// digits of the registration number to calculate the check digit.
//
// Algorithm source: Article 11 of the Slovenian Business Register decree
// (Uredba o vodenju in vzdrževanju Poslovnega registra Slovenije).
// https://pisrs.si/Pis.web/pregledPredpisa?id=URED7599
var registrationNumberMultipliers = []int{7, 6, 5, 4, 3, 2}

var identityDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeMaticna,
		Name: i18n.String{
			i18n.EN: "Registration Number",
			i18n.SL: "Matična številka",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Decree on the keeping and maintenance of the Slovenian Business Register",
					i18n.SL: "Uredba o vodenju in vzdrževanju Poslovnega registra Slovenije",
				},
				URL: "https://pisrs.si/Pis.web/pregledPredpisa?id=URED7599",
			},
		},
	},
}

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				org.IdentityTypeIn(IdentityTypeMaticna),
				rules.Field("code",
					rules.Assert("01", "Slovenian registration number is required",
						is.Present,
					),
					rules.AssertIfPresent("02", "invalid Slovenian registration number format",
						is.Matches(registrationNumberPattern),
					),
					rules.AssertIfPresent("03", "invalid Slovenian registration number check digit",
						is.StringFunc("checksum", registrationNumberChecksumValid),
					),
				),
			),
		),
	)
}

// registrationNumberChecksumValid reports whether the registration number's
// seventh digit is a valid modulo-11 check digit of the first six. The
// three-digit business-unit suffix carries no checksum of its own and is
// only validated for being numeric by the format rule. It guards its length
// so it is safe to call before the format rule has run.
func registrationNumberChecksumValid(code string) bool {
	if len(code) < 7 {
		return false
	}
	return validMod11(code[:7], registrationNumberMultipliers)
}
