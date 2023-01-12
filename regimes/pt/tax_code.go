package pt

import (
	"errors"
	"fmt"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

var (
	validPrefixes = map[string]bool{
		"1": true, "2": true, "3": true, "5": true, "6": true, "8": true,
		"45": true, "70": true, "71": true, "72": true, "74": true, "75": true,
		"77": true, "78": true, "79": true, "90": true, "91": true, "98": true, "99": true}
)

// validateTaxIdentity checks to ensure the NIT code looks okay.
func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.Required, validation.By(validateTaxCode)),
		validation.Field(&tID.Zone, isValidZoneCode),
	)
}

// normalizeTaxIdentity will remove any whitespace or separation characters from
// the tax code.
func normalizeTaxIdentity(tID *tax.Identity) error {
	return common.NormalizeTaxIdentity(tID)
}

var isValidZoneCode = validation.In(validZoneCodes()...)

func validZoneCodes() []interface{} {
	ls := make([]interface{}, len(zones))
	for i, v := range zones {
		ls[i] = v.Code
	}
	return ls
}

// based on example provided by https://pt.wikipedia.org/wiki/N%C3%BAmero_de_identifica%C3%A7%C3%A3o_fiscal
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
	if l != 9 {
		return errors.New("invalid length")
	}

	if !validPrefixes[code[:1]] && !validPrefixes[code[:2]] {
		return errors.New("invalid prefix")
	}

	// calculate check-digit
	sum := 0
	for i := 1; i < 9; i++ {
		v, err := strconv.Atoi(string(code[i-1]))
		if err != nil {
			return fmt.Errorf("invalid code: %w", err)
		}
		sum += v * (10 - i)
	}
	rmd := sum % 11
	ckd := 0
	switch rmd {
	case 0, 1:
		ckd = 0
	default:
		ckd = 11 - rmd
	}

	// compare the provided check digit with the calculated one
	compare, err := strconv.Atoi(string(code[8]))
	if err != nil {
		return fmt.Errorf("invalid check digit: %w", err)
	}
	if compare != ckd {
		return errors.New("checksum mismatch")
	}
	return nil
}
