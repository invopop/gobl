// Türkiye uses two tax identifiers: VKN (10-digit, companies) and TCKN
// (11-digit, individuals). Both are checksum-validated; the type is detected
// by length. The "TR" prefix and spaces are stripped during normalization.
// All Turkish invoices require the supplier to have a valid VKN or TCKN.
//
// Checksum algorithms based on:
// https://github.com/MhmtMutlu/tckn-vkn-validator

package tr

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var (
	// VKN (Vergi Kimlik Numarası): 10-digit company tax number issued by GIB,
	// required for all legal entities.
	vknRegexp = regexp.MustCompile(`^\d{10}$`)
	// TCKN (Türkiye Cumhuriyeti Kimlik Numarası): 11-digit national identity
	// number issued by the civil registry, used by sole traders and individuals.
	tcknRegexp = regexp.MustCompile(`^\d{11}$`)
)

var (
	errTaxCodeInvalidFormat   = errors.New("invalid format")
	errTaxCodeInvalidChecksum = errors.New("invalid check digit")
)

// normalizeTaxIdentity strips spaces and the "TR" country prefix, uppercases.
func normalizeTaxIdentity(tID *tax.Identity) {
	if tID == nil {
		return
	}
	tax.NormalizeIdentity(tID)
}

// validateTaxIdentity checks that the tax identity code is a valid VKN or TCKN.
func validateTaxIdentity(tID *tax.Identity) error {
	if tID == nil || tID.Code == cbc.CodeEmpty {
		return nil
	}
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateTaxCode)),
	)
}

// validateTaxCode detects VKN vs TCKN by length and verifies the checksum.
func validateTaxCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == cbc.CodeEmpty {
		return nil
	}
	s := code.String()
	switch {
	case vknRegexp.MatchString(s):
		return verifyVKN(s)
	case tcknRegexp.MatchString(s):
		return verifyTCKN(s)
	default:
		return errTaxCodeInvalidFormat
	}
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
		return errTaxCodeInvalidChecksum
	}
	return nil
}

// verifyTCKN validates the 11-digit Turkish national identity number checksum.
//
// Algorithm:
//  1. First digit must not be 0.
//  2. Digit 10 = (sum of odd-indexed digits * 7 - sum of even-indexed digits) % 10
//     (indices 0,2,4,6,8 are odd positions; 1,3,5,7 are even positions)
//  3. Digit 11 = sum of first 10 digits % 10
func verifyTCKN(s string) error {
	digits := stringToDigits(s)
	if digits[0] == 0 {
		return errTaxCodeInvalidFormat
	}
	oddSum := digits[0] + digits[2] + digits[4] + digits[6] + digits[8]
	evenSum := digits[1] + digits[3] + digits[5] + digits[7]
	d10 := (oddSum*7 - evenSum) % 10
	if d10 < 0 {
		d10 += 10
	}
	if d10 != digits[9] {
		return errTaxCodeInvalidChecksum
	}
	total := 0
	for i := 0; i < 10; i++ {
		total += digits[i]
	}
	if total%10 != digits[10] {
		return errTaxCodeInvalidChecksum
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
