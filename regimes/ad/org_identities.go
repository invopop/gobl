package ad

import (
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// IdentityTypeNRT represents the Número de Registre Tributari,
	// the sole tax identifier issued to all taxpayers in Andorra.
	IdentityTypeNRT cbc.Code = "NRT"
)

// identityDefinitions is referenced by ad.go and describes the identity
// types recognised in Andorra.
var identityDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeNRT,
		Name: i18n.String{
			i18n.EN: "Tax Register Number (NRT)",
			i18n.CA: "Número de Registre Tributari (NRT)",
			i18n.ES: "Número de Registro Tributario (NRT)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The Número de Registre Tributari (NRT) is the sole tax identifier
				issued to all taxpayers in Andorra — resident and non-resident,
				individuals and legal entities — by the Departament de Tributs i
				de Fronteres.

				The first letter encodes the taxpayer category: F for resident
				individuals, E for non-resident individuals, A for joint-stock
				companies (SA), L for limited liability companies (SL), and other
				letters for cooperatives, foundations, public bodies and
				special-purpose entities.

				The NRT is used for all tax declarations including IGI, corporate
				tax, and customs. Entities below the €40,000 annual IGI threshold
				hold an NRT but are not required to register for IGI.
			`),
		},
	},
}

// normalizeIdentity cleans NRT input before validation runs.
// Strips hyphens and spaces, uppercases, and removes any "AD" country prefix.
// Examples:
//
//	L-132950-X  →  L132950X
//	l132950x    →  L132950X
//	ADL132950X  →  L132950X
//  NRTL132950X →  L132950X

func normalizeIdentity(id *org.Identity) {
	if id == nil || id.Type != IdentityTypeNRT {
		return
	}
	code := strings.ToUpper(id.Code.String())
	code = tax.IdentityCodeBadCharsRegexp.ReplaceAllString(code, "")
	code = strings.TrimPrefix(code, string(l10n.AD)) // strip "AD" country prefix
	code = strings.TrimPrefix(code, "NRT")           // strip "NRT" label prefix
	id.Code = cbc.Code(code)
}

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				org.IdentityTypeIn(IdentityTypeNRT),
				rules.Field("code",
					rules.Assert("01", "identity code for type NRT must be valid",
						is.Func("valid NRT", orgIdentityCheckNRT),
					),
				),
			),
		),
	)
}

func orgIdentityCheckNRT(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return false
	}
	return reNRT.MatchString(code.String())
}