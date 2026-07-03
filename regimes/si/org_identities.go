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
	IdentityTypeMaticna cbc.Code = "MATICNA"
)

// registrationNumberLen is the length of a registration number: six
// digits, a modulo-11 check digit, and a three-digit suffix identifying
// the business unit ("000" for the head office).
const registrationNumberLen = 10

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
				Title: i18n.NewString("Uredba o vodenju in vzdrževanju Poslovnega registra Slovenije"),
				URL:   "https://pisrs.si/Pis.web/pregledPredpisa?id=URED7599",
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
					rules.Assert("01", "identity code for type MATICNA must be valid",
						is.Func("valid", isValidRegistrationNumber),
					),
				),
			),
		),
	)
}

// isValidRegistrationNumber expects a registration number in its full
// ten-digit form and validates the modulo-11 check digit that follows
// the six base digits. The business unit suffix carries no checksum of
// its own.
func isValidRegistrationNumber(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok || len(code) != registrationNumberLen {
		return false
	}
	for i := 7; i < registrationNumberLen; i++ {
		if code[i] < '0' || code[i] > '9' {
			return false
		}
	}
	return validateMod11(code[:7], registrationNumberMultipliers) == nil
}
