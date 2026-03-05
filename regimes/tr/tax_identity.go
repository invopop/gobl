package tr

// Checksum algorithm based on:
// https://github.com/MhmtMutlu/tckn-vkn-validator

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// VKN (Vergi Kimlik Numarası): 10-digit company tax number issued by GIB,
// required for all legal entities.
var vknRegexp = regexp.MustCompile(`^\d{10}$`)

var (
	errInvalidFormat   = errors.New("invalid format")
	errInvalidChecksum = errors.New("invalid check digit")
)

// normalizeTaxIdentity standardizes the tax identity code.
func normalizeTaxIdentity(tID *tax.Identity) {
	if tID == nil {
		return
	}
	tax.NormalizeIdentity(tID)
}

// validateTaxIdentity checks that the tax identity code is a valid VKN.
func validateTaxIdentity(tID *tax.Identity) error {
	if tID == nil || tID.Code == cbc.CodeEmpty {
		return nil
	}
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateVKNCode)),
	)
}

func validateVKNCode(value interface{}) error {
	code, _ := value.(cbc.Code)
	if code == cbc.CodeEmpty {
		return nil
	}
	s := code.String()
	if !vknRegexp.MatchString(s) {
		return errInvalidFormat
	}
	return verifyVKN(s)
}

// verifyVKN validates the 10-digit Turkish company tax number checksum.
//
// Algorithm:
//  1. For each of the first 9 digits (position i, 0-indexed):
//     v = (digit + (9 - i)) % 10
//     product = v * 2^(9-i)
//     contribution = product % 9, or 9 if product > 0 and product % 9 == 0
//  2. Check digit = (10 - sum % 10) % 10
func verifyVKN(s string) error {
	digits := stringToDigits(s)
	sum := 0
	for i := 0; i < 9; i++ {
		v := (digits[i] + (9 - i)) % 10
		product := v * pow2(9-i)
		contrib := product % 9
		if product > 0 && product%9 == 0 {
			contrib = 9
		}
		sum += contrib
	}
	check := (10 - sum%10) % 10
	if check != digits[9] {
		return errInvalidChecksum
	}
	return nil
}

func stringToDigits(s string) []int {
	digits := make([]int, len(s))
	for i, c := range s {
		digits[i] = int(c - '0')
	}
	return digits
}

func pow2(n int) int {
	result := 1
	for i := 0; i < n; i++ {
		result *= 2
	}
	return result
}
