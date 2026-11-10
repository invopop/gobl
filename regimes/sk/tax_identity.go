package sk

import (
	"regexp"
	"strconv"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// taxCodeRegexp matches a normalized Slovak IČ DPH: exactly 10 digits, the first
// of which is non-zero. The "SK" country prefix, spaces and hyphens are removed
// during normalization (tax.NormalizeIdentity), so only the digits reach this
// validation.
var taxCodeRegexp = regexp.MustCompile(`^[1-9]\d{9}$`)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "tax identity code for SK must be 10 digits starting with 1-9",
					is.MatchesRegexp(taxCodeRegexp),
				),
				rules.AssertIfPresent("02", "tax identity code for SK failed the checksum: the 10 digit number must be divisible by 11",
					is.StringFunc("checksum", validateTaxCodeChecksum),
				),
			),
		),
	)
}

// validateTaxCodeChecksum applies the Slovak IČ DPH checksum: the whole 10-digit
// number must be divisible by 11. Each assertion is evaluated independently, so a
// code that already failed the format rule still reaches this function; ParseInt
// rejects anything that is not a plain integer.
//
// Note: the Slovak tax authority does not publish this algorithm; the modulo-11
// rule is consistently documented by EU VAT-number validators and matches all
// known valid identifiers (e.g. SK2020273893).
//
// Reference: https://github.com/ltns35/go-vat/blob/main/countries/slovakia.go
// Reference: https://vatdb.com/guides/validate-sk-vat-number/
func validateTaxCodeChecksum(code string) bool {
	n, err := strconv.ParseInt(code, 10, 64)
	if err != nil {
		return false
	}
	return n%11 == 0
}
