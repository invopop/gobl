package saft

import (
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// TaxRateExtensions returns the mapping of tax rates defined in PT
// to their extension values used by SAF-T.
//
// Use this to lookup a tax rate key for a SAF-T tax rate code:
//
//	saft.TaxRateExtensions().Lookup("RED") // returns tax.RateReduced
func TaxRateExtensions() tax.Extensions {
	return taxRateMap
}

var taxRateMap = tax.Extensions{
	tax.RateReduced:      "RED",
	tax.RateIntermediate: "INT",
	tax.RateStandard:     "NOR",
	tax.RateExempt:       "ISE",
	tax.RateOther:        "OUT",
}

func normalizeTaxCombo(combo *tax.Combo) {
	// copy the SAF-T tax rate code to the line
	switch combo.Category {
	case tax.CategoryVAT:
		if combo.Rate.IsEmpty() {
			return
		}
		k, ok := taxRateMap[combo.Rate]
		if !ok {
			return
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
	switch c.Category {
	case tax.CategoryVAT:
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
	return nil
}
