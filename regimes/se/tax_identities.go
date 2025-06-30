package se

import (
	"errors"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

const (
	// The full length of a Swedish tax ID, including the check digits.
	taxCodeLength = 12
	// The length of the code before the check digits.
	taxCodeLengthWithoutCheckDigits = 10
	// The check digits of a Swedish tax ID.
	taxCodeCheckDigit = "01"
)

// validateTaxIdentity performs validation specific to Swedish tax IDs.
// Assumes the code has already been normalized.
func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code,
			validation.By(validateTaxCode),
		),
	)
}

// validateTaxCode validates the tax code for Swedish tax identities.
// Assumes the code has already been normalized, is made of 12 numeric characters,
// retaining the checksum at the end, plus 2 control digits "01".
func validateTaxCode(value any) error {
	// No need for nil check, as the identity type is validated before this function is called
	code, _ := value.(cbc.Code)
	if code == "" {
		return nil
	}

	// Normalised Swedish tax IDs must have a specific length.
	if len(code) != taxCodeLength {
		return errors.New("invalid length")
	}
	// Swedish tax IDs must finish in "01".
	if code[10:] != taxCodeCheckDigit {
		return errors.New("invalid check digit, expected 01")
	}
	// Swedish tax IDs must be exclusively numeric.
	if _, err := strconv.Atoi(string(code)); err != nil {
		return errors.New("invalid characters, expected numeric")
	}
	// The code prior to the check digit must be Luhn-valid.
	if !(code[:10]).IsValidLuhnChecksum() {
		return errors.New("invalid identification number checksum")
	}
	return nil
}
