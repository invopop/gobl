package lt

import (
	"regexp"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// taxCodeRegexp matches a normalized Lithuanian VAT payer code (PVM mokėtojo
// kodas): 9 digits with 1 as the 8th digit, or 12 digits with 1 as the 11th
// digit. The "LT" country prefix, spaces and hyphens are removed during
// normalization (tax.NormalizeIdentity), so only the digits reach this
// validation.
var taxCodeRegexp = regexp.MustCompile(`^(\d{7}1\d|\d{10}1\d)$`)

// taxCodeWeights and taxCodeAltWeights are the multipliers of the first and
// second check digit passes, applied in order to every digit except the last.
var (
	taxCodeWeights    = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2}
	taxCodeAltWeights = []int{3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4}
)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "tax identity code for LT must be 9 digits with 1 as the 8th digit, or 12 digits with 1 as the 11th digit",
					is.MatchesRegexp(taxCodeRegexp),
				),
				rules.AssertIfPresent("02", "tax identity code for LT failed the checksum",
					is.StringFunc("checksum", validateTaxCodeChecksum),
				),
			),
		),
	)
}

// validateTaxCodeChecksum applies the Lithuanian VAT payer code check digit:
// the weighted sum of every digit except the last, modulo 11, must equal the
// last digit. If the remainder is 10, the sum is computed again with the
// alternative weights, and if it is still 10 the check digit is 0. Each
// assertion is evaluated independently, so a code that already failed the
// format rule still reaches this function, which rejects anything that is not
// 9 or 12 digits.
//
// Note: Art. 74 of the Law on VAT No. IX-751 leaves the composition of the code
// to the State Tax Inspectorate (VMI), which does not publish this algorithm. It
// follows two independent implementations and has been verified against real
// taxpayer codes.
//
// Reference: https://github.com/arthurdejong/python-stdnum/blob/master/stdnum/lt/pvm.py
// Reference: https://github.com/ltns35/go-vat/blob/main/countries/lithuania.go
func validateTaxCodeChecksum(code string) bool {
	if len(code) != 9 && len(code) != 12 {
		return false
	}
	for _, r := range code {
		if r < '0' || r > '9' {
			return false
		}
	}
	check := taxCodeRemainder(code, taxCodeWeights)
	if check == 10 {
		check = taxCodeRemainder(code, taxCodeAltWeights)
		if check == 10 {
			check = 0
		}
	}
	return int(code[len(code)-1]-'0') == check
}

// taxCodeRemainder returns the weighted sum of every digit of the code except
// the last, modulo 11.
func taxCodeRemainder(code string, weights []int) int {
	sum := 0
	for i, w := range weights[:len(code)-1] {
		sum += int(code[i]-'0') * w
	}
	return sum % 11
}
