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

// Valid formats: 2/3/5 (10 digits), 3/3/5 (11 digits standard), or 3/4/4 (11 digits NW)
// See: https://de.wikipedia.org/wiki/Steuernummer
var taxNumberRegexPattern = regexp.MustCompile(`^(\d{2}/\d{3}/\d{5}|\d{3}/\d{3}/\d{5}|\d{3}/\d{4}/\d{4})$`)

var identityDefinitions = []*cbc.Definition{
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

	// If already matches the regex, it's already in valid format
	if taxNumberRegexPattern.MatchString(id.Code.String()) {
		return
	}

	// Normalize to standard format
	code := cbc.NormalizeNumericalCode(id.Code).String()
	if len(code) == 11 {
		// If 11 digits, return the standard format 123/456/78901 (3/3/5)
		code = fmt.Sprintf("%s/%s/%s", code[:3], code[3:6], code[6:])
	} else if len(code) == 10 {
		// If 10 digits, return the format 12/345/67890 (2/3/5)
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
