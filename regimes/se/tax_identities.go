package se

import (
	"errors"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

const (
	taxCodeCountryPrefix = "SE"
	taxCodeLength        = 14
	taxCodeCheckDigit    = "01"
)

// normalizeTaxIdentity performs normalization specific to Swedish tax IDs,
// ensuring the code is normalized and the country prefix is added if missing.
func normalizeTaxIdentity(id *tax.Identity) {
	tax.NormalizeIdentity(id)
	// TO-DO: decide if this is necessary. If not, we may need to remove the check digit suffix.
	// Re-add the SE prefix if missing.
	if id.Code.String()[:2] != taxCodeCountryPrefix {
		id.Code = cbc.Code(taxCodeCountryPrefix + id.Code.String())
	}
}

// validateTaxIdentity performs validation specific to Swedish tax IDs.
// Assumes the code has already been normalized.
func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code,
			validation.By(validateTaxCode),
		),
	)
}

// validateTaxCode validates the tax code for Swedish tax identities.
// Assumes the code has already been normalized, containing the country prefix
// and check digits.
func validateTaxCode(value any) error {
	code, ok := value.(cbc.Code)
	if !ok {
		return nil
	}
	if code == "" {
		return nil
	}
	// Normalised Swedish tax IDs must be 14 characters long.
	if len(code) != taxCodeLength {
		return errors.New("invalid length")
	}
	// Swedish tax IDs must start with "SE".
	if code[:2] != taxCodeCountryPrefix {
		return ErrInvalidTaxIDCountryPrefix
	}
	// Swedish tax IDs must finish in "01".
	if code[12:] != taxCodeCheckDigit {
		return errors.New("invalid check digit, expected 01")
	}
	// Swedish tax IDs must be exclusively numeric after the prefix.
	if _, err := strconv.Atoi(string(code[2:])); err != nil {
		return errors.New("invalid characters, expected numeric")
	}
	return nil
}
