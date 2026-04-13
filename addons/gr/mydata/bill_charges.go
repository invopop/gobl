package mydata

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func normalizeBillCharge(c *bill.Charge) {
	if c == nil {
		return
	}

	// Try to map the tax type from the key
	switch c.Key {
	case bill.ChargeKeyStampDuty:
		c.Ext.Set(ExtKeyTaxType, TaxTypeStampDuty)
	case bill.ChargeKeyTax:
		c.Ext.Set(ExtKeyTaxType, TaxTypeOtherTax)
	}
	if c.Ext.Has(ExtKeyTaxType) {
		return
	}

	// Try to map the tax type from the other extensions
	switch {
	case c.Ext.Has(ExtKeyFee):
		c.Ext.Set(ExtKeyTaxType, TaxTypeFee)
	case c.Ext.Has(ExtKeyOtherTax):
		c.Ext.Set(ExtKeyTaxType, TaxTypeOtherTax)
	case c.Ext.Has(ExtKeyStampDuty):
		c.Ext.Set(ExtKeyTaxType, TaxTypeStampDuty)
	}
}

func billChargeRules() *rules.Set {
	return rules.For(new(bill.Charge),
		rules.Field("ext",
			rules.Assert("01", "only one of fee, other-tax, or stamp-duty allowed",
				tax.ExtensionsAllowOneOf(ExtKeyFee, ExtKeyOtherTax, ExtKeyStampDuty),
			),
		),
		rules.When(is.Func("tax type is fee", chargeTaxTypeIsFee),
			rules.Field("ext",
				rules.Assert("02",
					fmt.Sprintf("charge with fee tax type requires '%s' extension", ExtKeyFee),
					tax.ExtensionsRequire(ExtKeyFee),
				),
			),
		),
		rules.When(is.Func("tax type is other tax", chargeTaxTypeIsOtherTax),
			rules.Field("ext",
				rules.Assert("03",
					fmt.Sprintf("charge with other-tax type requires '%s' extension", ExtKeyOtherTax),
					tax.ExtensionsRequire(ExtKeyOtherTax),
				),
			),
		),
		rules.When(is.Func("tax type is stamp duty", chargeTaxTypeIsStampDuty),
			rules.Field("ext",
				rules.Assert("04",
					fmt.Sprintf("charge with stamp-duty type requires '%s' extension", ExtKeyStampDuty),
					tax.ExtensionsRequire(ExtKeyStampDuty),
				),
			),
		),
	)
}

func chargeTaxTypeIsFee(val any) bool {
	c, ok := val.(*bill.Charge)
	return ok && c != nil && c.Ext.Get(ExtKeyTaxType) == TaxTypeFee
}

func chargeTaxTypeIsOtherTax(val any) bool {
	c, ok := val.(*bill.Charge)
	return ok && c != nil && c.Ext.Get(ExtKeyTaxType) == TaxTypeOtherTax
}

func chargeTaxTypeIsStampDuty(val any) bool {
	c, ok := val.(*bill.Charge)
	return ok && c != nil && c.Ext.Get(ExtKeyTaxType) == TaxTypeStampDuty
}
