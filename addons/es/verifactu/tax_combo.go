package verifactu

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func normalizeTaxCombo(tc *tax.Combo) {
	switch tc.Category {
	case tax.CategoryVAT, es.TaxCategoryIGIC:
		if tc.Country != "" && tc.Country != l10n.ES.Tax() {
			// Assume this is a not subject to VAT
			tc.Ext = tc.Ext.
				Set(ExtKeyOpClass, "N2").
				SetOneOf(ExtKeyRegime, "01", "17").
				Delete(ExtKeyExempt)
			return
		}

		prepareTaxComboKey(tc)

		// Try to automatically determine the regime if not already set.
		// This approach is not deterministic.
		if tc.Key == tax.KeyExport {
			tc.Ext = tc.Ext.SetIfEmpty(ExtKeyRegime, "02")
		}
		if tc.Surcharge != nil {
			tc.Ext = tc.Ext.SetIfEmpty(ExtKeyRegime, "18")
		}
		tc.Ext = tc.Ext.SetIfEmpty(ExtKeyRegime, "01")

		// Deterministically set the operation class and exemption code.
		switch tc.Key {
		case tax.KeyStandard, tax.KeyZero: // Default
			tc.Ext = tc.Ext.
				Set(ExtKeyOpClass, "S1").
				Delete(ExtKeyExempt)
		case tax.KeyReverseCharge:
			tc.Ext = tc.Ext.
				Set(ExtKeyOpClass, "S2").
				Delete(ExtKeyExempt)
		case tax.KeyOutsideScope:
			// Default to N2 (not subject due to place of supply rules) since this is most common
			// when providing services to non-EU customers. N1 can be used for other cases where
			// the operation falls outside VAT scope in Spain (e.g. company transfers).
			tc.Ext = tc.Ext.
				SetOneOf(ExtKeyOpClass, "N2", "N1").
				Delete(ExtKeyExempt)
		case tax.KeyExempt:
			tc.Ext = tc.Ext.
				SetOneOf(ExtKeyExempt, "E1", "E6").
				Delete(ExtKeyOpClass)
		case tax.KeyExport:
			tc.Ext = tc.Ext.
				SetOneOf(ExtKeyExempt, "E2", "E3", "E4").
				Delete(ExtKeyOpClass)
		case tax.KeyIntraCommunity:
			tc.Ext = tc.Ext.
				Set(ExtKeyExempt, "E5").
				Delete(ExtKeyOpClass)
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

// prepareTaxComboKey tries to reverse map existing extension keys into the
// appropriate tax combo key. This helps with the migration period when getting
// users to move to keys.
func prepareTaxComboKey(tc *tax.Combo) {
	if !tc.Key.IsEmpty() {
		return
	}
	switch tc.Ext.Get(ExtKeyExempt) {
	case "E1", "E6":
		tc.Key = tax.KeyExempt
	case "E2", "E3", "E4":
		tc.Key = tax.KeyExport
	case "E5":
		tc.Key = tax.KeyIntraCommunity
	}
	switch tc.Ext.Get(ExtKeyOpClass) {
	case "S2":
		tc.Key = tax.KeyReverseCharge
	case "N1", "N2":
		tc.Key = tax.KeyOutsideScope
	}
	if tc.Key.IsEmpty() {
		tc.Key = tax.KeyStandard // "S1" fallback
	}
}
