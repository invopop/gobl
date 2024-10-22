package pl

import (
	"errors"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

const (
	// IdentityKeyTaxNumber represents the Polish tax number (PESEL). It is not
	// required for invoices, but can be included for identification purposes.
	IdentityKeyTaxNumber cbc.Key = "pl-tax-number"
)

// Reference: https://en.wikipedia.org/wiki/PESEL

var identityKeyDefinitions = []*cbc.KeyDefinition{
	{
		Key: IdentityKeyTaxNumber,
		Name: i18n.String{
			i18n.EN: "Tax Number",
			i18n.PL: "Numer podatkowy",
		},
	},
}

func validateTaxNumber(id *org.Identity) error {
	if id == nil || id.Key != IdentityKeyTaxNumber {
		return nil
	}

	return validation.ValidateStruct(id,
		validation.Field(&id.Code, validation.By(validateTaxIdCode)),
	)
}

func validateTaxIdCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()

	if len(val) != 11 {
		return errors.New("length must be 11")
	}

	multipliers := []int{1, 3, 7, 9}
	sum := 0

	// Loop through the first 10 digits
	for i := 0; i < 10; i++ {
		digit, _ := strconv.Atoi(string(val[i]))
		sum += digit * multipliers[i%4]
	}

	modulo := sum % 10
	lastDigit, _ := strconv.Atoi(string(val[10]))

	if (modulo == 0 && lastDigit == 0) || lastDigit == 10-modulo {
		return nil
	}

	return errors.New("invalid checksum")
}
