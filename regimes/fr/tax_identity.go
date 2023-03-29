package fr

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var (
	taxCodeSIRENRegexp = regexp.MustCompile(`^\d{9}$`)
)

// validateTaxIdentity checks to ensure the SIRET code looks okay.
func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.Required, validation.By(validateSIRENTaxCode)),
	)
}

func validateSIRENTaxCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok {
		return nil
	}
	if code == "" {
		return nil
	}
	str := code.String()

	if !taxCodeSIRENRegexp.MatchString(str) {
		return errors.New("invalid format")
	}

	base := str[:8]
	chk := str[8:]
	v := common.ComputeLuhnCheckDigit(base)
	if chk != v {
		return errors.New("checksum mismatch")
	}

	return nil
}
