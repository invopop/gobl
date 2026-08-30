package pe

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// taxCodeWeights are the mod-11 multipliers used to validate the check digit
// of a Peruvian RUC (Registro Único de Contribuyentes), as assigned by SUNAT.
// See: https://www.sunat.gob.pe/
var taxCodeWeights = []int{5, 4, 3, 2, 7, 6, 5, 4, 3, 2}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Peruvian RUC",
					is.Func("valid mod-11 RUC", isValidRUC),
				),
			),
		),
	)
}

// normalizeTaxIdentity performs the standard tax identity normalization,
// which removes punctuation and the "PE" country prefix, leaving the plain
// 11-digit RUC.
func normalizeTaxIdentity(tID *tax.Identity) {
	tax.NormalizeIdentity(tID)
}

// isValidRUC reports whether the value is a valid Peruvian RUC: eleven
// digits, starting with a taxpayer type prefix assigned by SUNAT, and ending
// with a mod-11 check digit.
func isValidRUC(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok {
		return false
	}
	s := code.String()
	if len(s) != 11 {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	// The two-digit prefix identifies the taxpayer type: 10 for natural
	// persons (embeds the DNI), 20 for legal entities, and 15, 16 and 17
	// for non-domiciled taxpayers and other special registrations.
	switch s[:2] {
	case "10", "15", "16", "17", "20":
	default:
		return false
	}

	sum := 0
	for i, w := range taxCodeWeights {
		sum += int(s[i]-'0') * w
	}
	check := 11 - (sum % 11)
	switch check {
	case 10:
		check = 0
	case 11:
		check = 1
	}
	return int(s[10]-'0') == check
}
