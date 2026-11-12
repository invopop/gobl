// Package ae provides the tax identity validation specific to the United Arab Emirates.
package ae

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pkg/luhn"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

var (
	// trnRegexp is the format of the TRN: 15 digits, the first 10 of them the
	// TIN.
	trnRegexp = regexp.MustCompile(`^\d{15}$`)
	// tinRegexp is the format of the TIN, as the Federal Tax Authority
	// registered it under ISO 6523 (ICD 0235): 10 digits, the first one 1,
	// the UAE's identifier under the GCC agreement.
	tinRegexp = regexp.MustCompile(`^1\d{9}$`)
)

// tinNoBusiness is a TIN's 8-digit business identifier when it is all zeros,
// which the Federal Tax Authority never issues.
const tinNoBusiness = "00000000"

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn("AE"),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid UAE tax identity code format or check digit",
					is.Func("valid TRN", isValidTRN),
				),
			),
		),
	)
}

// isValidTRN reports whether a code is a TRN: 15 digits that start with a
// valid TIN. The last 5 digits tell the tax the number is for, and have no
// published list.
func isValidTRN(val any) bool {
	code, ok := val.(cbc.Code)
	return ok && trnRegexp.MatchString(code.String()) && validTIN(code[:10])
}

// validTIN reports whether a code is a TIN: 10 digits starting with 1, with a
// business identifier that is not all zeros and a Luhn check digit at the end.
func validTIN(code cbc.Code) bool {
	return tinRegexp.MatchString(code.String()) &&
		code[1:9] != tinNoBusiness &&
		luhn.Check(code)
}
