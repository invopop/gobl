package it

// The tax code here refers to Partita IVA, which is a distinct construct from
// Codice Fiscale. Italy operates with two types of tax identification codes.
// Though not all Italian persons possess Partita IVA, all parties engaged in
// economic activities within Italy are required to have one.

import (
	"errors"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

const (
	taxIDEvenChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	taxIDOddChars  = "BAKPLCQDREVOSFTGUHMINJWZYX"
	taxIDCharCode  = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	taxIDCRCMod    = 26
)

// validateTaxIdentity performs checks on the tax codes according to the type
// that was set. Additional validation is laid out at the invoice layer.
func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateTaxCode)),
	)
}

// source: https://it.wikipedia.org/wiki/Partita_IVA#Struttura_del_codice_identificativo_di_partita_IVA
func validateTaxCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	str := code.String()

	// Check code is just numbers
	for _, v := range str {
		x := v - 48
		if x < 0 || x > 9 {
			return errors.New("contains invalid characters")
		}
	}

	if len(str) != 11 {
		return errors.New("invalid length")
	}

	chk := computeLuhnCheckDigit(str[:10])
	if chk != str[10:] {
		return errors.New("invalid check digit")
	}

	return nil
}

// TODO: refactor this into a shareable method.
func computeLuhnCheckDigit(number string) string {
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

	return strconv.FormatInt(int64((10-(sum%10))%10), 10)
}
