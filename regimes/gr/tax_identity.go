package gr

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Reference: https://lytrax.io/blog/projects/greek-tin-validator-generator

var (
	taxCodeRegexp = regexp.MustCompile(`^\d{9}$`)
)

// normalizeTaxIdentity requires additional steps for Greece as the language code
// is used in the tax code.
func normalizeTaxIdentity(tID *tax.Identity) error {
	if tID == nil {
		return nil
	}
	// also allow for usage of "GR" which may be used in the tax code
	// by accident.
	if err := tax.NormalizeIdentity(tID, l10n.GR); err != nil {
		return err
	}
	tID.Country = "EL" // always override for greece
	return nil
}

// validateTaxIdentity checks to ensure the tax code looks okay.
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

	if !taxCodeRegexp.MatchString(val) {
		return errors.New("invalid format")
	}

	if !hasValidChecksum(val) {
		return errors.New("checksum mismatch")
	}

	return nil
}

func hasValidChecksum(val string) bool {
	digits := make([]int, 9)
	for i, char := range val {
		num, err := strconv.Atoi(string(char))
		if err != nil {
			return false
		}
		digits[i] = num
	}

	var sum int
	for i := 0; i < 8; i++ {
		sum += digits[i] * (1 << uint(8-i))
	}
	checkDigit := sum % 11 % 10

	return checkDigit == digits[8]
}
