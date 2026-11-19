// Package ae provides the tax identity validation specific to the United Arab Emirates.
package ae

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pkg/luhn"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// TRNPattern is the TRN format: 15 digits starting with the TIN.
const TRNPattern = `^1\d{14}$`

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn("AE"),
			rules.Field("code",
				rules.AssertIfPresent("01", "UAE tax identity code must be 15 digits starting with 1",
					is.Matches(TRNPattern),
				),
				rules.When(
					is.Matches(TRNPattern),
					rules.AssertIfPresent("02", "UAE tax identity code must start with a TIN with a valid check digit",
						is.Func("valid TIN check digit", tinCheckDigitIsValid),
					),
				),
			),
		),
	)
}

// tinCheckDigitIsValid checks the Luhn digit of the TIN, the first 10 digits of a TIN or TRN.
func tinCheckDigitIsValid(val any) bool {
	code, ok := val.(cbc.Code)
	return ok && len(code) >= 10 && luhn.Check(code[:10])
}
