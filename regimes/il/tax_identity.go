// Package il provides the tax identity validation specific to Israel
package il

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Israeli business tax ID (Mispar Osek Murshe / מספר עוסק מורשה) is a 9-digit number.
//
// Format: 9 numeric digits, no separators.
//
// The checksum algorithm for this number has not been found in any official
// ITA source consulted. Therefore, only format validation (9 numeric digits)
// is applied here. This eliminates the vast majority of invalid inputs safely,
// without risking rejection of legitimate tax IDs due to an unverified algorithm.
//
// Sources:
// - ITA API Specification v1.0 (July 2023), Table 2.1: Vat_Number field defined as N9
//   https://www.gov.il/BlobFolder/generalpage/israel-invoice-160723/he/IncomeTax_software-houses-en-040723.pdf

var (
	osekRegex = regexp.MustCompile(`^\d{9}$`)
)

// validateTaxIdentity checks to ensure the Israeli Osek Murshe format is correct.
func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateOsekCode)),
	)
}

// validateOsekCode checks that the tax ID is a valid 9-digit format.
func validateOsekCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()

	if !osekRegex.MatchString(val) {
		return errors.New("must be a 9-digit number")
	}
	return nil
}
