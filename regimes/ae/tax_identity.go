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
					rules.AssertIfPresent("03", "UAE tax identity code must start with a TIN whose business identifier is not all zeros",
						is.Func("TIN business identifier set", tinBusinessIDIsSet),
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

// tinNoBusiness is a TIN's 8-digit business identifier, between the leading 1
// and the check digit, when it is all zeros, which the Federal Tax Authority
// never issues.
const tinNoBusiness = "00000000"

// tinBusinessIDIsSet reports whether the business identifier of the TIN at the
// start of a code is not all zeros. It is used on TINs and on TRNs, which
// start with one, after their format is checked.
func tinBusinessIDIsSet(val any) bool {
	code, ok := val.(cbc.Code)
	return ok && len(code) >= 10 && code[1:9] != tinNoBusiness
}
