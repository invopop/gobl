package ro

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Romanian tax identification information:
// - CUI (Cod Unic de Înregistrare) / CIF (Cod de Identificare Fiscală)
// - Format: [RO] + 2-10 digits
//
// References:
// - ANAF Checksum Algorithm: Common fiscal algorithm for "Cifra de control"
// - Validates against Law 227/2015 requirements
var (
	taxCodeRegexp = regexp.MustCompile(`^(RO)?[0-9]{2,10}$`)
	// Weights for CUI checksum (reversed logic: 753217532 aligned right-to-left)
	// as we will range 0:9 to simplify the checksum logic.
	cifWeights = []int{7, 5, 3, 2, 1, 7, 5, 3, 2}
)

func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code,
			validation.Required,
			validation.By(validateTaxCode),
		),
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

	// Normalize locally for validation (remove RO prefix if present)
	// No need to check numeric content as regex ensures this
	val = strings.TrimPrefix(val, "RO")

	return validateCIFChecksum(val)
}

// validateCIFChecksum validates the Romanian CIF/CUI checksum.
// Algorithm:
// 1. Pad code with leading zeros to 9 digits (excluding control digit).
// 2. Multiply digits by weights [7,5,3,2,1,7,5,3,2].
// 3. Sum results.
// 4. Control = (Sum * 10) % 11.
// 5. If Control == 10, then Control = 0.
func validateCIFChecksum(code string) error {
	length := len(code)
	if length < 2 || length > 10 {
		return errors.New("invalid length")
	}

	// The last digit is the control digit
	controlDigit := int(code[length-1] - '0')

	// The data part is everything before the last digit
	dataStr := code[:length-1]

	// Pad with leading zeros to ensure exactly 9 digits
	// e.g. "123" becomes "000000123"
	paddedData := fmt.Sprintf("%09s", dataStr)

	var sum int
	for i := range 9 {
		digit := int(paddedData[i] - '0')
		sum += digit * cifWeights[i]
	}

	// This is equivalent to 11 - sum % 11, but adjusted for the modulo 10 rule
	// to simplify the logic. This replaces the use of
	// remainder = sum % 11
	// checkDigit = 11 - remainder
	//
	// Also, instead of having to have an if branch to make the calculated control 0 when
	// modulo is 10, we can use modulo 10 directly.
	calculatedControl := (sum * 10) % 11 % 10

	if controlDigit != calculatedControl {
		return errors.New("invalid checksum")
	}

	return nil
}
