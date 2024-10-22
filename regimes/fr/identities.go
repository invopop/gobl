package fr

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

const (
	// IdentityKeyTaxNumber represents the French tax reference number (numéro fiscal de référence).
	IdentityKeyTaxNumber cbc.Key = "fr-tax-number"
)

// https://www.oecd.org/content/dam/oecd/en/topics/policy-issue-focus/aeoi/france-tin.pdf

var badCharsRegexPattern = regexp.MustCompile(`[^\d]`)

var identityKeyDefinitions = []*cbc.KeyDefinition{
	{
		Key: IdentityKeyTaxNumber,
		Name: i18n.String{
			i18n.EN: "Tax Number",
			i18n.FR: "Numéro fiscal de référence",
		},
	},
}

// validateTaxNumber validates the French tax reference number.
func validateTaxNumber(id *org.Identity) error {
	if id == nil || id.Key != IdentityKeyTaxNumber {
		return nil
	}

	return validation.ValidateStruct(id,
		validation.Field(&id.Code, validation.By(validateTaxIDCode)),
	)
}

// validateTaxIDCode validates the normalized tax ID code.
func validateTaxIDCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()

	// Check length
	if len(val) != 13 {
		return errors.New("length must be 13 digits")
	}

	// Check that all characters are digits
	if _, err := strconv.Atoi(val); err != nil {
		return errors.New("must contain only digits")
	}

	// Check that the first digit is 0, 1, 2, or 3
	firstDigit := val[0]
	if firstDigit < '0' || firstDigit > '3' {
		return errors.New("first digit must be 0, 1, 2, or 3")
	}

	return nil
}

// normalizeTaxNumber removes any non-digit characters from the tax number.
func normalizeTaxNumber(id *org.Identity) {
	if id == nil || id.Key != IdentityKeyTaxNumber {
		return
	}
	code := id.Code.String()
	code = badCharsRegexPattern.ReplaceAllString(code, "")
	id.Code = cbc.Code(code)
}
