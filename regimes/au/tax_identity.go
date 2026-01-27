package au

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

/*
 * ABN (Australian Business Number) Validation
 *
 * Format: 11 digits
 * Algorithm: Modulus 89 weighted checksum
 *
 * Sources:
 *  - https://abr.business.gov.au/Help/AbnFormat
 *
 * The ABN is calculated using a modulus 89 algorithm:
 * 1. Subtract 1 from the first (left-most) digit of the ABN
 * 2. Multiply each digit (including the modified first digit) by a weight
 *    Weights: [10, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19]
 * 3. Sum the resulting products
 * 4. Divide the sum by 89
 * 5. If the remainder is zero, the ABN is valid
 */

const (
	abnPattern = `^\d{11}$`
)

var (
	abnRegexp = regexp.MustCompile(abnPattern)
	// Weights for ABN checksum calculation
	abnWeights = [11]int{10, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19}
)

func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code,
			validation.By(validateABN),
		),
	)
}

func validateABN(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok {
		return nil
	}
	if code == "" {
		return nil
	}

	str := code.String()

	// Check format: must be exactly 11 digits
	if !abnRegexp.MatchString(str) {
		return errors.New("invalid format")
	}

	// Validate checksum using modulus 89 algorithm
	if !validateABNChecksum(str) {
		return errors.New("checksum mismatch")
	}

	return nil
}

// validateABNChecksum implements the modulus 89 ABN validation algorithm
func validateABNChecksum(abn string) bool {
	if len(abn) != 11 {
		return false
	}

	// Convert string to digit array
	digits := make([]int, 11)
	for i, char := range abn {
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return false
		}
		digits[i] = digit
	}

	// Step 1: Subtract 1 from the first digit
	digits[0]--

	// Step 2 & 3: Apply weights and sum
	sum := 0
	for i := 0; i < 11; i++ {
		sum += digits[i] * abnWeights[i]
	}

	// Step 4 & 5: Check if sum mod 89 equals 0
	return sum%89 == 0
}

// normalizeParty normalizes the tax ID in a party
func normalizeParty(party *org.Party) {
	if party == nil || party.TaxID == nil {
		return
	}
	normalizeABN(party.TaxID)
}

// normalizeABN cleans and formats an ABN
func normalizeABN(tID *tax.Identity) {
	if tID == nil || tID.Code == "" {
		return
	}
	tax.NormalizeIdentity(tID)
}
