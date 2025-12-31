package mydata

import (
	"slices"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var taxCategoryExtensions = []cbc.Key{
	ExtKeyFee,
	ExtKeyOtherTax,
	ExtKeyStampDuty,
}

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
			validation.When(
				c.Ext[ExtKeyTaxType] == TaxTypeFee,
				tax.ExtensionsRequire(ExtKeyFee),
				tax.ExtensionsExclude(taxCategoryExtensionsExcept(ExtKeyFee)...),
			),
			validation.When(
				c.Ext[ExtKeyTaxType] == TaxTypeOtherTax,
				tax.ExtensionsRequire(ExtKeyOtherTax),
				tax.ExtensionsExclude(taxCategoryExtensionsExcept(ExtKeyOtherTax)...),
			),
			validation.When(
				c.Ext[ExtKeyTaxType] == TaxTypeStampDuty,
				tax.ExtensionsRequire(ExtKeyStampDuty),
				tax.ExtensionsExclude(taxCategoryExtensionsExcept(ExtKeyStampDuty)...),
			),
			validation.Skip,
		),
	)
}

func taxCategoryExtensionsExcept(key cbc.Key) []cbc.Key {
	return slices.DeleteFunc(slices.Clone(taxCategoryExtensions), func(k cbc.Key) bool {
		return k == key
	})
}
