package sa

import (
	"regexp"

	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Reference: https://zatca.gov.sa
// Reference: https://lookuptax.com/docs/tax-identification-number/saudi-arabia-tax-id-guide
//
// The Saudi Arabia VAT Tax Identification Number (TIN) is a 15-digit number
// with the following structure:
//
//   - Position 1:     GCC member state code (always 3 for Saudi Arabia)
//   - Positions 2-9:  Taxpayer serial number
//   - Position 10:    Check digit (algorithm undisclosed)
//   - Positions 11-13: Subsidiary/branch code
//   - Positions 14-15: Tax type indicator

var (
	// ZATCA requires the first digit to be 3 (GCC code for Saudi Arabia)
	// and the last digit to be 3 (VAT tax type).
	tinRegex = regexp.MustCompile(`^3\d{13}3$`)
)

// validateTaxIdentity checks to ensure the SA TIN format is correct.
func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code,
			validation.Match(tinRegex).Error("must be a 15-digit number starting and ending with 3"),
		),
	)
}
