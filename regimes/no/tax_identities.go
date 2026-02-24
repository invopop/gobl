package no

import (
	"errors"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateTaxIdentity(tID *tax.Identity) error {
	if tID == nil {
		return nil
	}
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code,
			validation.By(validateTaxCode),
		),
	)
}

func validateTaxCode(value any) error {
	code, _ := value.(cbc.Code)
	if code == "" {
		return nil
	}

	digits, ok := cleanNorwayTaxCode(code)
	if !ok {
		return errors.New("invalid length")
	}
	if !validOrgNrMod11(digits) {
		return errors.New("invalid checksum")
	}
	return nil
}

// normalizeTaxIdentity converts common representations like "NO#########MVA" into the
// canonical internal representation: 9 digits (orgnr).
func normalizeTaxIdentity(tID *tax.Identity) {
	if tID == nil {
		return
	}

	// Apply common normalization first (trims, normalizes country, etc.)
	tax.NormalizeIdentity(tID)

	if tID.Code == "" {
		return
	}

	digits, ok := cleanNorwayTaxCode(tID.Code)
	if ok {
		tID.Code = cbc.Code(digits)
	}
}
