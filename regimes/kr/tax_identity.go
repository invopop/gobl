package kr

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// brnWeights are the multipliers applied to the first eight digits of a Korean
// Business Registration Number (사업자등록번호) when validating its check digit.
var brnWeights = []int{1, 3, 7, 1, 3, 7, 1, 3}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Korean business registration number",
					is.Func("valid BRN check digit", isValidBRN),
				),
			),
		),
	)
}

// isValidBRN reports whether the value is a valid Korean Business Registration
// Number: ten digits ending in the check digit defined by the National Tax
// Service. The first eight digits are weighted by brnWeights; the ninth digit
// is multiplied by five and the tens and units of that product are added
// separately; the tenth digit is the check digit that makes the total a
// multiple of ten.
func isValidBRN(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok {
		return false
	}
	s := code.String()
	if len(s) != 10 {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	sum := 0
	for i, w := range brnWeights {
		sum += int(s[i]-'0') * w
	}
	ninth := int(s[8]-'0') * 5
	sum += ninth/10 + ninth%10
	check := (10 - sum%10) % 10
	return check == int(s[9]-'0')
}
