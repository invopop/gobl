package sdi

import (
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateTaxCombo(val any) error {
	c, ok := val.(*tax.Combo)
	if !ok {
		return nil
	}
	switch c.Category {
	case tax.CategoryVAT:
		return validation.ValidateStruct(c,
			validation.Field(&c.Ext,
				validation.When(
					c.Percent == nil,
					tax.ExtensionsRequires(ExtKeyExempt),
				),
				validation.Skip,
			),
		)
	// ensure retained taxes have the required extension
	case it.TaxCategoryIRPEF, it.TaxCategoryIRES, it.TaxCategoryINPS, it.TaxCategoryENPAM, it.TaxCategoryENASARCO:
		return validation.ValidateStruct(c,
			validation.Field(&c.Ext,
				tax.ExtensionsRequires(ExtKeyRetained),
				validation.Skip,
			),
		)
	}
	return nil
}
