package si

import (
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// taxCodePattern matches a Slovenian tax number (davčna številka): eight
// digits, the first of which is never zero, followed by a modulo-11 check
// digit.
//
// The official register documents the format and the modulo-11 control,
// but not the formula itself.
// Format source: https://www.fu.gov.si/en/taxes_and_other_duties/work_with_us/entry_into_the_tax_register_and_tax_number
// Algorithm reference: https://github.com/arthurdejong/python-stdnum (stdnum/si/ddv.py)
const taxCodePattern = `^[1-9]\d{7}$`

// taxCodeMultipliers are the weights applied to the first seven digits of
// the tax number to calculate the check digit.
var taxCodeMultipliers = []int{8, 7, 6, 5, 4, 3, 2}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Slovenian VAT identity format",
					is.Matches(taxCodePattern),
				),
				rules.AssertIfPresent("02", "invalid Slovenian VAT identity check digit",
					is.StringFunc("checksum", taxCodeChecksumValid),
				),
			),
		),
	)
}

// taxCodeChecksumValid reports whether the tax number's final digit is a
// valid modulo-11 check digit. It runs alongside the format rule, so it
// guards against malformed input rather than assuming it.
func taxCodeChecksumValid(code string) bool {
	return validMod11(code, taxCodeMultipliers)
}

// validMod11 reports whether the last digit of code is the modulo-11 check
// digit of the preceding ones, weighted with the given multipliers. It
// expects exactly len(weights)+1 numeric characters and returns false for
// anything else, so it is safe to call before the format rule has run.
//
// Both Slovenian identifiers (the tax number and the registration number)
// share this scheme with different weights: a remainder of 0 is never
// assigned (no valid number produces it), a remainder of 1 means a check
// digit of 0, and any other remainder means a check digit of 11 minus the
// remainder.
func validMod11(code string, weights []int) bool {
	if len(code) != len(weights)+1 {
		return false
	}
	sum := 0
	for i, w := range weights {
		if code[i] < '0' || code[i] > '9' {
			return false
		}
		sum += int(code[i]-'0') * w
	}
	check := code[len(weights)]
	if check < '0' || check > '9' {
		return false
	}
	switch r := sum % 11; r {
	case 0:
		return false
	case 1:
		return check == '0'
	default:
		return int(check-'0') == 11-r
	}
}
