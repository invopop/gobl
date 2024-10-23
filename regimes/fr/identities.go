package fr

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

// Identification keys used for additional codes not
// covered by the standard fields.
const (
	IdentityTypeSIREN cbc.Code = "SIREN" // SIREN is the main local tax code used in france, we use the normalized VAT version for the tax ID.
	IdentityTypeSIRET cbc.Code = "SIRET" // SIRET number combines the SIREN with a branch number.
	IdentityTypeRCS   cbc.Code = "RCS"   // Trade and Companies Register.
	IdentityTypeRM    cbc.Code = "RM"    // Directory of Traders.
	IdentityTypeNAF   cbc.Code = "NAF"   // Identifies the main branch of activity of the company or self-employed person.
	IdentityTypeSPI   cbc.Code = "SPI"   // Système de Pilotage des Indices
	IdentityTypeNIF   cbc.Code = "NIF"   // Numéro d'identification fiscale (people)
)

var (
	identityTypeSPIPattern = regexp.MustCompile(`^[0-3]\d{12}$`)
)

var badCharsRegexPattern = regexp.MustCompile(`[^\d]`)

var identityTypeDefinitions = []*cbc.ValueDefinition{
	{
		// https://www.oecd.org/content/dam/oecd/en/topics/policy-issue-focus/aeoi/france-tin.pdf
		Value: IdentityTypeSPI.String(),
		Name: i18n.String{
			i18n.EN: "Index Steering System",
			i18n.FR: "Système de Pilotage des Indices",
		},
	},
}

func normalizeIdentity(id *org.Identity) {
	if id == nil {
		return
	}
	switch id.Type {
	case IdentityTypeSPI:
		code := id.Code.String()
		code = badCharsRegexPattern.ReplaceAllString(code, "")
		id.Code = cbc.Code(code)
	}
}

// validateIdentity performs basic validation checks on identities provided.
func validateIdentity(id *org.Identity) error {
	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.By(identityValidator(id.Type)),
			validation.Skip,
		),
	)
}

func identityValidator(typ cbc.Code) validation.RuleFunc {
	return func(value interface{}) error {
		switch typ {
		case IdentityTypeSPI:
			return validation.Validate(value, validation.Match(identityTypeSPIPattern))
		//TODO: Add the other types
		default:
			return nil
		}
	}
}
