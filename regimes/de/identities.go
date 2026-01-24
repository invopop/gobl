package de

import (
	"fmt"
	"regexp"
	"strings"

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

// Valid formats: 2/3/5 (10 digits), 3/3/5 (11 digits standard), or 3/4/4 (11 digits NRW)
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

	// Check if input already has the NRW format (3/4/4)
	// If so, preserve it as-is
	original := id.Code.String()
	if strings.Count(original, "/") == 2 {
		parts := strings.Split(original, "/")
		if len(parts) == 3 {
			// Extract only digits from each part
			p1 := cbc.NormalizeNumericalCode(cbc.Code(parts[0])).String()
			p2 := cbc.NormalizeNumericalCode(cbc.Code(parts[1])).String()
			p3 := cbc.NormalizeNumericalCode(cbc.Code(parts[2])).String()

			// Check if it matches NRW format (3/4/4)
			if len(p1) == 3 && len(p2) == 4 && len(p3) == 4 {
				id.Code = cbc.Code(fmt.Sprintf("%s/%s/%s", p1, p2, p3))
				return
			}
		}
	}

	// Otherwise, normalize to standard format
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
