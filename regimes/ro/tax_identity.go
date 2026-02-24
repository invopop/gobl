package ro

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var (
	taxCodeWeights = []int{7, 5, 3, 2, 1, 7, 5, 3, 2}
	taxCodeRegexp  = regexp.MustCompile(`^\d{2,10}$`)
)

// validateTaxIdentity checks to ensure the CUI/CIF tax code looks okay.
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

	if !taxCodeRegexp.MatchString(val) {
		return errors.New("invalid format")
	}

	return validateCUIChecksum(val)
}

// validateCUIChecksum validates a Romanian CUI/CIF using the standard
// weighted checksum algorithm:
//  1. Separate the last digit as the expected check digit.
//  2. Right-align the remaining digits against the weight table and
//     multiply each digit by its corresponding weight, summing the results.
//  3. Compute: (sum * 10) % 11. If the result is 10, use 0.
//  4. Compare with the expected check digit.
func validateCUIChecksum(val string) error {
	last := len(val) - 1
	expected := int(val[last] - '0')

	// The weights are right-aligned: a 4-digit body uses the last 4 weights.
	body := val[:last]
	offset := len(taxCodeWeights) - len(body)

	sum := 0
	for i, ch := range body {
		digit := int(ch - '0')
		sum += digit * taxCodeWeights[offset+i]
	}

	check := (sum * 10) % 11
	if check == 10 {
		check = 0
	}

	if check != expected {
		return errors.New("checksum mismatch")
	}

	return nil
}
