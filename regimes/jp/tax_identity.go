package jp

import (
	"errors"
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var (
	// Corporate Number: exactly 13 digits.
	taxCodePattern = regexp.MustCompile(`^\d{13}$`)
)

// normalizeTaxIdentity removes whitespace, hyphens, and the "T" prefix
// used in Registration Numbers (適格請求書発行事業者番号, Tekikaku Seikyūsho Hakkō Jigyōsha Bangō).
func normalizeTaxIdentity(tID *tax.Identity) {
	if tID == nil {
		return
	}
	tax.NormalizeIdentity(tID)
	// Remove the "T" prefix used in Qualified Invoice Registration Numbers.
	// The underlying number is still a Corporate Number.
	code := tID.Code.String()
	code = strings.TrimPrefix(code, "T")
	tID.Code = cbc.Code(code)
}

// validateTaxIdentity checks that the tax identity contains a valid
// 13-digit Japanese Corporate Number (法人番号, Hōjin Bangō) with correct checksum.
func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateTaxCode)),
	)
}

func validateTaxCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()

	if !taxCodePattern.MatchString(val) {
		return errors.New("must be a 13-digit number")
	}

	if !validateCorporateNumberChecksum(val) {
		return errors.New("invalid checksum")
	}

	return nil
}

// validateCorporateNumberChecksum verifies the check digit of a 13-digit
// Corporate Number using the official algorithm:
//
//	check_digit = 9 - (SUM(Pn * Qn) for n=1..12) mod 9
//
// where Pn is the n-th digit from the RIGHT of the 12 base digits (digits 2-13),
// and Qn is 1 if n is odd, 2 if n is even.
func validateCorporateNumberChecksum(code string) bool {
	if len(code) != 13 {
		return false
	}

	// The first digit is the check digit; digits 2-13 are the base number.
	checkDigit := int(code[0] - '0')
	base := code[1:] // 12 digits

	sum := 0
	for n := 1; n <= 12; n++ {
		// Pn: n-th digit from the right of the 12-digit base.
		p := int(base[12-n] - '0')
		// Qn: 1 if n is odd, 2 if n is even.
		q := 1
		if n%2 == 0 {
			q = 2
		}
		sum += p * q
	}

	expected := 9 - (sum % 9)
	return checkDigit == expected
}
