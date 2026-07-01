package is

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	rulesis "github.com/invopop/gobl/rules/is" // aliased: rules/is collides with this package's name
	"github.com/invopop/gobl/tax"
)

const (
	// kennitalaLength is the number of digits in a normalized kennitala (DDMMYY-RRCV).
	kennitalaLength = 10
	// kennitalaWeightedDigits is the count of leading digits (positions 1–8) that
	// feed into the MOD-11 check-digit calculation.
	kennitalaWeightedDigits = 8
	// kennitalaCheckDigitIndex is the 0-based index of the check digit (position 9).
	kennitalaCheckDigitIndex = 8
	// kennitalaModulus is the modulus of the MOD-11 check-digit algorithm.
	kennitalaModulus = 11
	// kennitalaNoCheckDigitRemainder is the (S mod 11) remainder for which the check
	// digit would be 10 — not representable as a single digit, so Registers Iceland
	// never issues such numbers and they must be rejected.
	kennitalaNoCheckDigitRemainder = 1

	// A company kennitala encodes its day-of-month offset by 40, so the day field is
	// 41–71 and the first digit is 4–7. A natural person's day is 01–31 (first digit
	// 0–3); temporary "kerfiskennitala" numbers use a first digit of 8 or 9.
	companyFirstDigitMin   = 4
	companyFirstDigitMax   = 7
	personFirstDigitMin    = 0
	personFirstDigitMax    = 3
	temporaryFirstDigitMin = 8
	temporaryFirstDigitMax = 9
)

// kennitalaWeights are the multipliers applied to digits 1–8 when computing the
// MOD-11 check digit. It is a var, not a const, because Go constants cannot hold
// composite types such as arrays.
var kennitalaWeights = [kennitalaWeightedDigits]int{3, 2, 7, 6, 5, 4, 3, 2}

// ValidKennitala reports whether code is a structurally valid kennitala: exactly
// ten digits whose ninth digit matches the MOD-11 check digit computed from the
// first eight. It does not distinguish persons from companies — use
// Company or Person for that. The tenth digit (century
// indicator) is intentionally excluded from the checksum.
func ValidKennitala(code cbc.Code) bool {
	s := code.String()
	if len(s) != kennitalaLength || !isAllDigits(s) {
		return false
	}
	sum := 0
	for i := 0; i < kennitalaWeightedDigits; i++ {
		sum += int(s[i]-'0') * kennitalaWeights[i]
	}
	remainder := sum % kennitalaModulus
	if remainder == kennitalaNoCheckDigitRemainder {
		return false
	}
	check := (kennitalaModulus - remainder) % kennitalaModulus
	return int(s[kennitalaCheckDigitIndex]-'0') == check
}

// Company reports whether code's first digit (4–7) identifies a company
// rather than a person or a temporary number. It does not validate the checksum.
func Company(code cbc.Code) bool {
	d, ok := firstDigit(code)
	return ok && d >= companyFirstDigitMin && d <= companyFirstDigitMax
}

// Person reports whether code's first digit (0–3) identifies a natural
// person. It does not validate the checksum.
func Person(code cbc.Code) bool {
	d, ok := firstDigit(code)
	return ok && d >= personFirstDigitMin && d <= personFirstDigitMax
}

func isTemporaryKennitala(code cbc.Code) bool {
	d, ok := firstDigit(code)
	return ok && d >= temporaryFirstDigitMin && d <= temporaryFirstDigitMax
}

func firstDigit(code cbc.Code) (int, bool) {
	s := code.String()
	if len(s) == 0 || s[0] < '0' || s[0] > '9' {
		return 0, false
	}
	return int(s[0] - '0'), true
}

func isAllDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// isCanonicalKennitala reports whether code is a 10-digit numeric string — the
// precondition for the checksum and classification checks. It lets each downstream
// assertion skip cleanly, so a malformed code is reported only by the format or
// length assertion.
func isCanonicalKennitala(code cbc.Code) bool {
	s := code.String()
	return isAllDigits(s) && len(s) == kennitalaLength
}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid kennitala format",
					rulesis.Func("numeric", kennitalaFormatValid),
				),
				rules.AssertIfPresent("02", "invalid kennitala length",
					rulesis.Func("length", kennitalaLengthValid),
				),
				rules.AssertIfPresent("03", "invalid checksum for kennitala",
					rulesis.Func("checksum", kennitalaChecksumValid),
				),
				rules.AssertIfPresent("04", "person kennitala is not valid for VAT",
					rulesis.Func("not person", kennitalaNotPerson),
				),
				rules.AssertIfPresent("05", "temporary kennitala is not valid for VAT",
					rulesis.Func("not temporary", kennitalaNotTemporary),
				),
			),
		),
	)
}

func kennitalaFormatValid(value any) bool {
	code, ok := value.(cbc.Code)
	return ok && isAllDigits(code.String())
}

func kennitalaLengthValid(value any) bool {
	code, _ := value.(cbc.Code)
	s := code.String()
	if !isAllDigits(s) {
		return true // not numeric: the format assertion reports this
	}
	return len(s) == kennitalaLength
}

func kennitalaChecksumValid(value any) bool {
	code, _ := value.(cbc.Code)
	if !isCanonicalKennitala(code) {
		return true // format/length assertions report these
	}
	return ValidKennitala(code)
}

func kennitalaNotPerson(value any) bool {
	code, _ := value.(cbc.Code)
	if !isCanonicalKennitala(code) {
		return true
	}
	return !Person(code)
}

func kennitalaNotTemporary(value any) bool {
	code, _ := value.(cbc.Code)
	if !isCanonicalKennitala(code) {
		return true
	}
	return !isTemporaryKennitala(code)
}
