//go:build ignore

package template

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var taxCodePatterns = []*regexp.Regexp{
	regexp.MustCompile(`^\d{8}$`), // Replace with actual tax ID pattern
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
	for _, re := range taxCodePatterns {
		if re.MatchString(val) {
			match = true
			break
		}
	}
	if !match {
		return errors.New("invalid format")
	}

	return nil
}
