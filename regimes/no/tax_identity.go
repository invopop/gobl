// Package no provides the tax identity validation specific to Norway.
package no

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var (
	// TRN in Norway is a 9-digit number with a checksum validation
	trnRegex = regexp.MustCompile(`^\d{9}$`)
)

// validateTaxIdentity checks to ensure the Norwegian TRN format is correct.
func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateTRNCode)),
	)
}

// validateTRNCode checks that the TRN is a valid 9-digit format with a checksum.
func validateTRNCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()

	// Check if TRN matches the 9-digit pattern
	if !trnRegex.MatchString(val) {
		return errors.New("must be a 9-digit number")
	}

	// Validate checksum
	if !validateChecksum(val) {
		return errors.New("invalid checksum for TRN")
	}

	return nil
}

// validateChecksum validates the TRN checksum using weight factors.
func validateChecksum(trn string) bool {
	// Use Norwegian organization number weights: [3,2,7,6,5,4,3,2] according to https://www.bits.no/document/standard-for-account-number-in-the-norwegian-banking-community-ver10/
	weights := []int{3, 2, 7, 6, 5, 4, 3, 2}
	sum := 0
	for i, r := range trn[:8] {
		digit, _ := strconv.Atoi(string(r))
		sum += digit * weights[i]
	}
	remainder := sum % 11
	checkDigit := 11 - remainder
	if checkDigit == 11 {
		checkDigit = 0
	} else if checkDigit == 10 {
		return false
	}
	lastDigit, _ := strconv.Atoi(string(trn[8]))
	return checkDigit == lastDigit
}
