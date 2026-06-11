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
				rules.AssertIfPresent("01", "invalid Norwegian VAT number",
					is.Func("valid MVA-suffixed mod-11 org number", isValidVATCode),
				),
			),
		),
	)
}

// normalizeTaxIdentity performs standard tax identity normalization, and then
// ensures the "MVA" suffix that Norwegian VAT numbers carry (e.g.
// "NO 923 456 783 MVA"): a VAT identity given as a bare organisation number
// gains the suffix, so the serialized form is always NO<orgnr>MVA as the
// EHF and Peppol national rules require.
func normalizeTaxIdentity(tID *tax.Identity) {
	tax.NormalizeIdentity(tID)
	if tID.Code != "" && !strings.HasSuffix(string(tID.Code), "MVA") {
		tID.Code += "MVA"
	}
}

// isValidVATCode reports whether the value is a valid Norwegian VAT number
// code: a mod-11 organisation number followed by the "MVA" suffix.
func isValidVATCode(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok || !strings.HasSuffix(string(code), "MVA") {
		return false
	}
	return isValidOrgNumber(cbc.Code(strings.TrimSuffix(string(code), "MVA")))
}

// isValidOrgNumber reports whether the value is a valid Norwegian
// organisasjonsnummer: nine digits with a mod-11 check digit. The leading-digit
// "8 or 9" range is an allocation convention, not a validation rule, so it is
// deliberately not enforced.
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
