// Package ae provides the tax identity validation specific to the United Arab Emirates.
package ae

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pkg/luhn"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// TRNPattern is the format of the TRN: 15 digits, the first 10 of them the
// TIN (Ministry of Finance, UAE Electronic Invoice Mandatory Field
// Requirements, definition 11). Each tax registration, such as VAT or
// Corporate Tax, has its own TRN. PINT AE also requires a VAT TRN to end in
// 03 (IBR-132-AE), which the PINT addon checks. The first 10 digits end in
// a Luhn check digit, which the pattern does not check.
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

// tinCheckDigitIsValid reports whether the first 10 digits of a code, the
// TIN, end in a valid Luhn check digit, as the Federal Tax Authority
// registered the TIN under ISO 6523 (Peppol participant identifier scheme
// AE:TIN, ICD 0235). It is used on TINs and on TRNs, which start with one
// (see TRNPattern), after their format is checked.
func tinCheckDigitIsValid(val any) bool {
	code, ok := val.(cbc.Code)
	return ok && len(code) >= 10 && luhn.Check(code[:10])
}
