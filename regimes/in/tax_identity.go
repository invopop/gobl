// Package in provides the tax identity validation specific to India.
package in

import (
	"errors"
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var (
	taxCodeRegexp = regexp.MustCompile(`^[0-9]{2}[A-Z]{5}[0-9]{4}[A-Z]{1}[1-9A-Z]{1}Z[0-9A-Z]{1}$`)

	conversionTable = map[rune]int{
		'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
		'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15, 'G': 16, 'H': 17, 'I': 18,
		'J': 19, 'K': 20, 'L': 21, 'M': 22, 'N': 23, 'O': 24, 'P': 25, 'Q': 26, 'R': 27,
		'S': 28, 'T': 29, 'U': 30, 'V': 31, 'W': 32, 'X': 33, 'Y': 34, 'Z': 35,
	}
)

func normalizeTaxIdentity(tID *tax.Identity) {
	if tID == nil {
		return
	}
	tax.NormalizeIdentity(tID, l10n.IN)
	tID.Code = cbc.Code(strings.ToUpper(tID.Code.String()))
	tID.Country = "IN"
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

	if !taxCodeRegexp.MatchString(val) {
		return errors.New("invalid GSTIN format")
	}

	if !hasValidChecksum(val) {
		return errors.New("checksum mismatch")
	}

	return nil
}

func hasValidChecksum(gstin string) bool {
	if len(gstin) != 15 {
		return false
	}

	sum := 0
	for i, char := range gstin[:14] {
		value, exists := conversionTable[char]
		if !exists {
			return false
		}

		multiplier := 1
		if i%2 != 0 {
			multiplier = 2
		}

		product := value * multiplier
		sum += product/36 + product%36
	}

	remainder := sum % 36
	calculatedChecksum := (36 - remainder) % 36
	checksumChar := findCharByValue(calculatedChecksum)

	return checksumChar == rune(gstin[14])
}

func findCharByValue(value int) rune {
	for char, num := range conversionTable {
		if num == value {
			return char
		}
	}
	return ' '
}
