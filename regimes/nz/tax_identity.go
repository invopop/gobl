package nz

import (
	"regexp"

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
		validation.Field(&tID.Code,
			validation.Match(irdRegex).Error("must be 8 or 9 digits"),
		),
	)
}
