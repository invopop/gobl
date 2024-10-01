package de

import (
	"errors"
	"regexp"

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
	IdentityKeyTaxID     cbc.Key = "de-tax-id"
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
	if id == nil || id.Key != IdentityKeyTaxNumber {
		return
	}
	code := id.Code.String()
	code = regexp.MustCompile(`[^\d]`).ReplaceAllString(code, "")
	id.Code = cbc.Code(code)
}

// Validate checks the document type and determines if it can be validated.
func validateIdentity(id *org.Identity) error {
	if id == nil || id.Key != IdentityKeyTaxNumber {
		return nil
	}
	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.Required,
			validation.By(validateIDNumber),
			validation.Skip,
		),
	)
}

// ValidateTaxNumber checks the document type and determines if it can be validated.
func validateIDNumber(value interface{}) error {
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
