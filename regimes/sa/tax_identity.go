package sa

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Saudi VAT Identification Number is a 15-digit number starting and ending with 3.
const (
	vatIDPattern = `^3[0-9]{13}3$`
)

var (
	vatIDRegex = regexp.MustCompile(vatIDPattern)
)

func normalizeTaxIdentity(tID *tax.Identity) {
	if tID.Code == "" {
		return
	}
	tax.NormalizeIdentity(tID)
}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(countryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Saudi tax identity code",
					is.Func("valid Saudi VAT number", validateTaxCode),
				),
			),
		),
	)
}

func validateTaxCode(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return false
	}
	return vatIDRegex.MatchString(code.String())
}
