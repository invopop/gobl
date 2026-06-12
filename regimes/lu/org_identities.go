package lu

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// IdentityTypeRCS represents a Luxembourg Registre de Commerce et des Sociétés
	// (RCS) number, assigned by the Luxembourg Business Registers (LBR) to
	// businesses registered in Luxembourg.
	//
	// Format: one letter (B, F, G, or H) followed by up to six digits.
	// Examples: B263475, F12345, G45678
	//
	// Source: https://www.lbr.lu
	IdentityTypeRCS cbc.Code = "RCS"
)

// rcsRegexp validates a normalised RCS number: register letter + 1–6 digits.
var rcsRegexp = regexp.MustCompile(`^[BFGH]\d{1,6}$`)

var identityTypeDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeRCS,
		Name: i18n.String{
			i18n.EN: "RCS Number",
			i18n.FR: "Numéro RCS",
			i18n.LB: "RCS-Nummer",
		},
		Desc: i18n.String{
			i18n.EN: "Luxembourg company registration number assigned by the Luxembourg Business Registers (LBR).",
			i18n.FR: "Numéro d'immatriculation au Registre de Commerce et des Sociétés attribué par le Registre luxembourgeois des entreprises.",
		},
	},
}

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				org.IdentityTypeIn(IdentityTypeRCS),
				rules.Field("code",
					rules.Assert("01", "invalid RCS number",
						is.Func("valid RCS format", isValidRCSCode),
					),
				),
			),
		),
	)
}

// normalizeOrgIdentity strips spaces and non-alphanumeric characters from an
// RCS number and uppercases it (e.g. "b 263 475" → "B263475").
func normalizeOrgIdentity(id *org.Identity) {
	if id.Type != IdentityTypeRCS {
		return
	}
	id.Code = cbc.NormalizeAlphanumericalCode(id.Code)
}

func isValidRCSCode(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return false
	}
	return rcsRegexp.MatchString(code.String())
}
