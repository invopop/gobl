package mydata

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
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

func validateBillCharge(c *bill.Charge) error {
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
