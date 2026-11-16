package sk

import (
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

/*
	The Slovak IČO (Identifikačné číslo organizácie) is the number the
Statistical Office of the Slovak Republic assigns to every legal person and
entrepreneur. It is 8 digits long, the last one being a check digit
calculated with weights 8 down to 2 and a MOD 11 algorithm, the same as the
Czech IČO.

	Source: https://docs.peppol.eu/edelivery/codelists/v9.7/Peppol%20Code%20Lists%20-%20Participant%20identifier%20schemes%20v9.7.html
	(Participant identifier scheme 0158, SK:ICO)
*/

const (
	// IdentityTypeICO is the Slovak IČO
	IdentityTypeICO cbc.Code = "ICO"
)

// We use a regex to reject anything outside 8 plain digits.
var icoRegexp = regexp.MustCompile(`^\d{8}$`)

// Separators an IČO is written with, as in "36 396 567". Any other character
// stays and fails the format rule.
var icoSeparators = strings.NewReplacer(" ", "", " ", "", ".", "", "-", "", "/", "")

var identityTypeDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeICO,
		Name: i18n.String{
			i18n.EN: "Organisation Identification Number (IČO)",
			i18n.SK: "Identifikačné číslo organizácie (IČO)",
		},
	},
}

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				org.IdentityTypeIn(IdentityTypeICO),
				rules.Field("code",
					rules.AssertIfPresent("01", "identity code for type ICO must be 8 digits",
						is.MatchesRegexp(icoRegexp),
					),
					rules.AssertIfPresent("02", "identity code for type ICO failed the check digit",
						is.StringFunc("check digit", validateICOCheckDigit),
					),
				),
			),
		),
	)
}

func normalizeOrgIdentity(id *org.Identity) {
	if id == nil || id.Type != IdentityTypeICO {
		return
	}

	// Remove the separators only
	id.Code = cbc.Code(icoSeparators.Replace(id.Code.String()))
}

func validateICOCheckDigit(code string) bool {
	// A malformed code is left to the format rule
	if !icoRegexp.MatchString(code) {
		return true
	}

	// Weighted sum of the first 7 digits, weights 8 down to 2
	sum := 0
	for i, c := range code[:7] {
		sum += int(c-'0') * (8 - i)
	}

	// MOD 11 reduction: 11 minus the remainder, modulo 10
	return int(code[7]-'0') == (11-sum%11)%10
}
