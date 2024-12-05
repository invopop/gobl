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

	if err := hasValidChecksum(val); err != nil {
		return err
	}

	return nil
}

func hasValidChecksum(gstin string) error {
	if len(gstin) != 15 {
		return errors.New("invalid GSTIN length")
	}

	sum := 0
	for i, char := range gstin[:14] {
		value := charToValue(char)

		multiplier := 1
		if i%2 != 0 {
			multiplier = 2
		}

		product := value * multiplier
		sum += product/36 + product%36
	}

	remainder := sum % 36
	calculatedChecksum := (36 - remainder) % 36
	checksumChar := valueToChar(calculatedChecksum)

	if checksumChar != rune(gstin[14]) {
		return errors.New("checksum mismatch")
	}

	return nil
}

func charToValue(char rune) int {
	if char >= '0' && char <= '9' {
		return int(char - '0')
	}
	return int(char - 'A' + 10)
}

func valueToChar(value int) rune {
	if value >= 0 && value <= 9 {
		return rune('0' + value)
	}
	return rune('A' + value - 10)
}
