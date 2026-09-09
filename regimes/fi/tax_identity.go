package fi

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

var taxCodeRegexps = []*regexp.Regexp{
	regexp.MustCompile(`^\d{8}$`),
}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Finnish Y-tunnus tax identity code",
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

// Finland's Y-tunnus (Business ID) check digit validation.
//
// Format: 7 digits + hyphen + check digit. Hyphen removed during normalization.
// Validation: MOD 11 with weights [7, 9, 10, 5, 8, 4, 2, 1].
//
// Digit conversion via val[i]-'0' assumes the input contains only
// ASCII digits. This is guaranteed by the regex validation in validateTaxCode.
//
// Reference: https://www.vero.fi/globalassets/tietoa-verohallinnosta/ohjelmistokehittajille/yritys--ja-yhteisötunnuksen-ja-henkilötunnuksen-tarkistusmerkin-tarkistuslaskenta.pdf
func validateTaxCodeChecksum(val string) bool {
	weights := []int{7, 9, 10, 5, 8, 4, 2, 1}
	sum := 0

	for i := range 8 {
		sum += int(val[i]-'0') * weights[i]
	}

	return sum%11 == 0
}
