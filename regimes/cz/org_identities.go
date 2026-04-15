package cz

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
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

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.Func("is IČO", isICOIdentity),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Czech IČO code",
					is.Func("valid", isValidICOCode),
				),
			),
		),
	)
}

func isICOIdentity(val any) bool {
	id, _ := val.(*org.Identity)
	return id != nil && id.Key == IdentityKeyICO
}

func isValidICOCode(val any) bool {
	code, _ := val.(cbc.Code)
	if code == "" || len(code.String()) != 8 {
		return false
	}
	return validateLegalEntityCode(code.String()) == nil
}

// normalizeIdentity normalizes IČO codes by stripping non-numeric characters.
func normalizeIdentity(id *org.Identity) {
	if id == nil || id.Key != IdentityKeyICO {
		return
	}
	id.Code = cbc.NormalizeNumericalCode(id.Code)
}
