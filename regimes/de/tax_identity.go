package de

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Reference: https://github.com/ltns35/go-vat/blob/main/countries/germany.go

var (
	taxCodeRegexps = []*regexp.Regexp{
		regexp.MustCompile(`^[1-9]\d{8}$`),
	}
)

// validateTaxIdentity checks to ensure the NIT code looks okay.
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

	match := false
	for _, re := range taxCodeRegexps {
		if re.MatchString(val) {
			match = true
			break
		}
	}
	if !match {
		return errors.New("invalid format")
	}

	return validateTaxCodeChecksum(val)
}

func validateTaxCodeChecksum(val string) error {
	p := 10
	sum := 0
	cd := 0
	for i := 0; i < 8; i++ {
		digit, err := strconv.Atoi(string(val[i]))
		if err != nil {
			return errors.New("invalid digit")
		}
		sum = (digit + p) % 10
		if sum == 0 {
			sum = 10
		}
		p = (2 * sum) % 11
	}

	if 11-p == 10 {
		cd = 0
	} else {
		cd = 11 - p
	}

	ecd, err := strconv.Atoi(string(val[8]))
	if err != nil {
		return errors.New("invalid checksum")
	}
	if cd != ecd {
		return errors.New("checksum mismatch")
	}

	return nil
}
