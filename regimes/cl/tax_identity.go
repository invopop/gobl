package cl

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var (
	// rutPattern matches normalized Chilean RUT format (no separators).
	// Format: 6-8 digits followed by a check digit (0-9 or K/k), for a total of 7-9 characters.
	// Modern RUTs typically have 8-9 digits total; older RUTs may have 7.
	// Leading zeros are not allowed as real Chilean RUTs don't start with zero.
	// Examples: "713254975" (9), "77668208K" (9), "10000009" (8)
	rutPattern = regexp.MustCompile(`^([1-9]\d{5,7})([\dKk])$`)
)

// validateTaxIdentity checks to ensure the RUT code looks okay.
func validateTaxIdentity(tID *tax.Identity) error {
	if tID == nil {
		return nil
	}
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code,
			validation.When(
				tID.Country.In("CL"),
				validation.By(validateRUT),
			),
		),
	)
}

// normalizeTaxIdentity removes formatting characters from Chilean RUT numbers.
//
// Normalization process:
//  1. Removes all dots (.) and hyphens (-)
//  2. Converts lowercase 'k' to uppercase 'K' for check digits
//
// Examples:
//   - "71.325.497-5" → "713254975"
//   - "77.668.208-k" → "77668208K"
//   - "12.345.678-5" → "123456785"
//
// Reference: https://en.wikipedia.org/wiki/National_identification_number#Chile
func normalizeTaxIdentity(tID *tax.Identity) {
	if tID == nil {
		return
	}
	tax.NormalizeIdentity(tID)
}

// validateRUT validates a Chilean RUT (Rol Único Tributario) tax identification number.
//
// The RUT consists of a number (6-8 digits) followed by a check digit calculated
// using the modulo 11 algorithm. The check digit can be 0-9 or K (for check value 10).
// Total length: 7-9 characters.
//
// Validation process:
//  1. Verify the RUT matches the expected format (6-8 digits + check digit = 7-9 total)
//  2. Calculate the expected check digit using modulo 11 algorithm
//  3. Compare the calculated check digit with the provided one
//
// Returns nil if valid, or an error describing the validation failure.
//
// References:
//   - RUT Format: https://es.wikipedia.org/wiki/Rol_%C3%9Anico_Tributario
//   - RUT Length Updates (7-9 digits): https://help.getcirrus.com/es/articles/8515715-modificaciones-para-manejo-de-rut-de-7-a-9-digitos-chile
func validateRUT(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}

	str := strings.ToUpper(string(code))

	matches := rutPattern.FindStringSubmatch(str)
	if matches == nil {
		return errors.New("invalid RUT format")
	}

	number := matches[1]
	checkDigit := matches[2]

	expected, err := calculateRUTCheckDigit(number)
	if err != nil {
		return err
	}

	if expected != checkDigit {
		return errors.New("invalid RUT check digit")
	}

	return nil
}

// calculateRUTCheckDigit calculates the check digit for a Chilean RUT using the modulo 11 algorithm.
//
// The algorithm works as follows:
//  1. Process each digit from right to left
//  2. Multiply each digit by a factor that cycles through 2, 3, 4, 5, 6, 7
//  3. Sum all the products
//  4. Calculate: 11 - (sum mod 11)
//  5. Special cases: result 11 → "0", result 10 → "K"
//
// Example for RUT 71325497:
//
//	7*2 + 9*3 + 4*4 + 5*5 + 2*6 + 3*7 + 1*2 + 7*3 = 14+27+16+25+12+21+2+21 = 138
//	11 - (138 mod 11) = 11 - 6 = 5
//	Result: "713254975"
//
// Parameters:
//
//	rut: The numeric portion of the RUT as a string (without check digit)
//
// Returns:
//
//	The calculated check digit as a string ("0"-"9" or "K"), or an error if conversion fails
//
// Reference: https://es.wikipedia.org/wiki/Rol_%C3%9Anico_Tributario#Algoritmo
func calculateRUTCheckDigit(rut string) (string, error) {
	sum := 0
	multiplier := 2

	// Process digits from right to left
	for i := len(rut) - 1; i >= 0; i-- {
		digit := int(rut[i] - '0')

		if digit < 0 || digit > 9 {
			return "", strconv.ErrSyntax
		}

		sum += digit * multiplier

		multiplier++
		if multiplier > 7 {
			multiplier = 2
		}
	}

	remainder := 11 - (sum % 11)

	switch remainder {
	case 11:
		return "0", nil
	case 10:
		return "K", nil
	default:
		return strconv.Itoa(remainder), nil
	}
}
