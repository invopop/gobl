package ee

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

var taxCodeRegexps = []*regexp.Regexp{
	regexp.MustCompile(`^\d{9}$`),
}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Estonian VAT identity code",
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

	for _, re := range taxCodeRegexps {
		if re.MatchString(val) {
			return validateTaxCodeChecksum(val)
		}
	}
	return false
}

// Estonia's VAT number (käibemaksukohustuslase number)
// Finding any info on the algorithm was quite difficult,
// using what was stated in the docs for vat-validator
// Reference: https://vat-validator.readthedocs.io/en/latest/_modules/vat_validator/countries.html
// Format: EE + 9 digits
// Validation: weighted sum of first 8 digits with weights [3, 7, 1, 3, 7, 1, 3, 7]
//
// 9th digit is check digit, must equal (10 - (sum % 10)) % 10

func validateTaxCodeChecksum(val string) bool {
	weights := []int{3, 7, 1, 3, 7, 1, 3, 7}
	sum := 0

	for i := range 8 {
		sum += int(val[i]-'0') * weights[i]
	}

	expected := (10 - (sum % 10)) % 10

	return int(val[8]-'0') == expected
}
