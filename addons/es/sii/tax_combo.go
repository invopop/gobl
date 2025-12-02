package sii

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
				Set(ExtKeyOutsideScope, "location").
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
		tc.Ext = tc.Ext.SetIfEmpty(ExtKeyRegime, "01")

		// Deterministically set the operation class and exemption code.
		switch tc.Key {
		case tax.KeyStandard, tax.KeyZero, tax.KeyReverseCharge:
			tc.Ext = tc.Ext.Delete(ExtKeyExempt)
		case tax.KeyOutsideScope:
			// Default to `location` (not subject due to place of supply rules) since this is most common
			// when providing services to non-EU customers. `other` can be used for other cases where the
			// operation falls outside VAT scope in Spain (e.g. company transfers).
			tc.Ext = tc.Ext.
				SetOneOf(ExtKeyOutsideScope, "location", "other").
				Delete(ExtKeyExempt)
		case tax.KeyExempt:
			tc.Ext = tc.Ext.
				SetOneOf(ExtKeyExempt, "E1", "E6").
				Delete(ExtKeyOutsideScope)
		case tax.KeyExport:
			tc.Ext = tc.Ext.
				SetOneOf(ExtKeyExempt, "E2", "E3", "E4").
				Delete(ExtKeyOutsideScope)
		case tax.KeyIntraCommunity:
			tc.Ext = tc.Ext.
				Set(ExtKeyExempt, "E5").
				Delete(ExtKeyOutsideScope)
		}
	}
}

func validateTaxCombo(tc *tax.Combo) error {
	if !tc.Category.In(tax.CategoryVAT, es.TaxCategoryIGIC) {
		return nil
	}
	return validation.ValidateStruct(tc,
		validation.Field(&tc.Ext,
			validation.Required,
			// Regime is always required for VAT and IGIC
			tax.ExtensionsRequire(ExtKeyRegime),
			validation.When(
				tc.Percent != nil,
				tax.ExtensionsExclude(ExtKeyExempt),
			),
			validation.When(
				// ExtKeyOutsideScope and ExtKeyExempt are mutually exclusive
				tc.Ext.Has(ExtKeyOutsideScope),
				tax.ExtensionsExclude(ExtKeyExempt),
			),
			// https://sede.agenciatributaria.gob.es/static_files/Sede/Procedimiento_ayuda/G417/FicherosSuministros/V_1_1/Validaciones_ErroresSII_v1.1.pdf (Page 51, point 11)
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
	if tc.Ext.Has(ExtKeyOutsideScope) {
		tc.Key = tax.KeyOutsideScope
	}
	if tc.Key.IsEmpty() {
		tc.Key = tax.KeyStandard
	}
}
