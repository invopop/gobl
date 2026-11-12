package hr

import (
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func taxComboRules() *rules.Set {
	return rules.For(new(tax.Combo),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.Field("ext",
				rules.Assert("01", "Croatian pass-through item code is not valid",
					tax.ExtensionHasValidCode(ExtKeyPassThrough),
				),
			),
			rules.When(
				is.Func("has pass-through item", taxComboHasPassThrough),
				rules.Assert("02", "Croatian pass-through item requires a VAT combo with the exempt key",
					is.Func("VAT exempt", taxComboIsVATExempt),
				),
			),
		),
	)
}

func taxComboHasPassThrough(val any) bool {
	tc, ok := val.(*tax.Combo)
	return ok && tc != nil && tc.Ext.Has(ExtKeyPassThrough)
}

func taxComboIsVATExempt(val any) bool {
	tc, ok := val.(*tax.Combo)
	return ok && tc != nil && tc.Category == tax.CategoryVAT && tc.Key == tax.KeyExempt
}
