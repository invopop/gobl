package de

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
	// IdentityKeyTaxNumber represents the German tax number (Steuernummer) issued to
	// people that can be included on invoices inside Germany. For international
	// sales, the registered VAT number (Umsatzsteueridentifikationsnummer) should
	// be used instead.
	IdentityKeyTaxNumber cbc.Key = "de-tax-number"
	// IdentityKeyTaxID represents the German tax ID (Steuerliche Identifikationsnummer)
	IdentityKeyTaxID cbc.Key = "de-tax-id"
)

var identityKeyDefinitions = []*cbc.KeyDefinition{
	{
		Key: IdentityKeyTaxNumber,
		Name: i18n.String{
			i18n.EN: "Tax Number",
			i18n.DE: "Steuernummer",
		},
	},
	{
		Key: IdentityKeyTaxID,
		Name: i18n.String{
			i18n.EN: "Tax ID",
			i18n.DE: "Steuerliche Identifikationsnummer",
		},
	},
}

// Normalize will attempt to clean the object passed to it.
func normalizeIdentity(id *org.Identity) {
	if id == nil || (id.Key != IdentityKeyTaxNumber && id.Key != IdentityKeyTaxID) {
		return
	}
	code := id.Code.String()
	code = regexp.MustCompile(`[^\d]`).ReplaceAllString(code, "")
	id.Code = cbc.Code(code)
}

// Validate checks the document type and determines if it can be validated.
func validateIdentity(id *org.Identity) error {
	if id == nil || (id.Key != IdentityKeyTaxNumber && id.Key != IdentityKeyTaxID) {
		return nil
	}
	if id.Key == IdentityKeyTaxNumber {
		return validation.ValidateStruct(id,
			validation.Field(&id.Code,
				validation.Required,
				validation.By(validateTaxNumber),
				validation.Skip,
			),
		)
	}
	if id.Key == IdentityKeyTaxID {
		return validation.ValidateStruct(id,
			validation.Field(&id.Code,
				validation.Required,
				validation.By(validateTaxID),
				validation.Skip,
			),
		)
	}
	return nil
}

// Validation for German Steuernummer
func validateTaxNumber(value interface{}) error {
	val, ok := value.(cbc.Code)
	if !ok || val == cbc.CodeEmpty {
		return nil
	}
	code := val.String()
	if match, _ := regexp.MatchString(`^\d+$`, code); !match {
		return errors.New("invalid format: tax number should only contain digits")
	}

	if len(code) < 10 || len(code) > 13 {
		return errors.New("invalid length")
	}
	// Check if the first digit is not 0
	if code[0] == '0' {
		return errors.New("invalid format: first digit cannot be 0")
	}

	return nil
}

// Validation for German Steuerliche Identifikationsnummer
func validateTaxID(value interface{}) error {
	val, ok := value.(cbc.Code)
	if !ok || val == cbc.CodeEmpty {
		return nil
	}
	code := val.String()
	if match, _ := regexp.MatchString(`^\d+$`, code); !match {
		return errors.New("invalid format: tax number should only contain digits")
	}

	var codeNumber string
	switch len(code) {
	case 10, 11:
		codeNumber = code
	case 12, 13:
		codeNumber = code[2:]
	default:
		return errors.New("invalid length")
	}
	// Check if the first digit is not 0
	if codeNumber[0] == '0' {
		return errors.New("invalid format: first digit cannot be 0")
	}
	// Check digit occurrence rule for the first ten digits: https://de.wikipedia.org/wiki/Steuerliche_Identifikationsnummer
	digitCount := make(map[rune]int)
	for _, digit := range codeNumber[:10] {
		digitCount[digit]++
	}

	twiceOrThrice := false
	zeroOccurrence := 0
	singleOccurrence := 0

	for _, count := range digitCount {
		switch count {
		case 2, 3:
			if twiceOrThrice {
				return errors.New("invalid format: more than one digit appears twice or thrice")
			}
			twiceOrThrice = true
		case 0:
			zeroOccurrence++
		case 1:
			singleOccurrence++
		default:
			return errors.New("invalid format: a digit appears more than three times")
		}
	}

	if !twiceOrThrice || zeroOccurrence > 2 || singleOccurrence < 7 || singleOccurrence > 8 {
		return errors.New("invalid format: digit occurrence rule not satisfied")
	}

	// Extract digits and check digit ISO/IEC 7064 MOD 11,10
	digits := make([]int, len(codeNumber)-1)
	for i := 0; i < len(codeNumber)-1; i++ {
		digit, _ := strconv.Atoi(string(codeNumber[i]))
		digits[i] = digit
	}
	RealCheckDigit, _ := strconv.Atoi(string(codeNumber[len(codeNumber)-1]))
	prod := 10
	for _, digit := range digits {
		sum := (digit + prod) % 10
		if sum == 0 {
			sum = 10
		}
		prod = (sum * 2) % 11
	}
	checkDigit := 11 - prod
	if checkDigit == 10 {
		checkDigit = 0
	}
	if RealCheckDigit == checkDigit {
		return nil
	}

	return errors.New("checksum mismatch")
}
