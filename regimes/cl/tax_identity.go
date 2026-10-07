package cl

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// rutWeights are the mod-11 multipliers used to calculate a Chilean RUT's
// check digit ("dígito verificador"), applied right-to-left and repeating.
var rutWeights = []int{2, 3, 4, 5, 6, 7}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Chilean RUT",
					is.Func("valid mod-11 RUT", isValidRUT),
				),
			),
		),
	)
}

// normalizeTaxIdentity performs the standard tax identity normalization,
// which removes punctuation (dots and the hyphen) and the "CL" country
// prefix, leaving the plain RUT body plus its uppercase check digit.
func normalizeTaxIdentity(tID *tax.Identity) {
	tax.NormalizeIdentity(tID)
}

// isValidRUT reports whether the value is a valid Chilean RUT: a numeric
// body of up to 8 digits followed by a mod-11 check digit (0-9 or K).
func isValidRUT(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok {
		return false
	}
	s := code.String()
	// Chilean law does not fix an exact digit count for the RUT body. Numbers
	// have been assigned sequentially for decades, so older RUTs can be shorter
	// than the 7-8 digits typical of RUTs issued today. We accept a body of
	// 1 to 8 digits plus 1 check digit (total length 2-9) as a practical bound,
	// not a legal requirement, to avoid rejecting older but mathematically
	// valid RUTs.
	if len(s) < 2 || len(s) > 9 {
		return false
	}
	body, dv := s[:len(s)-1], s[len(s)-1]
	for _, r := range body {
		if r < '0' || r > '9' {
			return false
		}
	}
	sum := 0
	count := 0
	for i := len(body) - 1; i >= 0; i-- {
		digit := int(body[i] - '0')
		weight := rutWeights[count%len(rutWeights)]
		sum += digit * weight
		count++
	}
	check := 11 - (sum % 11)
	switch check {
	case 11:
		check = 0
	case 10:
		return dv == 'K'
	}
	return int(dv-'0') == check
}
