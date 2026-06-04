package no

import (
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// taxCodeWeights are the mod-11 multipliers for Norwegian organisasjonsnummer
// validation, as specified by Brønnøysundregistrene.
// See: https://www.brreg.no/en/about-us-2/our-registers/about-the-central-coordinating-register-for-legal-entities-ccr/about-the-organisation-number/
var taxCodeWeights = []int{3, 2, 7, 6, 5, 4, 3, 2}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid organisasjonsnummer",
					is.Func("valid mod-11 org number", isValidOrgNumber),
				),
			),
		),
	)
}

// normalizeTaxIdentity performs standard tax identity normalization, and then
// removes the "MVA" suffix common in Norwegian VAT numbers
// (e.g. "NO 923 456 783 MVA").
func normalizeTaxIdentity(tID *tax.Identity) {
	if tID == nil {
		return
	}
	tax.NormalizeIdentity(tID)
	tID.Code = cbc.Code(strings.TrimSuffix(string(tID.Code), "MVA"))
}

// isValidOrgNumber reports whether the value is a valid Norwegian
// organisasjonsnummer: nine digits validated by a mod-11 check digit. The same
// number is the basis for both the tax identity and the `ON` organization
// identity (org.nr + "MVA" forms the VAT number).
//
// Per Brønnøysundregistrene the only structural rule is the mod-11 check digit;
// the leading-digit "8 or 9" range is an allocation convention, not a
// validation rule, so we deliberately do not enforce it (it would reject
// otherwise-valid numbers allocated outside that range).
func isValidOrgNumber(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok {
		return false
	}
	s := code.String()
	if len(s) != 9 {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	sum := 0
	for i, w := range taxCodeWeights {
		sum += int(s[i]-'0') * w
	}
	remainder := sum % 11
	check := 0
	if remainder != 0 {
		check = 11 - remainder
	}
	if check == 10 {
		return false
	}
	return int(s[8]-'0') == check
}
