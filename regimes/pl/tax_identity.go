package pl

import (
	"errors"
	"regexp"
	"strconv"
	"unicode"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

/*
 * Sources of data:
 *
 *  - https://pl.wikipedia.org/wiki/Numer_identyfikacji_podatkowej
 *
 */

const (
	taxIdentityPattern = `^[1-9]((\d[1-9])|([1-9]\d))\d{7}$`
)

var (
	taxIdentityRegexp = regexp.MustCompile(taxIdentityPattern)
)

func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code,
			validation.Required,
			validation.By(validateTaxCode),
		),
	)
}

func validateTaxCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok {
		return nil
	}
	if code == "" {
		return nil
	}

	if taxIdentityRegexp.MatchString(code.String()) {
		if validateNIPChecksum(code) {
			return nil
		}
		return errors.New("checksum mismatch")
	}

	return errors.New("invalid format")
}

func validateNIPChecksum(code cbc.Code) bool {
	nipStr := code.String()
	if len(nipStr) != 10 {
		return false
	}

	for _, char := range nipStr {
		if !unicode.IsDigit(char) {
			return false
		}
	}

	digits := make([]int, 10)
	for i, char := range nipStr {
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return false
		}
		digits[i] = digit
	}

	weights := [9]int{6, 5, 7, 2, 3, 4, 5, 6, 7}
	checkSum := 0
	for i, digit := range digits[:9] {
		checkSum += digit * weights[i]
	}
	checkSum %= 11

	return checkSum == digits[9]
}
