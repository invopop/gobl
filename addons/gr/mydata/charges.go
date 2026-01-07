package mydata

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func normalizeCharge(c *bill.Charge) {
	if c == nil {
		return
	}
	if c.Ext.Has(ExtKeyTaxType) {
		return
	}
	switch {
	case c.Ext.Has(ExtKeyFee):
		c.Ext.Set(ExtKeyTaxType, TaxTypeFee)
	case c.Ext.Has(ExtKeyOtherTax):
		c.Ext.Set(ExtKeyTaxType, TaxTypeOtherTax)
	case c.Ext.Has(ExtKeyStampDuty):
		c.Ext.Set(ExtKeyTaxType, TaxTypeStampDuty)
	}
}

func validateCharge(c *bill.Charge) error {
	if c == nil {
		return nil
	}
	return validation.ValidateStruct(c,
		validation.Field(&c.Ext,
			tax.ExtensionsAllowOneOf(ExtKeyFee, ExtKeyOtherTax, ExtKeyStampDuty),
			validation.When(
				c.Ext.Get(ExtKeyTaxType) == TaxTypeFee,
				tax.ExtensionsRequire(ExtKeyFee),
			),
			validation.When(
				c.Ext.Get(ExtKeyTaxType) == TaxTypeOtherTax,
				tax.ExtensionsRequire(ExtKeyOtherTax),
			),
			validation.When(
				c.Ext.Get(ExtKeyTaxType) == TaxTypeStampDuty,
				tax.ExtensionsRequire(ExtKeyStampDuty),
			),
			validation.Skip,
		),
	)
}
