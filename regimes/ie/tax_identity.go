package ie

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Source: https://github.com/ltns35/go-vat

var (
	taxCodeRegexp = `^(\d{7}[A-Z]{1,2}|\d{1}[A-Z]{1}\d{5}[A-Z]{1})$`
)

// validateTaxIdentity checks to ensure the NIT code looks okay.
func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateTaxCode)),
	)
}

func validateTaxCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()

	if !regexp.MustCompile(taxCodeRegexp).MatchString(val) {
		return errors.New("invalid format")
	}

	return nil
}
