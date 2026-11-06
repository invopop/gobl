package cz

import (
	"regexp"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// A DIČ is the code "CZ" followed by a root: the IČO of a legal person, the
// birth number of a natural person, or an identifier assigned by the tax
// office (Act 280/2009 Sb., § 130). The prefix is removed during
// normalization, so only the root reaches validation.
var taxCodeRegexp = regexp.MustCompile(`^\d{8,10}$`)

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "tax identity code for CZ must be 8, 9 or 10 digits",
					is.MatchesRegexp(taxCodeRegexp),
				),
				rules.AssertIfPresent("02", "tax identity code for CZ failed the check digit",
					is.StringFunc("check digit", validateCheckDigit),
				),
			),
		),
	)
}

// validateCheckDigit applies the check appropriate to the length of a
// well-formed code. A malformed code is left to the format rule.
func validateCheckDigit(val string) bool {
	switch {
	case !taxCodeRegexp.MatchString(val):
		return true
	case len(val) == 8:
		return validateCompanyNumber(val)
	case len(val) == 10:
		return validateBirthNumber(val)
	}
	// Nine digits carry no check digit: birth numbers issued before 1954
	// (Act 133/2000 Sb., § 13(4)) and the tax office's own identifiers.
	return true
}

// The IČO (identifikační číslo osoby) check digit.
//
// Format: 7 digits + check digit. Validation: MOD 11 with weights 8 down to 2,
// where the check digit is 11 minus the remainder, modulo 10.
//
// The Czech Statistical Office, which assigns the IČO, does not publish the
// algorithm. It is consistently implemented by public validators and holds
// for every known IČO, including the EET sample identity CZ00000019.
//
// Reference: https://github.com/arthurdejong/python-stdnum/blob/master/stdnum/cz/dic.py
func validateCompanyNumber(val string) bool {
	sum := 0
	for i, c := range val[:7] {
		sum += int(c-'0') * (8 - i)
	}
	return int(val[7]-'0') == (11-sum%11)%10
}

// The birth number (rodné číslo) check digit.
//
// Format: 9 digits + check digit, the first two being the year of birth.
// Validation: the check digit is the remainder of the first nine digits modulo
// 11, which makes the whole number divisible by 11 (Act 133/2000 Sb., § 13(3)).
// Numbers issued between 1954 and 1985 wrote a remainder of 10 as 0; later
// numbers with that remainder were never issued.
//
// Reference: https://www.zakonyprolidi.cz/cs/2000-133#p13
// Reference: https://cs.wikipedia.org/wiki/Rodn%C3%A9_%C4%8D%C3%ADslo
func validateBirthNumber(val string) bool {
	rem := 0
	for _, c := range val[:9] {
		rem = (rem*10 + int(c-'0')) % 11
	}
	check := int(val[9] - '0')
	if rem == 10 {
		year := val[:2]
		return check == 0 && year >= "54" && year <= "85"
	}
	return check == rem
}
