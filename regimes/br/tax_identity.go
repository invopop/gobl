package br

import (
	"errors"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Reference: https://pt.wikipedia.org/wiki/Cadastro_Nacional_da_Pessoa_Jurídica

func validateTaxIdentity(tID *tax.Identity) error {
	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateTaxCode)),
	)
}

func validateTaxCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()

	// Verify length
	if len(val) != 14 {
		return errors.New("must have 14 digits")
	}

	// Verify first verification digit
	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	if err := verifyDigit(val, weights1, 12); err != nil {
		return err
	}

	// Verify second verification digit
	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	if err := verifyDigit(val, weights2, 13); err != nil {
		return err
	}

	return nil
}

func verifyDigit(cnpj string, weights []int, position int) error {
	sum := 0
	for i := 0; i < len(weights); i++ {
		digit, err := strconv.Atoi(string(cnpj[i]))
		if err != nil {
			return errors.New("must contain only digits")
		}
		sum += digit * weights[i]
	}

	remainder := sum % 11
	var expectedDigit int
	if remainder < 2 {
		expectedDigit = 0
	} else {
		expectedDigit = 11 - remainder
	}

	actualDigit, err := strconv.Atoi(string(cnpj[position]))
	if err != nil {
		return errors.New("must contain only digits")
	}

	if actualDigit != expectedDigit {
		return errors.New("verification digit mismatch")
	}

	return nil
}
