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
				rules.Assert("01", "VAT number must be 15 digits starting/ending with 3 (BR-KSA-40)",
					is.Func("when present, VAT id must be valid", vatIDValid)),
			),
		),
	)
}

func vatIDValid(val any) bool {
	code, ok := val.(cbc.Code)
	if !ok || code.IsEmpty() {
		return true
	}
	match, _ := regexp.MatchString(vatIDPattern, code.String())
	return match
}
