package pay

import (
	"regexp"

	"github.com/invopop/gobl/rules/is"
)

var (
	// ibanPattern reflects the ISO 13616 structure: country code, check
	// digits, and 11 to 30 characters of national account data. Per-country
	// lengths are not enforced.
	ibanPattern = regexp.MustCompile(`^[A-Z]{2}\d{2}[A-Z0-9]{11,30}$`)

	// bicPattern reflects the ISO 9362 structure: institution and country
	// codes, a location code, and an optional branch code.
	bicPattern = regexp.MustCompile(`^[A-Z]{6}[A-Z0-9]{2}([A-Z0-9]{3})?$`)
)

var (
	// IsIBAN confirms the value is an International Bank Account Number.
	IsIBAN = is.AllOf(
		is.MatchesRegexp(ibanPattern),
		is.StringFunc("has matching check digits", ibanCheckDigitsMatch),
	)

	// IsBIC confirms the value is a Business Identifier Code.
	IsBIC = is.MatchesRegexp(bicPattern)
)

// ibanCheckDigitsMatch runs the ISO 7064 mod 97-10 calculation over the
// account number moved in front of the country code and check digits, with
// letters converted to two digits each (A is 10, Z is 35).
func ibanCheckDigitsMatch(iban string) bool {
	if len(iban) < 5 {
		return false // too short to carry check digits and an account number
	}
	n := 0
	for _, r := range iban[4:] + iban[:4] {
		if r >= 'A' && r <= 'Z' {
			n = (n*100 + int(r-'A') + 10) % 97
		} else {
			n = (n*10 + int(r-'0')) % 97
		}
	}
	return n == 1
}
