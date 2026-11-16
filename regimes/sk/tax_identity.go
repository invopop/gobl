package sk

import (
	"regexp"
	"strconv"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// taxCodeRegexp matches a normalized Slovak DIČ: exactly 10 digits, the first of
// which gives the type of person, from 1 (natural person) to 7 (VAT group). An
// IČ DPH is the prefix "SK" followed by the same DIČ. The prefix, spaces and
// hyphens are removed during normalization (tax.NormalizeIdentity), so only the
// digits reach this validation.
//
// Reference: Peppol participant identifier scheme 0245 (SK:DIC), code lists v9.7.
var taxCodeRegexp = regexp.MustCompile(`^[1-7]\d{9}$`)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "tax identity code for SK must be 10 digits starting with 1-7",
					is.MatchesRegexp(taxCodeRegexp),
				),
				rules.AssertIfPresent("02", "tax identity code for SK failed the checksum: the 10 digit number must be divisible by 11",
					is.StringFunc("checksum", validateTaxCodeChecksum),
				),
			),
		),
	)
}

// validateTaxCodeChecksum applies the Slovak DIČ checksum: the whole 10-digit
// number must be divisible by 11. Each assertion is evaluated independently, so a
// code that already failed the format rule still reaches this function; ParseInt
// rejects anything that is not a plain integer.
//
// Reference: Peppol participant identifier scheme 0245 (SK:DIC), code lists v9.7,
// "The DIČ number is divisible by 11".
func validateTaxCodeChecksum(code string) bool {
	n, err := strconv.ParseInt(code, 10, 64)
	if err != nil {
		return false
	}
	return n%11 == 0
}
