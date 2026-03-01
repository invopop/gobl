package fi

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var taxCodeRegexps = []*regexp.Regexp{
	regexp.MustCompile(`^\d{8}$`),
}

func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateTaxCode)),
	)
}

func validateTaxCode(value any) error {
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

// Finland's Y-tunnus (Business ID) check digit validation.
//
// Format: 7 digits + hyphen + check digit. Hyphen removed during normalization.
// Validation: MOD 11 with weights [7, 9, 10, 5, 8, 4, 2, 1].
//
// Digit conversion via val[i]-'0' assumes the input contains only
// ASCII digits. This is guaranteed by the regex validation in validateTaxCode.
//
// Reference: https://www.vero.fi/globalassets/tietoa-verohallinnosta/ohjelmistokehittajille/yritys--ja-yhteisötunnuksen-ja-henkilötunnuksen-tarkistusmerkin-tarkistuslaskenta.pdf
func validateTaxCodeChecksum(val string) error {
	weights := []int{7, 9, 10, 5, 8, 4, 2, 1}
	sum := 0

	for i := range 8 {
		sum += int(val[i]-'0') * weights[i]
	}

	if sum%11 != 0 {
		return errors.New("checksum mismatch")
	}

	return nil
}
