package si

import (
	"errors"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// taxCodeLen is the expected length of a Slovenian tax number (davčna
// številka): seven digits followed by a modulo-11 check digit. The first
// digit is never zero.
const taxCodeLen = 8

// taxCodeMultipliers are the weights applied to the first seven digits of
// the tax number to calculate the check digit.
//
// The official register documents the format and the modulo-11 control,
// but not the formula itself.
// Format source: https://www.fu.gov.si/en/taxes_and_other_duties/work_with_us/entry_into_the_tax_register_and_tax_number
// Algorithm reference: https://github.com/arthurdejong/python-stdnum (stdnum/si/ddv.py)
var taxCodeMultipliers = []int{8, 7, 6, 5, 4, 3, 2}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Slovenian VAT identity code",
					is.Func("valid", isValidTaxIdentityCode),
				),
			),
		),
	)
}

func isValidTaxIdentityCode(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return false
	}
	return validateTaxCode(code) == nil
}

func validateTaxCode(code cbc.Code) error {
	if code == "" {
		return nil
	}
	if len(code) != taxCodeLen {
		return errors.New("invalid length")
	}
	if code[0] == '0' {
		return errors.New("invalid format")
	}
	return validateMod11(code, taxCodeMultipliers)
}

// validateMod11 checks that code is made up of digits and that its last
// digit is the modulo-11 check digit of the preceding ones, weighted with
// the given multipliers. It expects exactly len(weights)+1 characters.
//
// Both Slovenian identifiers (the tax number and the registration number)
// share this scheme with different weights: a remainder of 0 is never
// assigned (no valid number produces it), a remainder of 1 means a check
// digit of 0, and any other remainder means a check digit of 11 minus the
// remainder.
func validateMod11(code cbc.Code, weights []int) error {
	sum := 0
	for i, w := range weights {
		c := code[i]
		if c < '0' || c > '9' {
			return errors.New("invalid characters")
		}
		sum += int(c-'0') * w
	}
	check := code[len(weights)]
	if check < '0' || check > '9' {
		return errors.New("invalid characters")
	}
	switch r := sum % 11; r {
	case 0:
		return errors.New("invalid checksum")
	case 1:
		if check != '0' {
			return errors.New("checksum mismatch")
		}
	default:
		if int(check-'0') != 11-r {
			return errors.New("checksum mismatch")
		}
	}
	return nil
}
