package mydata

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/gr"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var taxComboMapVAT = map[cbc.Key]cbc.Code{
	tax.RateStandard:                            "1",
	tax.RateReduced:                             "2",
	tax.RateSuperReduced:                        "3",
	tax.RateStandard.With(gr.TaxRateIsland):     "4",
	tax.RateReduced.With(gr.TaxRateIsland):      "5",
	tax.RateSuperReduced.With(gr.TaxRateIsland): "6",
	tax.RateExempt:                              "7",
}

func normalizeTaxCombo(tc *tax.Combo) {
	// copy the SAF-T tax rate code to the line
	switch tc.Category {
	case tax.CategoryVAT:
		if k, ok := taxComboMapVAT[tc.Rate]; ok {
			if tc.Ext == nil {
				tc.Ext = make(tax.Extensions)
			}
			tc.Ext[ExtKeyVATRate] = k
		}
	}
}

func validateTaxCombo(tc *tax.Combo) error {
	if tc == nil {
		return nil
	}
	switch tc.Category {
	case tax.CategoryVAT:
		return validation.ValidateStruct(tc,
			validation.Field(&tc.Ext,
				tax.ExtensionsRequire(ExtKeyVATRate),
				validation.When(
					tc.Percent == nil,
					tax.ExtensionsRequire(ExtKeyExemption),
				),
				validation.When(
					// MyDATA uses income category and type for accounting purposes
					// and for them to be grouped with taxes. We ensure they're present
					// here so that the
					tc.Ext.Has(ExtKeyIncomeCat) || tc.Ext.Has(ExtKeyIncomeType),
					tax.ExtensionsRequire(ExtKeyIncomeCat, ExtKeyIncomeType),
				),
				validation.Skip,
			),
		)
	}
	return nil
}
