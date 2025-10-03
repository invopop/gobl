package verifactu

import (
	"fmt"

	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

func validateParty(p *org.Party) error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Name,
			validation.By(validateNoForbiddenChars),
			validation.Skip,
		),
	)
}

// forbiddenChars contains characters that are not allowed in certain string fields
var forbiddenChars = []rune{'<', '>', '"', '\'', '='}

// validateNoForbiddenChars validates that a string doesn't contain any of the forbidden characters: < > " ' =
func validateNoForbiddenChars(val any) error {
	str, _ := val.(string)

	for _, char := range str {
		for _, forbidden := range forbiddenChars {
			if char == forbidden {
				return fmt.Errorf("contains forbidden character: %c", char)
			}
		}
	}

	return nil
}
