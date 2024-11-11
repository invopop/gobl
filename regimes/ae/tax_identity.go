// Package ae provides the tax identity validation specific to the United Arab Emirates.
package ae

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var (
	// TRN in UAE is a 15-digit number
	trnRegex = regexp.MustCompile(`^\d{15}$`)
)

// validateTaxIdentity checks to ensure the UAE TRN format is correct.
func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateTRNCode)),
	)
}

// validateTRNCode checks that the TRN is a valid 15-digit format.
func validateTRNCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()

	// Check if TRN matches the 15-digit pattern
	if !trnRegex.MatchString(val) {
		return errors.New("invalid format: TRN must be a 15-digit number")
	}

	return nil
}
