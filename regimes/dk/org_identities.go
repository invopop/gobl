package dk

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

const (
	// IdentityTypeCVR represents the Danish "CVR-nummer" (Centrale Virksomhedsregister),
	// the Central Business Register number used to identify businesses in Denmark.
	IdentityTypeCVR cbc.Code = "CVR"
)

var identityTypeDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeCVR,
		Name: i18n.String{
			i18n.EN: "CVR Number",
			i18n.DA: "CVR-nummer",
		},
	},
}

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.Func("is CVR", isCVRIdentity),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Danish CVR identity code",
					is.Func("valid", isValidCVRCode)),
			),
		),
	)
}

func isCVRIdentity(val any) bool {
	id, _ := val.(*org.Identity)
	return id != nil && id.Type == IdentityTypeCVR
}

func isValidCVRCode(val any) bool {
	code, _ := val.(cbc.Code)
	return validateTaxCode(code) == nil
}
