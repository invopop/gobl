package arca

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func billChargeRules() *rules.Set {
	return rules.For(new(bill.Charge),
		rules.When(is.Func("key is tax", chargeKeyIsTax),
			rules.Field("ext",
				rules.Assert("01",
					fmt.Sprintf("tax charge requires '%s' extension", ExtKeyTaxType),
					tax.ExtensionsRequire(ExtKeyTaxType),
				),
			),
		),
		rules.When(is.Func("has tax type ext", chargeHasTaxTypeExt),
			rules.Field("percent",
				rules.Assert("02", "percent is required when tax type is set", is.Present),
			),
		),
		rules.When(is.Func("tax type is other", chargeTaxTypeIsOther),
			rules.Field("reason",
				rules.Assert("03", "reason is required when tax type is 'other'", is.Present),
			),
		),
	)
}

func chargeKeyIsTax(val any) bool {
	c, ok := val.(*bill.Charge)
	return ok && c != nil && c.Key.Has(bill.ChargeKeyTax)
}

func chargeHasTaxTypeExt(val any) bool {
	c, ok := val.(*bill.Charge)
	return ok && c != nil && c.Ext.Has(ExtKeyTaxType)
}

func chargeTaxTypeIsOther(val any) bool {
	c, ok := val.(*bill.Charge)
	return ok && c != nil && c.Ext.Get(ExtKeyTaxType) == "99"
}
