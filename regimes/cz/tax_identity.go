package cz

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Reference: https://github.com/ltns35/go-vat

var (
	taxCodePattern = regexp.MustCompile(`^\d{8,10}$`)

	errTaxCodeInvalidFormat    = errors.New("invalid format")
	errTaxCodeChecksumMismatch = errors.New("checksum mismatch")
)

// validateTaxIdentity checks to ensure the Czech DIČ code is valid.
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

	if !taxCodePattern.MatchString(val) {
		return errTaxCodeInvalidFormat
	}

	switch len(val) {
	case 8:
		return validateLegalEntityCode(val)
	case 9:
		// 9-digit codes include special IDs (starting with 6) and
		// older individual birth number formats; no checksum required.
		return nil
	case 10:
		return validateIndividualCode(val)
	}

	return nil
}

// validateLegalEntityCode validates an 8-digit legal entity DIČ using
// modulo-11 checksum with weights [8,7,6,5,4,3,2].
func validateLegalEntityCode(val string) error {
	weights := []int{8, 7, 6, 5, 4, 3, 2}
	total := 0

	for i := range 7 {
		total += int(val[i]-'0') * weights[i]
	}

	expected := 11 - (total % 11)
	if expected == 10 {
		expected = 0
	} else if expected == 11 {
		expected = 1
	}

	checkDigit := int(val[7] - '0')
	if checkDigit != expected {
		return errTaxCodeChecksumMismatch
	}

	return nil
}

// validateIndividualCode validates a 10-digit individual DIČ (derived from
// Rodné číslo). Must be divisible by 11.
func validateIndividualCode(val string) error {
	// The regex already ensures only digits; ParseUint will not fail.
	n, _ := strconv.ParseUint(val, 10, 64)

	if n%11 != 0 {
		return errTaxCodeChecksumMismatch
	}

	return nil
}
