package ad

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// NRT (Número de Registre Tributari) format:
// 1 letter prefix (A, C, D, F, L, U, V) indicating entity type,
// followed by 6 digits, and 1 uppercase control letter.
var nrtRegex = regexp.MustCompile(`^[ACDFLUV]\d{6}[A-Z]$`)

func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateNRT)),
	)
}

func validateNRT(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == cbc.CodeEmpty {
		return nil
	}
	if !nrtRegex.MatchString(string(code)) {
		return errors.New("invalid NRT format")
	}
	return nil
}
