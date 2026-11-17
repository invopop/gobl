// Package ae provides the tax identity validation specific to the United Arab Emirates.
package ae

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pkg/luhn"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// trnPattern is the format of the TRN: 15 digits, the first 10 of them the
// TIN. The last 5 tell the tax the number is for, and have no published list.
const trnPattern = `^1\d{14}$`

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn("AE"),
			rules.Field("code",
				rules.AssertIfPresent("01", "UAE tax identity code must be 15 digits starting with 1",
					is.Matches(trnPattern),
				),
				rules.When(
					is.Matches(trnPattern),
					rules.AssertIfPresent("02", "UAE tax identity code must start with a TIN with a valid check digit",
						is.Func("valid TIN check digit", tinCheckDigitIsValid),
					),
				),
			),
		),
	)
}

// tinCheckDigitIsValid reports whether the first 10 digits of a code, the
// TIN, end in a valid Luhn check digit. It is used on TINs and on TRNs, which
// start with one, after their format is checked.
func tinCheckDigitIsValid(val any) bool {
	code, ok := val.(cbc.Code)
	return ok && len(code) >= 10 && luhn.Check(code[:10])
}
