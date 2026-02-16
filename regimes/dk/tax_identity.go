package dk

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Reference: https://github.com/ltns35/go-vat

var (
	taxCodeRegexps = []*regexp.Regexp{
		regexp.MustCompile(`^\d{8}$`),
	}
)

// validateTaxIdentity checks to ensure the Danish CVR code is valid.
func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateTaxCode)),
	)
}

func validateTaxCode(value any) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()

	match := false
	for _, re := range taxCodeRegexps {
		if re.MatchString(val) {
			match = true
			break
		}
	}
	if !match {
		return errors.New("invalid format")
	}

	return validateTaxCodeChecksum(val)
}

func validateTaxCodeChecksum(val string) error {
	// Danish CVR numbers use modulo-11 checksum with multipliers [2, 7, 6, 5, 4, 3, 2, 1]
	multipliers := []int{2, 7, 6, 5, 4, 3, 2, 1}
	total := 0

	for i := range 8 {
		digit, err := strconv.Atoi(string(val[i]))
		if err != nil {
			return errors.New("invalid digit")
		}
		total += digit * multipliers[i]
	}

	if total%11 != 0 {
		return errors.New("checksum mismatch")
	}

	return nil
}
