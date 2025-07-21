package br

import (
	"errors"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Both CNPJ and CPF are supported as tax identifiers because, in Brazil, service invoices (NFS-e) may be issued by either legal entities (CNPJ) or individuals (CPF).
// CNPJ Reference: https://pt.wikipedia.org/wiki/Cadastro_Nacional_da_Pessoa_Jurídica
// CPF validation uses the same Mod11 algorithm as CNPJ, with different weights.
// Matches standard validators and test cases.
// Ref: https://pt.wikipedia.org/wiki/Cadastro_de_Pessoas_F%C3%ADsicas
// (Note: Some sources mention other variants, but Mod11 is the most consistent in practice.)

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

	switch len(val) {
	case 14:
		// CNPJ validation
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
	case 11:
		// CPF validation
		weights1 := []int{10, 9, 8, 7, 6, 5, 4, 3, 2}
		if err := verifyDigit(val, weights1, 9); err != nil {
			return err
		}
		weights2 := []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2}
		if err := verifyDigit(val, weights2, 10); err != nil {
			return err
		}
		//
	default:
		return errors.New("must have 11 (CPF) or 14 (CNPJ) digits")
	}
	return nil
}

func verifyDigit(val string, weights []int, position int) error {
	sum := 0
	for i := 0; i < len(weights); i++ {
		digit, err := strconv.Atoi(string(val[i]))
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

	actualDigit, err := strconv.Atoi(string(val[position]))
	if err != nil {
		return errors.New("must contain only digits")
	}

	if actualDigit != expectedDigit {
		return errors.New("verification digit mismatch")
	}

	return nil
}
