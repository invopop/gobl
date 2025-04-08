package fr

import (
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

const (
	// SIREN is the main local tax code used in france, we use the normalized VAT version for the tax ID.
	IdentityKeySiren cbc.Key = "fr-siren"
	// SIRET is the SIREN with a branch number.
	IdentityKeySiret cbc.Key = "fr-siret"
)

var (
	taxCodeSIRENRegexp = regexp.MustCompile(`^\d{9}$`)
	taxCodeSIRETRegexp = regexp.MustCompile(`^\d{14}$`)
)

var identityKeyDefinitions = []*cbc.Definition{
	{
		Key: IdentityKeySiren,
		Name: i18n.String{
			i18n.EN: "SIREN",
			i18n.FR: "SIREN",
		},
	},
	{
		Key: IdentityKeySiret,
		Name: i18n.String{
			i18n.EN: "SIRET",
			i18n.FR: "SIRET",
		},
	},
}

func normalizeIdentity(id *org.Identity) {
	if id == nil || (id.Key != IdentityKeySiren && id.Key != IdentityKeySiret) {
		return
	}
	code := strings.ToUpper(id.Code.String())
	code = tax.IdentityCodeBadCharsRegexp.ReplaceAllString(code, "")
	code = strings.TrimPrefix(code, string(l10n.FR))
	id.Code = cbc.Code(code)
}

func validateIdentity(id *org.Identity) error {
	if id == nil || (id.Key != IdentityKeySiren && id.Key != IdentityKeySiret) {
		return nil
	}

	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.Required,
			validation.By(validateSIRENTaxCode),
			validation.Skip,
		),
	)
}
