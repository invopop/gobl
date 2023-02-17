package it

// The tax code here refers to Partita IVA, which is a distinct construct from
// Codice Fiscale. Italy operates with two types of tax identification codes.
// Though not all Italian persons possess Partita IVA, all parties engaged in
// economic activities within Italy are required to have one.

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

// validateTaxIdentity looks at the provided identity's code and performs the
// calculations required to determine if it is valid.
// These methods assume the code has already been normalized
// and thus only contains upper-case letters and numbers with
// no white space.
func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.Required, validation.By(validateTaxCode)))
}

// normalizeTaxIdentity removes any whitespace or separation characters and ensures all letters are
// uppercase.
func normalizeTaxIdentity(tID *tax.Identity) error {
	return common.NormalizeTaxIdentity(tID)
}

// source: https://it.wikipedia.org/wiki/Partita_IVA#Struttura_del_codice_identificativo_di_partita_IVA
func validateTaxCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return errors.New("code: cannot be blank")
	}

	for _, v := range code {
		x := v - 48
		if x < 0 || x > 9 {
			return errors.New("contains invalid characters")
		}
	}

	if len(code) != 11 {
		return errors.New("invalid length")
	}

	if computeLuhnCheckDigit(string(code[:10])) != int(code[10]-'0') {
		return errors.New("invalid check digit")
	}

	return nil
}

// Luhn Algorithm https://en.wikipedia.org/wiki/Luhn_algorithm
func computeLuhnCheckDigit(number string) int {
	sum := 0
	pos := 0

	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')

		if pos%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		pos++
	}

	return (10 - (sum % 10)) % 10
}
