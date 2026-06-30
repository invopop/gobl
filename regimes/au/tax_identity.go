package au

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// abnRegexp matches the basic shape of an Australian Business Number: 11 digits.
var abnRegexp = regexp.MustCompile(`^[0-9]{11}$`)

// abnWeights are the position weights used by the ABN check-digit algorithm.
var abnWeights = []int{10, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19}

// normalizeTaxIdentity strips whitespace, separators and country prefix from an
// Australian Business Number.
func normalizeTaxIdentity(tID *tax.Identity) {
	tax.NormalizeIdentity(tID)
}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.Assert("01", "invoice tax id code must be a valid 11-digit ABN",
					is.MatchesRegexp(abnRegexp)),
				rules.Assert("02", "invoice tax id code checksum must be valid",
					is.Func("valid", isValidTaxIdentityCode)),
			),
		),
	)
}

// isValidTaxIdentityCode reports whether the value is an ABN with valid check
// digits. Empty or non-conforming values are left to the format rule above.
func isValidTaxIdentityCode(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return true
	}
	return validABN(code.String())
}

// validABN reports whether code is a structurally valid ABN with correct check
// digits, using the ATO weighted modulus-89 algorithm.
func validABN(code string) bool {
	if len(code) != len(abnWeights) {
		return false
	}
	sum := 0
	for i, w := range abnWeights {
		d := int(code[i] - '0')
		if i == 0 {
			d-- // subtract 1 from the leading digit
		}
		sum += d * w
	}
	return sum%89 == 0
}
