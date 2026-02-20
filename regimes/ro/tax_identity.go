package ro

/*
 * Sources of data:
 *
 *  - Romanian Ministry of Finance / ANAF
 *  - https://www.anaf.ro
 *
 */

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var (
	taxCodeRegexp = regexp.MustCompile(`^\d{2,10}$`)
)

func normalizeTaxIdentity(tID *tax.Identity) {
	tax.NormalizeIdentity(tID)
}

func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateTaxCode)),
	)
}

func validateTaxCode(value any) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()

	if !taxCodeRegexp.MatchString(val) {
		return errors.New("invalid format")
	}

	return validateCUIChecksum(val)
}

func validateCUIChecksum(val string) error {
	number, err := strconv.Atoi(val)
	if err != nil {
		return errors.New("invalid format")
	}

	key := 753217532

	controlDigit := number % 10
	number /= 10

	sum := 0
	for number > 0 {
		sum += (number % 10) * (key % 10)
		number /= 10
		key /= 10
	}

	rest := (sum * 10) % 11
	if rest == 10 {
		rest = 0
	}

	if rest != controlDigit {
		return errors.New("checksum mismatch")
	}

	return nil
}
