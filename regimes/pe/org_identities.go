package pe

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// identityTypeDefinitions lists the non-tax identity documents most commonly
// used to identify individuals on Peruvian invoices, from SUNAT's
// Catalogue 06. The catalogue's numeric codes are a concern of the electronic
// document formats, so they are left to a future addon.
// Source: https://www.sunat.gob.pe/legislacion/superin/2021/anexo-026-2021.pdf
// (Annex VI, SUNAT Catalogue 06).
var identityTypeDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeDNI,
		Name: i18n.String{
			i18n.EN: "National Identity Document",
			i18n.ES: "Documento Nacional de Identidad",
		},
		Desc: i18n.String{
			i18n.EN: "Eight-digit identity document issued by RENIEC to Peruvian citizens.",
			i18n.ES: "Documento de identidad de ocho dígitos emitido por RENIEC a los ciudadanos peruanos.",
		},
	},
	{
		Code: IdentityTypeCE,
		Name: i18n.String{
			i18n.EN: "Foreigner Resident Card",
			i18n.ES: "Carné de Extranjería",
		},
		Desc: i18n.String{
			i18n.EN: "Identity document issued by Migraciones to foreign residents in Peru.",
			i18n.ES: "Documento de identidad emitido por Migraciones a los extranjeros residentes en el Perú.",
		},
	},
	{
		Code: IdentityTypePassport,
		Name: i18n.String{
			i18n.EN: "Passport",
			i18n.ES: "Pasaporte",
		},
	},
}

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				org.IdentityTypeIn(IdentityTypeDNI),
				rules.Field("code",
					rules.Assert("01", "invalid DNI",
						is.Func("valid 8-digit DNI", isValidDNI),
					),
				),
			),
		),
	)
}

// normalizeOrgIdentity strips non-numeric characters from a DNI. Foreigner
// cards and passports don't have a single stable format, so they are left as
// provided.
func normalizeOrgIdentity(id *org.Identity) {
	if id.Type != IdentityTypeDNI {
		return
	}
	id.Code = cbc.NormalizeNumericalCode(id.Code)
}

// isValidDNI reports whether the value is a valid Peruvian DNI number:
// exactly eight digits. The physical document carries a check digit, but it
// is not part of the number itself, so only the format is validated here.
func isValidDNI(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok {
		return false
	}
	s := code.String()
	if len(s) != 8 {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
