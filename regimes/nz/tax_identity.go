package nz

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var (
	// IRD numbers are typically 8 or 9 digits.
	irdRegex = regexp.MustCompile(`^\d{8,9}$`)
)

// validateTaxIdentity checks to ensure the NZ IRD format is correct.
func validateTaxIdentity(tID *tax.Identity) error {
	if tID == nil {
		return nil
	}
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateIRDCode)),
	)
}

// validateIRDCode checks that the IRD number is a valid 8-9 digit format.
func validateIRDCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()

	if !irdRegex.MatchString(val) {
		return errors.New("must be 8 or 9 digits")
	}

	return nil
}
