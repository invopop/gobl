package lu

import (
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "Luxembourg TVA number must be exactly 8 digits",
					is.Length(8, 8),
				),
				rules.AssertIfPresent("02", "Luxembourg TVA number must contain only digits",
					is.Digit,
				),
				rules.AssertIfPresent("03", "invalid Luxembourg TVA number: mod-89 checksum mismatch",
					is.Func("valid mod-89 checksum", tvaChecksumValid),
				),
			),
		),
	)
}

// normalizeTaxIdentity strips whitespace, dashes, and the "LU" country prefix
// from a Luxembourg TVA number, leaving the raw 8-digit code.
func normalizeTaxIdentity(tID *tax.Identity) {
	tax.NormalizeIdentity(tID)
}

// tvaChecksumValid reports whether the last two digits of a Luxembourg TVA
// number equal the first six digits mod 89.
//
// Source: https://www.aed.public.lu/en/tva/numero-tva.html
func tvaChecksumValid(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok {
		return true
	}
	val := code.String()
	if len(val) != 8 {
		return true // length rule handles this
	}
	base, err1 := strconv.Atoi(val[:6])
	check, err2 := strconv.Atoi(val[6:])
	if err1 != nil || err2 != nil {
		return true // digit rule handles this
	}
	return base%89 == check
}
