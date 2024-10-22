package at

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
	// IdentityKeyTaxNumber represents the Austrian tax number (Steuernummer) issued to
	// people that can be included on invoices inside Austria. For international
	// sales, the registered VAT number (Umsatzsteueridentifikationsnummer) should
	// be used instead.
	IdentityKeyTaxNumber cbc.Key = "at-tax-number"
)

// https://www.oecd.org/content/dam/oecd/en/topics/policy-issue-focus/aeoi/austria-tin.pdf

var badCharsRegexPattern = regexp.MustCompile(`[^\d]`)

var identityKeyDefinitions = []*cbc.KeyDefinition{
	{
		Key: IdentityKeyTaxNumber,
		Name: i18n.String{
			i18n.EN: "Tax Number",
			i18n.DE: "Steuernummer",
		},
	},
}

func normalizeTaxNumber(id *org.Identity) {
	if id == nil || id.Key != IdentityKeyTaxNumber {
		return
	}
	code := id.Code.String()
	code = badCharsRegexPattern.ReplaceAllString(code, "")
	id.Code = cbc.Code(code)
}

func validateTaxNumber(id *org.Identity) error {
	if id == nil || id.Key != IdentityKeyTaxNumber {
		return nil
	}

	return validation.ValidateStruct(id,
		validation.Field(&id.Code, validation.By(validateTaxIDCode)),
	)
}

// validateAustrianTaxIdCode validates the normalized tax ID code.
func validateTaxIDCode(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()

	// Austrian Steuernummer format: must have 9 digits (2 for tax office + 7 for taxpayer ID)
	if len(val) != 9 {
		return errors.New("length must be 9 digits")
	}

	// Split into tax office code and taxpayer number
	taxOffice, _ := strconv.Atoi(val[:2])
	taxpayerNumber, _ := strconv.Atoi(val[2:])

	// Perform basic checks
	if taxOffice < 1 || taxOffice > 99 {
		return errors.New("invalid tax office code")
	}

	if taxpayerNumber <= 0 {
		return errors.New("invalid taxpayer number")
	}

	return nil
}
