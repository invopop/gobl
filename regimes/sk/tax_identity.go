package sk

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// taxCodeRegexp matches a normalized Slovak IČ DPH: exactly 10 digits, the first
// of which is non-zero. The "SK" country prefix, spaces and hyphens are removed
// during normalization (tax.NormalizeIdentity), so only the digits reach this
// validation.
var taxCodeRegexp = regexp.MustCompile(`^[1-9]\d{9}$`)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "tax identity code for SK must be 10 digits starting with 1-9 and divisible by 11",
					is.Func("valid", validateTaxCode),
				),
			),
		),
	)
}

func validateTaxCode(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok {
		return false
	}
	val := code.String()

	if !taxCodeRegexp.MatchString(val) {
		return false
	}
	return validateTaxCodeChecksum(val)
}

// Slovak IČ DPH checksum: the whole 10-digit number must be divisible by 11.
//
// A 10-digit number can reach ~10^10, which overflows a 32-bit int, so the digits
// are accumulated into an int64. Digit conversion via val[i]-'0' assumes the input
// contains only ASCII digits, which is guaranteed by the regex in validateTaxCode.
//
// Note: the Slovak tax authority does not publish this algorithm; the modulo-11
// rule is consistently documented by EU VAT-number validators and matches all
// known valid identifiers (e.g. SK2020273893).
//
// Reference: https://github.com/ltns35/go-vat/blob/main/countries/slovakia.go
// Reference: https://vatdb.com/guides/validate-sk-vat-number/
func validateTaxCodeChecksum(val string) bool {
	var n int64
	for i := range len(val) {
		n = n*10 + int64(val[i]-'0')
	}
	return n%11 == 0
}
