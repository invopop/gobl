package cy

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// TINs issued since 27 March 2023 start with 6. Older TINs remain valid,
// so the prefix cannot be enforced without knowing the registration date.
// Source: Cyprus Tax Department.
// https://www.gov.cy/oikonomia/anakoinosi-tou-tmimatos-forologias-anaforika-me-ti-nea-morfi-arithmou-forologikis-taftotitas-a-f-t/
var taxCodeRegexp = regexp.MustCompile(`^\d{8}[A-Z]$`)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Cyprus TIN identity code",
					is.Func("valid", isValidTaxIdentityCode),
				),
			),
		),
	)
}

func isValidTaxIdentityCode(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return false
	}
	return taxCodeRegexp.MatchString(code.String())
}
