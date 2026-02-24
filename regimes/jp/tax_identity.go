package jp

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var (
	taxCodePattern = regexp.MustCompile(`^\d{13}$`)

	errInvalidFormat   = errors.New("invalid format")
	errInvalidChecksum = errors.New("checksum mismatch")
)

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
		return errInvalidFormat
	}

	// First digit must be non-zero (it's the check digit)
	if val[0] == '0' {
		return errInvalidFormat
	}

	return validateCorporateNumberChecksum(val)
}

// validateCorporateNumberChecksum validates the check digit of a 13-digit Japanese Corporate Number
// https://www.houjin-bangou.nta.go.jp/documents/checkdigit.pdf
// https://github.com/kufu/tsubaki/blob/master/lib/tsubaki/corporate_number.rb
//
// The check digit (first digit) is calculated as:
//
//	9 - ((Σ Pn × Qn) mod 9)
//
// where Pn is the n-th digit of the 12-digit base number (counting from
// the lowest/rightmost digit), and Qn = 1 if n is odd, 2 if n is even.
func validateCorporateNumberChecksum(val string) error {
	sum := 0
	base := val[1:] // 12-digit base number (digits 2-13)

	for i := 0; i < 12; i++ {
		digit, err := strconv.Atoi(string(base[11-i]))
		if err != nil {
			return errInvalidFormat
		}
		// n is 1-indexed: position 1 is the rightmost digit
		n := i + 1
		q := 1
		if n%2 == 0 {
			q = 2
		}
		sum += digit * q
	}

	expected := 9 - (sum % 9)

	checkDigit, err := strconv.Atoi(string(val[0]))
	if err != nil {
		return errInvalidFormat
	}

	if checkDigit != expected {
		return errInvalidChecksum
	}

	return nil
}
