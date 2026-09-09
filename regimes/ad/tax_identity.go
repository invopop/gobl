package ad

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// NRT (Número de Registre Tributari) format:
// A leading letter that can be:
//
//	F: Individual Residents
//	E: Non-resident Individuals
//	L: Limited Liability Companies (S.L.)
//	A: Joint-stock Corporations (S.A.)
//
// Followed by six digits, and ending with a control letter.
// Example: L-123456-A (often displayed with hyphens)
var (
	nrtRegexp = regexp.MustCompile(`^[AEFL][0-9]{6}[A-Z]$`)
)

func validateTaxIdentity(t *tax.Identity) error {
	if t == nil {
		return errors.New("tax identity cannot be nil")
	}
	return validation.ValidateStruct(t,
		validation.Field(&t.Code,
			validation.Required,
			validation.Match(nrtRegexp),
		),
	)
}
