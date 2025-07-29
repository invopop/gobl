package verifactu

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Simple map of the tax rates to operation classes.
var taxCategoryOpClassMap = tax.Extensions{
	tax.RateStandard:     "S1",
	tax.RateReduced:      "S1",
	tax.RateSuperReduced: "S1",
	tax.RateZero:         "S1",
	tax.TagReverseCharge: "S2",
}

var taxCategoryExemptMap = tax.Extensions{
	tax.RateExempt:                                      "E1",
	tax.RateExempt.With(tax.TagExport):                  "E2",
	tax.RateExempt.With(tax.TagExport).With(tax.TagEEA): "E5",
}

func normalizeTaxCombo(tc *tax.Combo) {
	if tc.Country != "" && tc.Country != l10n.ES.Tax() {
		return
	}

	switch tc.Category {
	case tax.CategoryVAT, es.TaxCategoryIGIC:
		ext := make(tax.Extensions)

		// Try to automatically determine the regime if not already set.
		if !tc.Ext.Has(ExtKeyRegime) {
			// Set default tax regime to "01" (General regime operation)
			ext[ExtKeyRegime] = "01"
			// If the operation is an export, set the regime to "02" (Export operation)
			if tc.Rate.Has(tax.TagExport) {
				ext[ExtKeyRegime] = "02"
			}
			// If the operation is a surcharge and the category is VAT, set the regime to "18" (Surcharge operation)
			if tc.Surcharge != nil && tc.Category == tax.CategoryVAT {
				ext[ExtKeyRegime] = "18"
			}
		}

		if !tc.Rate.IsEmpty() {
			// Assign operation class if applicable
			if v := taxCategoryOpClassMap.Get(tc.Rate); v != "" {
				ext[ExtKeyOpClass] = v
			}

			// Assign exempt code only if no operation class is present
			if !ext.Has(ExtKeyOpClass) {
				if v := taxCategoryExemptMap.Get(tc.Rate); v != "" {
					ext[ExtKeyExempt] = v
				}
			}
		}

		if len(ext) > 0 {
			tc.Ext = tc.Ext.Merge(ext)
		}

		if tc.Ext.Has(ExtKeyOpClass) {
			// cannot have exempt reason alongside operation class
			delete(tc.Ext, ExtKeyExempt)
		}
	}
}

func validateTaxCombo(tc *tax.Combo) error {
	if !tc.Category.In(tax.CategoryVAT, es.TaxCategoryIGIC) {
		return nil
	}
	return validation.ValidateStruct(tc,
		validation.Field(&tc.Ext,
			// Regime is always required for VAT and IGIC
			tax.ExtensionsRequire(ExtKeyRegime),
			validation.When(
				tc.Percent != nil, // Taxed
				tax.ExtensionsRequire(ExtKeyOpClass),
				validation.Required,
			),
			validation.When(
				// Cannot use both exempt and operation class at same time.
				tc.Ext.Has(ExtKeyOpClass),
				tax.ExtensionsExclude(ExtKeyExempt),
			),
			// https://www.agenciatributaria.es/static_files/AEAT_Desarrolladores/EEDD/IVA/VERI-FACTU/Validaciones_Errores_Veri-Factu.pdf (Page 10, section 15.5)
			validation.When(
				tc.Ext.Get(ExtKeyRegime).In("01"),
				tax.ExtensionsExcludeCodes(ExtKeyExempt, "E2", "E3"),
			),
			validation.Skip,
		),
	)
}
