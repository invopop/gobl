package verifactu

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Simple map of the tax rates to operation classes. These only apply to tax combos
// in Spain and only for the most basic of situations.
var taxCategoryOpClassMap = tax.Extensions{
	tax.RateStandard:     "S1",
	tax.RateReduced:      "S1",
	tax.RateSuperReduced: "S1",
	tax.RateZero:         "S1",
}

func normalizeTaxCombo(tc *tax.Combo) {
	if tc.Country != "" && tc.Country != l10n.ES.Tax() {
		return
	}
	switch tc.Category {
	case tax.CategoryVAT, es.TaxCategoryIGIC:
		if tc.Rate.IsEmpty() {
			return
		}
		v := taxCategoryOpClassMap.Get(tc.Rate)
		if v == "" {
			return
		}
		tc.Ext = tc.Ext.Merge(
			tax.Extensions{ExtKeyOpClass: v},
		)
		// Set default tax regime to "01" (General regime operation) if not specified
		if !tc.Ext.Has(ExtKeyTaxRegime) {
			tc.Ext = tc.Ext.Merge(
				tax.Extensions{ExtKeyTaxRegime: "01"},
			)
		}
	}
}

func validateTaxCombo(tc *tax.Combo) error {
	if !tc.Category.In(tax.CategoryVAT, es.TaxCategoryIGIC) {
		return nil
	}
	return validation.ValidateStruct(tc,
		validation.Field(&tc.Ext,
			validation.When(
				tc.Percent != nil, // Taxed
				tax.ExtensionsRequires(ExtKeyOpClass),
			),
			validation.When(
				tc.Percent == nil && !tc.Ext.Has(ExtKeyOpClass),
				tax.ExtensionsRequires(ExtKeyExempt),
			),
			tax.ExtensionsRequires(ExtKeyTaxRegime),
			validation.Skip,
		),
	)
}
