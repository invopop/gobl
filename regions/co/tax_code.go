package co

import (
	"errors"
	"fmt"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regions/common"
)

var (
	nitMultipliers = []int{3, 7, 13, 17, 19, 23, 29, 37, 41, 43, 47, 53, 59, 67, 71}
)

// ValidateTaxIdentity checks to ensure the NIT code looks okay.
func ValidateTaxIdentity(tID *org.TaxIdentity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateTaxCode)),
	)
}

// NormalizeTaxIdentity will remove any whitespace or separation characters from
// the tax code.
func NormalizeTaxIdentity(tID *org.TaxIdentity) error {
	if err := common.NormalizeTaxIdentity(tID); err != nil {
		return err
	}
	return nil
}

func validateTaxCode(value interface{}) error {
	code, ok := value.(string)
	if !ok {
		return nil
	}
	if code == "" {
		return nil
	}
	for _, v := range code {
		x := v - 48
		if x < 0 || x > 9 {
			return errors.New("contains invalid characters")
		}
	}
	l := len(code)
	if l > 10 {
		return errors.New("too long")
	}
	if l < 9 {
		return errors.New("too short")
	}

	return validateDigits(code[0:l-1], code[l-1:l])
}

func validateDigits(code, check string) error {
	ck, err := strconv.Atoi(check)
	if err != nil {
		return fmt.Errorf("invalid check: %w", err)
	}

	sum := 0
	l := len(code)
	for i, v := range code {
		// 48 == ASCII "0"
		sum += int(v-48) * nitMultipliers[l-i-1]
	}
	sum = sum % 11
	if sum >= 2 {
		sum = 11 - sum
	}

	if sum != ck {
		return errors.New("checksum mismatch")
	}

	return nil
}
