package cr

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Costa Rican tax identities reuse the party's identity document and carry no
// check digit, so validation is structural ("Anexos y Estructuras v4.4", campo
// "Número de cédula" and the OECD TIN profile for Costa Rica):
//
//   - Cédula física: 9 digits, no leading zero.
//   - Cédula jurídica or NITE: 10 characters (v4.4 permits letters).
//   - DIMEX: 11 or 12 digits, no leading zero.

func normalizeTaxIdentity(tID *tax.Identity) {
	tax.NormalizeIdentity(tID)
}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "must be 9 to 12 characters long",
					is.Length(9, 12),
				),
				rules.AssertIfPresent("02", "must contain only digits, or digits and letters when 10 characters long",
					codeTest("valid characters", hasValidCharacters),
				),
				rules.AssertIfPresent("03", "must not start with a zero",
					codeTest("no leading zero", hasNoLeadingZero),
				),
			),
		),
	)
}

// codeTest adapts a check on the code's string form into a rule test,
// converting the field value once for all assertions that use it.
func codeTest(name string, fn func(string) bool) is.FuncTest {
	return is.Func(name, func(value any) bool {
		code, ok := value.(cbc.Code)
		if !ok || code == "" {
			return false
		}
		return fn(code.String())
	})
}

// hasValidCharacters checks the content for the format implied by the length:
// 10-character codes (cédula jurídica or NITE) may combine digits and letters,
// every other format is digits only. Codes with an invalid length are left to
// the length rule.
func hasValidCharacters(s string) bool {
	switch len(s) {
	case 9, 11, 12:
		return isDigits(s)
	case 10:
		return isAlphanumeric(s)
	}
	return true
}

// hasNoLeadingZero rejects leading zeros on cédula física and DIMEX codes, the
// two formats the v4.4 annex forbids them for.
func hasNoLeadingZero(s string) bool {
	switch len(s) {
	case 9, 11, 12:
		return s[0] != '0'
	}
	return true
}

func isDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func isAlphanumeric(s string) bool {
	for _, r := range s {
		if (r < '0' || r > '9') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}
