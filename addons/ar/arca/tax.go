package arca

import (
	"github.com/invopop/gobl/regimes/ar"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func normalizeTaxCombo(tc *tax.Combo) {
	if tc == nil {
		return
	}

	if tc.Key == tax.KeyZero {
		tc.Ext = tc.Ext.Set(ExtKeyVATRate, "3")
		return
	}

	switch tc.Rate {
	case tax.RateReduced:
		tc.Ext = tc.Ext.Set(ExtKeyVATRate, "4")
	case tax.RateGeneral:
		tc.Ext = tc.Ext.Set(ExtKeyVATRate, "5")
	case ar.RateIncreased:
		tc.Ext = tc.Ext.Set(ExtKeyVATRate, "6")
	}
}

func taxComboRules() *rules.Set {
	return rules.For(new(tax.Combo),
		rules.When(is.Func("vat", taxComboIsVAT),
			rules.Field("ext",
				rules.Assert("01", "ar-arca-vat-rate: required", tax.ExtensionsRequire(ExtKeyVATRate)),
			),
		),
	)
}

func taxComboIsVAT(val any) bool {
	tc, ok := val.(*tax.Combo)
	return ok && tc != nil && tc.Category == tax.CategoryVAT
}
