package saft

import (
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func normalizeTaxCombo(combo *tax.Combo) {
	// copy the SAF-T tax rate code to the line
	if combo.Category == tax.CategoryVAT {
		var k tax.ExtValue
		switch combo.Rate {
		case tax.RateReduced:
			k = "RED"
		case tax.RateIntermediate:
			k = "INT"
		case tax.RateStandard:
			k = "NOR"
		case tax.RateExempt:
			k = "ISE"
		default:
			k = "OUT"
		}
		if combo.Ext == nil {
			combo.Ext = make(tax.Extensions)
		}
		combo.Ext[ExtKeyTaxRate] = k
	}
}

func validateTaxCombo(val any) error {
	c, ok := val.(*tax.Combo)
	if !ok {
		return nil
	}
	if c.Category != tax.CategoryVAT {
		return nil
	}
	return validation.ValidateStruct(c,
		validation.Field(&c.Ext,
			// NOTE! We know that some tax rate is required in portugal, but
			// we don't know what it should be for foreign countries.
			// Until this is known, we're removing the validation for the
			// country tax rate.
			validation.When(
				c.Country == "",
				tax.ExtensionsRequires(ExtKeyTaxRate),
			),
			validation.When(
				c.Percent == nil,
				tax.ExtensionsRequires(ExtKeyExemption),
			),
			validation.Skip,
		),
	)
}
