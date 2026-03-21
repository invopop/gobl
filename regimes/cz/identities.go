package cz

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

const (
	// IdentityKeyICO represents the Czech business registration number
	// (Identifikační číslo osoby) assigned to all businesses at formation.
	// For international trade, the DIČ (VAT number) should be used via
	// the TaxID field instead.
	IdentityKeyICO cbc.Key = "cz-ico"
)

var identityDefinitions = []*cbc.Definition{
	{
		Key: IdentityKeyICO,
		Name: i18n.String{
			i18n.EN: "Business Registration Number",
			i18n.CS: "Identifikační číslo osoby",
		},
	},
}

func validateIdentity(id *org.Identity) error {
	if id == nil || id.Key != IdentityKeyICO {
		return nil
	}
	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.Required,
			validation.By(validateICO),
			validation.Skip,
		),
	)
}

// validateICO validates an IČO using the same modulo-11 algorithm
// as the 8-digit legal entity DIČ (they share the same format).
func validateICO(value any) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	s := code.String()
	if len(s) != 8 {
		return errTaxCodeInvalidFormat
	}
	return validateLegalEntityCode(s)
}
