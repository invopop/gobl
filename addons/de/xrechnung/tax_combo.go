package xrechnung

import (
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func TaxRateExtensions() tax.Extensions {
	return taxRateMap
}

var taxRateMap = tax.Extensions{
	tax.RateStandard: "S",
	tax.RateZero:     "Z",
	tax.RateExempt:   "E",
}

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
