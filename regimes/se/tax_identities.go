package se

import (
	"errors"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/internal"
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

var (
	// ErrInvalidTaxIDLength is returned when the tax ID is not the correct length.
	ErrInvalidTaxIDLength = errors.New("invalid length")
	// ErrInvalidTaxIDCountryPrefix is returned when the tax ID does not start with "SE".
	ErrInvalidTaxIDCountryPrefix = errors.New("invalid country prefix")
	// ErrInvalidTaxIDCheckDigit is returned when the tax ID does not end with "01".
	ErrInvalidTaxIDCheckDigit = errors.New("invalid check digit")
	// ErrInvalidTaxIDCharacters is returned when the tax ID contains non-numeric characters.
	ErrInvalidTaxIDCharacters = errors.New("invalid characters, expected numeric")
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
	code, ok := value.(cbc.Code)
	if !ok {
		return nil
	}
	if code == "" {
		return nil
	}

	// Normalised Swedish tax IDs must have a specific length.
	if len(code) != taxCodeLength {
		return ErrInvalidTaxIDLength
	}
	// Swedish tax IDs must finish in "01".
	if code[10:] != taxCodeCheckDigit {
		return ErrInvalidTaxIDCheckDigit
	}
	// Swedish tax IDs must be exclusively numeric.
	if _, err := strconv.Atoi(string(code)); err != nil {
		return ErrInvalidTaxIDCharacters
	}
	// The code prior to the check digit must be Luhn valid.
	if !internal.ValidateLuhn(string(code[:10])) {
		return ErrInvalidChecksum
	}
	return nil
}
