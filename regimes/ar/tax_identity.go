package ar

import (
	"errors"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// CUIT (Clave Única de Identificación Tributaria) and CUIL (Clave Única de Identificación Laboral)
// are 11-digit tax identification codes used in Argentina.
//
// Format: XX-XXXXXXXX-X (typically displayed with hyphens, but stored without)
//
// The validation algorithm uses modulo 11 with specific multipliers.
// Multipliers sequence: [5, 4, 3, 2, 7, 6, 5, 4, 3, 2]
//
// Validation steps:
//  1. Multiply each of the first 10 digits by the corresponding multiplier
//  2. Sum all the results
//  3. Calculate: checkDigit = 11 - (sum % 11)
//  4. Special cases:
//     - If checkDigit = 11, then checkDigit = 0
//     - If checkDigit = 10, the prefix must be adjusted:
//     * For individuals: change 20 to 23 (men) or 27 to 28 (women)
//     * For companies: change 30 to 33
//     And then checkDigit = 9
//
// Common prefixes:
//   - 20: Male individual (CUIL)
//   - 23: Male individual with check digit conflict resolution
//   - 27: Female individual (CUIL)
//   - 28: Female individual with check digit conflict resolution
//   - 30: Company/Legal entity (CUIT)
//   - 33: Company with check digit conflict resolution
//   - 34: Foreign entity
//
// References:
// - AFIP Official Documentation: https://www.afip.gob.ar/
// - Algorithm explanation: https://whiz.tools/en/legal-business/argentinian-cuit-cuil-generator-validator
// - Validation details: https://lib.rs/crates/ar_cuil_cuit_validator
// - Tax ID guide: https://lookuptax.com/docs/tax-identification-number/Argentina-tax-id-guide
// - Implementation reference: https://www.lawebdelprogramador.com/codigo/Visual-Basic/160-Verificar-el-CUIT-CUIL-Argentina.html

func normalizeTaxIdentity(tID *tax.Identity) {
	// Use standard normalization to remove spaces and convert to uppercase
	tax.NormalizeIdentity(tID)
	// Remove hyphens commonly used in CUIT/CUIL formatting
	code := tID.Code.String()
	// Remove all hyphens
	normalized := ""
	for _, char := range code {
		if char != '-' {
			normalized += string(char)
		}
	}
	tID.Code = cbc.Code(normalized)
}

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

	// CUIT/CUIL must be exactly 11 digits
	if len(val) != 11 {
		return errors.New("must have 11 digits")
	}

	// Verify all characters are digits
	for _, char := range val {
		if char < '0' || char > '9' {
			return errors.New("must contain only digits")
		}
	}

	// Validate using modulo 11 algorithm
	return validateCUITCUIL(val)
}

func validateCUITCUIL(code string) error {
	// Multipliers for CUIT/CUIL validation
	multipliers := []int{5, 4, 3, 2, 7, 6, 5, 4, 3, 2}

	// Calculate the sum
	sum := 0
	for i := 0; i < 10; i++ {
		digit, err := strconv.Atoi(string(code[i]))
		if err != nil {
			return errors.New("must contain only digits")
		}
		sum += digit * multipliers[i]
	}

	// Calculate expected check digit
	remainder := sum % 11
	expectedCheckDigit := 11 - remainder

	// Special cases for check digit
	if expectedCheckDigit == 11 {
		expectedCheckDigit = 0
	} else if expectedCheckDigit == 10 {
		// When check digit would be 10, the prefix must be adjusted
		// and the check digit becomes 9
		// We need to verify the prefix has been properly adjusted
		prefix := code[0:2]
		// Valid prefixes that handle check digit 10 case
		validPrefixes := []string{"23", "28", "33"}
		hasValidPrefix := false
		for _, vp := range validPrefixes {
			if prefix == vp {
				hasValidPrefix = true
				break
			}
		}
		if !hasValidPrefix {
			return errors.New("invalid prefix for check digit conflict")
		}
		expectedCheckDigit = 9
	}

	// Get actual check digit (last digit)
	actualCheckDigit, err := strconv.Atoi(string(code[10]))
	if err != nil {
		return errors.New("must contain only digits")
	}

	// Verify check digit matches
	if actualCheckDigit != expectedCheckDigit {
		return errors.New("verification digit mismatch")
	}

	return nil
}
