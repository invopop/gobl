package my

import (
	"errors"
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var (
	// Malaysian Business Registration Numbers are typically 12-digit numbers
	taxCodeMYNumeric = regexp.MustCompile(`^\d{12}$`)

	// Alphanumeric SST/Service Tax IDs — e.g., SST1234567890 or W10-12345678-123
	taxCodeSST    = regexp.MustCompile(`^SST\d{10,12}$`)
	taxCodeWStyle = regexp.MustCompile(`^[A-Z0-9]{1,4}-\d{8}-\d{3}$`)
)

// normalizeTaxIdentity performs basic uppercasing and trims whitespace for consistency.
func normalizeTaxIdentity(tID *tax.Identity) {
	if tID == nil {
		return
	}
	code := strings.TrimSpace(tID.Code.String()) // trim spaces
	code = strings.ReplaceAll(code, " ", "")     // remove internal spaces too
	tID.Code = cbc.Code(strings.ToUpper(code))   // uppercase
}

// validateTaxIdentity validates Malaysian tax IDs.
func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateMYTaxCode)),
	)
}

func validateMYTaxCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	str := strings.ToUpper(strings.TrimSpace(code.String()))

	if taxCodeMYNumeric.MatchString(str) || taxCodeSST.MatchString(str) || taxCodeWStyle.MatchString(str) {
		return nil
	}

	return errors.New("invalid MY tax ID format")
}
