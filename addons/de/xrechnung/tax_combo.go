package xrechnung

import (
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// TaxRateExtensions returns the mapping of tax rates defined in DE
// to their extension values used by XRechnung.
func TaxRateExtensions() tax.Extensions {
	return taxRateMap
}

var taxRateMap = tax.Extensions{
	tax.RateStandard: "S",
	tax.RateZero:     "Z",
	tax.RateExempt:   "E",
}

// NormalizeTaxCombo adds the XRechnung tax rate code to the tax combo.
func NormalizeTaxCombo(combo *tax.Combo) {
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

// ValidateTaxCombo validates percentage is included as BR-DE-14 indicates
func ValidateTaxCombo(tc *tax.Combo) error {
	if tc == nil {
		return nil
	}
	// BR-DE-14: Percentage required for VAT
	return validation.ValidateStruct(tc,
		validation.Field(&tc.Percent,
			validation.When(tc.Category == tax.CategoryVAT,
				validation.Required),
		),
	)
}
