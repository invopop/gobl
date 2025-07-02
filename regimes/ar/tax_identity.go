// Package ar provides CUIT validation and normalization for Argentina.
package ar

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Tax Identity Regexp
var (
	TaxCodeBadCharsRegexp = regexp.MustCompile(`\D`)
)

// ValidateTaxIdentity performs checks on the tax identity.
func ValidateTaxIdentity(tID *tax.Identity) error {
	if tID == nil {
		return validation.NewError("identity", "tax identity is required")
	}
	return validation.ValidateStruct(tID,
		validation.Field(
			&tID.Code,
			validation.Required,
			validation.By(ValidateTaxCode)),
	)
}

// ValidateTaxCode checks if the provided CUIT code is structurally valid.
func ValidateTaxCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}

	str := code.String()

	// Check code is just numbers
	for _, v := range str {
		x := v - 48
		if x < 0 || x > 9 {
			return errors.New("contains invalid characters")
		}
	}

	if len(str) != 11 {
		return errors.New("invalid length")
	}

	chk := ComputeModulo11CheckDigit(str[:10])
	if chk != str[10:] {
		return errors.New("invalid check digit")
	}

	return nil
}

// NormalizeTaxIdentity ensures the tax code is good for AR
func NormalizeTaxIdentity(tID *tax.Identity) {
	if tID == nil {
		return
	}
	tID.Code = NormalizeTaxCode(tID.Code)
}

// NormalizeTaxCode removes formatting characters and returns a cleaned CUIT.
func NormalizeTaxCode(code cbc.Code) cbc.Code {
	c := strings.ToUpper(code.String())
	c = TaxCodeBadCharsRegexp.ReplaceAllString(c, "")
	return cbc.Code(c)
}

// ComputeModulo11CheckDigit expects as argument a number string excluding the check
// digit. The returned integer should be checked against the check digit by the
// caller.
// Modulo 11 Algorithm definition: https://en.wikipedia.org/wiki/Check_digit
func ComputeModulo11CheckDigit(number string) string {
	weights := []int{5, 4, 3, 2, 7, 6, 5, 4, 3, 2}
	sum := 0

	for i := 0; i < 10; i++ {
		d := int(number[i] - '0')
		sum += d * weights[i]
	}

	mod := sum % 11
	var digit int
	switch mod {
	case 0:
		digit = 0
	case 1:
		digit = 9
	default:
		digit = 11 - mod
	}

	return strconv.Itoa(digit)
}
