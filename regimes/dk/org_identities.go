package dk

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// IdentityTypeCVR represents the Danish "CVR-nummer" (Centrale Virksomhedsregister),
	// the Central Business Register number used to identify businesses in Denmark.
	IdentityTypeCVR cbc.Code = "CVR"

	// IdentityTypeCPR represents the Danish "CPR-nummer" (Det Centrale Personregister),
	// the personal identification number used to identify individuals in Denmark.
	IdentityTypeCPR cbc.Code = "CPR"
)

var identityTypeDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeCVR,
		Name: i18n.String{
			i18n.EN: "CVR Number",
			i18n.DA: "CVR-nummer",
		},
	},
	{
		Code: IdentityTypeCPR,
		Name: i18n.String{
			i18n.EN: "CPR Number",
			i18n.DA: "CPR-nummer",
		},
	},
}

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				org.IdentityTypeIn(IdentityTypeCVR),
				rules.Field("code",
					rules.AssertIfPresent("01", "invalid Danish CVR identity code",
						is.Func("valid", isValidCVRCode)),
				),
			),
			rules.When(
				org.IdentityTypeIn(IdentityTypeCPR),
				rules.Field("code",
					rules.AssertIfPresent("02", "invalid Danish CPR identity code",
						is.Func("valid", isValidCPRCode)),
				),
			),
		),
	)
}

func isValidCVRCode(val any) bool {
	code, _ := val.(cbc.Code)
	return validateTaxCode(code) == nil
}

// isValidCPRCode checks a Danish CPR number is 10 digits. There is no
// checksum or other validation.
func isValidCPRCode(val any) bool {
	return is.Length(10, 10).Check(val) && is.Digit.Check(val)
}

// normalizeOrgIdentity strips separators from a Danish CPR number, e.g.
// "150605-4321" becomes "1506054321".
func normalizeOrgIdentity(id *org.Identity) {
	if id.Type != IdentityTypeCPR {
		return
	}
	id.Code = cbc.NormalizeNumericalCode(id.Code)
}
