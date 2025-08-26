package fr

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var (
	taxCodeVATRegexp = regexp.MustCompile(`^[A-Z0-9]{2}\d{9}$`)
	taxCodeSIRENRegexp = regexp.MustCompile(`^\d{9}$`)
)

// normalizeTaxIdentity normalizes the SIREN code, if there are any errors,
// these will be picked up by validation.
func normalizeTaxIdentity(tID *tax.Identity) {
	if tID.Code == "" {
		return
	}
	tax.NormalizeIdentity(tID)

	str := tID.Code.String()
	// Check if we have a valid SIREN so we can try and
	// normalize with the check digit.
	if err := validateSIRENTaxCode(tID.Code); err != nil {
		return
	}
	chk := calculateVATCheckDigit(str)
	tID.Code = cbc.Code(fmt.Sprintf("%s%s", chk, str))
}

// validateTaxIdentity checks to ensure the SIRET code looks okay.
func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateVATTaxCode)),
	)
}

func validateVATTaxCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	str := code.String()

	if !taxCodeVATRegexp.MatchString(str) {
		return errors.New("invalid format")
	}

	// Extract the last nine digits as an integer.
	siren := str[2:] // extract last nine digits
	chk := calculateVATCheckDigit(siren)
	expectStr := str[:2] // compare with first two digits
	if chk != expectStr {
		return errors.New("checksum mismatch")
	}

	return nil
}

func calculateVATCheckDigit(str string) string {
	// Assume we have a SIREN
	total, _ := strconv.Atoi(str)
	total = (total*100 + 12) % 97

	return fmt.Sprintf("%02d", total)
}

func validateSIRENTaxCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	str := code.String()

	if !taxCodeSIRENRegexp.MatchString(str) {
		return errors.New("invalid format")
	}

	base := str[:8]
	chk := str[8:]
	v := computeLuhnCheckDigit(base)
	if chk != v {
		return errors.New("checksum mismatch")
	}

	return nil
}

// TODO: refactor this into a shareable method.
func computeLuhnCheckDigit(number string) string {
	sum := 0
	pos := 0

	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')

		if pos%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		pos++
	}

	return strconv.FormatInt(int64((10-(sum%10))%10), 10)
}
