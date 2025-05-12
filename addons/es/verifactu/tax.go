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
		ext := make(tax.Extensions)

		// Set default tax regime to "01" (General regime operation) if not already specified
		if !tc.Ext.Has(ExtKeyRegime) {
			ext[ExtKeyRegime] = "01"
		}

		if !tc.Rate.IsEmpty() {
			if v := taxCategoryOpClassMap.Get(tc.Rate); v != "" {
				ext[ExtKeyOpClass] = v
			}
		}

		if len(ext) > 0 {
			tc.Ext = tc.Ext.Merge(ext)
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
				tax.ExtensionsRequire(ExtKeyOpClass),
			),
			validation.When(
				tc.Percent == nil && !tc.Ext.Has(ExtKeyOpClass),
				tax.ExtensionsRequire(ExtKeyExempt),
			),
			tax.ExtensionsRequire(ExtKeyRegime),
			// https://www.agenciatributaria.es/static_files/AEAT_Desarrolladores/EEDD/IVA/VERI-FACTU/Validaciones_Errores_Veri-Factu.pdf (Page 10, section 15.5)
			validation.When(
				(tc.Category == tax.CategoryVAT || tc.Category == es.TaxCategoryIGIC) && tc.Ext.Get(ExtKeyRegime) == "01",
				tax.ExtensionsExcludeCodes(ExtKeyExempt, "E2", "E3"),
			),
			validation.Skip,
		),
	)
}
