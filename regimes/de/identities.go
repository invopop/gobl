package de

import (
	"fmt"
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
)

var taxNumberRegexPattern = regexp.MustCompile(`^\d{2,3}/\d{3}/\d{5}$`)
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

// Normalize for German Steuernummer
func normalizeTaxNumber(id *org.Identity) {
	if id == nil || id.Key != IdentityKeyTaxNumber {
		return
	}
	code := id.Code.String()
	code = badCharsRegexPattern.ReplaceAllString(code, "")
	if len(code) == 11 {
		// If 11 digits, return the format 123/456/78901
		code = fmt.Sprintf("%s/%s/%s", code[:3], code[3:6], code[6:])
	} else if len(code) == 10 {
		// If 10 digits, return the format 12/345/67890
		code = fmt.Sprintf("%s/%s/%s", code[:2], code[2:5], code[5:])
	}
	id.Code = cbc.Code(code)
}

// Validation for German Steuernummer
func validateTaxNumber(id *org.Identity) error {
	if id == nil || id.Key != IdentityKeyTaxNumber {
		return nil
	}
	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.Required,
			validation.Match(taxNumberRegexPattern),
			validation.Skip,
		),
	)
}
