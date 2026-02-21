package hr

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// The personal identification number consists of eleven digits,
// ten digits that are determined at random and one control digit
// Reference: https://narodne-novine.nn.hr/clanci/sluzbeni/2008_05_60_2033.html
var (
	oibRegexp = regexp.MustCompile(`^\d{11}$`)
)

// validateTaxIdentity checks to ensure the Croatian OIB tax code is valid.
func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateTaxCode)),
	)
}

// validateTaxCode validates the format and checksum of a Croatian OIB code.
func validateTaxCode(value any) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()

	if !oibRegexp.MatchString(val) {
		return errors.New("invalid format, expected 11 digits")
	}

	return validateOIBChecksum(val)
}

// validateOIBChecksum verifies the OIB check digit using the ISO 7064 MOD 11.10 algorithm.
// Reference: https://narodne-novine.nn.hr/clanci/sluzbeni/2009_01_1_6.html
func validateOIBChecksum(val string) error {
	r := 10
	for i := range 10 {
		a := (r + int(val[i]-'0')) % 10
		if a == 0 {
			a = 10
		}
		r = (a * 2) % 11
	}
	check := 11 - r
	if check == 10 {
		check = 0
	}
	if check != int(val[10]-'0') {
		return errors.New("invalid checksum")
	}
	return nil
}
