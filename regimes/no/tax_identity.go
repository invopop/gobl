package no

import (
	"errors"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// taxCodeWeights are the mod-11 multipliers for Norwegian organisasjonsnummer
// validation, as specified by Brønnøysundregistrene.
// See: https://www.brreg.no/en/about-us-2/our-registers/about-the-central-coordinating-register-for-legal-entities-ccr/about-the-organisation-number/
var taxCodeWeights = []int{3, 2, 7, 6, 5, 4, 3, 2}

// normalizeTaxIdentity performs standard tax identity normalization, and then
// removes the "MVA" suffix common in Norwegian VAT numbers.
func normalizeTaxIdentity(tID *tax.Identity) {
	if tID == nil {
		return
	}
	tax.NormalizeIdentity(tID)
	// Strip the "MVA" suffix common in Norwegian VAT numbers (e.g. "NO 923 456 783 MVA").
	tID.Code = cbc.Code(strings.TrimSuffix(string(tID.Code), "MVA"))
}

// validateTaxIdentity checks the Norwegian organisasjonsnummer using mod-11.
func validateTaxIdentity(tID *tax.Identity) error {
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

	if len(code) != 9 {
		return errors.New("must have 9 digits")
	}

	// All characters must be digits.
	for _, r := range code {
		if r < '0' || r > '9' {
			return errors.New("must only contain digits")
		}
	}

	// First digit must be 8 or 9.
	if code[0] != '8' && code[0] != '9' {
		return errors.New("first digit must be 8 or 9")
	}

	// Mod-11 check digit validation.
	sum := 0
	for i, w := range taxCodeWeights {
		sum += int(code[i]-'0') * w
	}
	remainder := sum % 11
	var check int
	if remainder == 0 {
		check = 0
	} else {
		check = 11 - remainder
	}
	if check == 10 {
		return errors.New("invalid check digit")
	}
	if int(code[8]-'0') != check {
		return errors.New("checksum mismatch")
	}

	return nil
}
