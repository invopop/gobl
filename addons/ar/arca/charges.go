package arca

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func billChargeRules() *rules.Set {
	return rules.For(new(bill.Charge),
		// Tax charges require the tax type extension
		rules.When(is.Func("has tax key", chargeHasTaxKey),
			rules.Field("ext",
				rules.Assert("01", "ar-arca-tax-type: required", tax.ExtensionsRequire(ExtKeyTaxType)),
			),
		),
		// When tax type extension is set, percent is required
		rules.When(is.Func("has tax type ext", chargeHasTaxTypeExt),
			rules.Field("percent",
				rules.Assert("02", "cannot be blank", is.Present),
			),
		),
		// When tax type is "other" (99), reason is required
		rules.When(is.Func("tax type other", chargeHasTaxTypeOther),
			rules.Field("reason",
				rules.Assert("03", "reason is required when tax type is 'other'", is.Present),
			),
		),
	)
}

func chargeHasTaxKey(val any) bool {
	c, ok := val.(*bill.Charge)
	return ok && c != nil && c.Key.Has(bill.ChargeKeyTax)
}

func chargeHasTaxTypeExt(val any) bool {
	c, ok := val.(*bill.Charge)
	return ok && c != nil && c.Ext.Has(ExtKeyTaxType)
}

func chargeHasTaxTypeOther(val any) bool {
	c, ok := val.(*bill.Charge)
	return ok && c != nil && c.Ext.Get(ExtKeyTaxType) == "99"
}
