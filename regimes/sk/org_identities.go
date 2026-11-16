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

const (
	// IdentityTypeICO is the IČO (identifikačné číslo organizácie), the number the
	// Statistical Office of the Slovak Republic assigns to every legal person and
	// entrepreneur, and which goes on commercial documents (§ 3a Commercial Code).
	IdentityTypeICO cbc.Code = "ICO"
)

// icoRegexp matches a normalized IČO: exactly 8 digits.
//
// Reference: Peppol participant identifier scheme 0158 (SK:ICO), code lists v9.7.
var icoRegexp = regexp.MustCompile(`^\d{8}$`)

// icoSeparators are the characters an IČO is written with, such as the spaces
// of "36 396 567". Any other character stays and fails the format rule.
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

// normalizeOrgIdentity removes the separators from an IČO, which is often
// written in groups of three digits such as "36 396 567".
func normalizeOrgIdentity(id *org.Identity) {
	if id == nil || id.Type != IdentityTypeICO {
		return
	}
	id.Code = cbc.Code(icoSeparators.Replace(id.Code.String()))
}

// validateICOCheckDigit applies the IČO check digit. A malformed code is left
// to the format rule.
//
// Format: 7 digits + check digit. Validation: MOD 11 with weights 8 down to 2,
// where the check digit is 11 minus the remainder, modulo 10. The IČO shares
// this algorithm with the Czech IČO (see regimes/cz).
//
// Reference: Peppol participant identifier scheme 0158 (SK:ICO), code lists v9.7,
// "Check number: mod11 (weights: 8, 7, 6, 5, 4, 3, 2)".
func validateICOCheckDigit(code string) bool {
	if !icoRegexp.MatchString(code) {
		return true
	}
	sum := 0
	for i, c := range code[:7] {
		sum += int(c-'0') * (8 - i)
	}
	return int(code[7]-'0') == (11-sum%11)%10
}
