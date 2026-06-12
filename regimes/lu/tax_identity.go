package lu

import (
	"errors"
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
				rules.AssertIfPresent("01", "invalid Luxembourg TVA number",
					is.Func("valid mod-89 TVA code", isValidTVACode),
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

// isValidTVACode reports whether the value is a valid Luxembourg TVA number:
// exactly 8 digits where the last two digits equal the first six digits mod 89.
//
// Source: https://www.aed.public.lu/en/tva/numero-tva.html
func isValidTVACode(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok {
		return false
	}
	return validateTVACode(code) == nil
}

func validateTVACode(code cbc.Code) error {
	val := code.String()
	if len(val) != 8 {
		return errors.New("must be exactly 8 digits")
	}
	for _, r := range val {
		if r < '0' || r > '9' {
			return errors.New("must contain only digits")
		}
	}
	base, _ := strconv.Atoi(val[:6])
	check, _ := strconv.Atoi(val[6:])
	if base%89 != check {
		return errors.New("checksum mismatch")
	}
	return nil
}
