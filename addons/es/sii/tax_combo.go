package sii

import (
	"fmt"
	"strings"

	"github.com/invopop/gobl/cbc"
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
				Set(ExtKeyNotSubject, "N2").
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
		case tax.KeyStandard, tax.KeyZero: // Default
			tc.Ext = tc.Ext.
				Set(ExtKeyNotExempt, "S1").
				Delete(ExtKeyExempt)
		case tax.KeyReverseCharge:
			tc.Ext = tc.Ext.
				Set(ExtKeyNotExempt, "S2").
				Delete(ExtKeyExempt)
		case tax.KeyOutsideScope:
			// Default to N2 (not subject due to place of supply rules) since this is most common
			// when providing services to non-EU customers. N1 can be used for other cases where
			// the operation falls outside VAT scope in Spain (e.g. company transfers).
			tc.Ext = tc.Ext.
				SetOneOf(ExtKeyNotSubject, "N2", "N1").
				Delete(ExtKeyExempt)
		case tax.KeyExempt:
			tc.Ext = tc.Ext.
				SetOneOf(ExtKeyExempt, "E1", "E6").
				Delete(ExtKeyNotExempt).
				Delete(ExtKeyNotSubject)
		case tax.KeyExport:
			tc.Ext = tc.Ext.
				SetOneOf(ExtKeyExempt, "E2", "E3", "E4").
				Delete(ExtKeyNotExempt).
				Delete(ExtKeyNotSubject)
		case tax.KeyIntraCommunity:
			tc.Ext = tc.Ext.
				Set(ExtKeyExempt, "E5").
				Delete(ExtKeyNotExempt).
				Delete(ExtKeyNotSubject)
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
			extensionsRequireOneOf(ExtKeyNotSubject, ExtKeyExempt, ExtKeyNotExempt),
			validation.When(
				tc.Percent != nil, // Subject and not exempt
				extensionsRequireOneOf(ExtKeyNotSubject, ExtKeyNotExempt),
			),
			tax.ExtensionsExcludeCodes(ExtKeyNotExempt, "S3"), // S3 can only be set by the converter when both S1 and S2 are present
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
	switch tc.Ext.Get(ExtKeyNotExempt) {
	case "S1":
		tc.Key = tax.KeyStandard
	case "S2":
		tc.Key = tax.KeyReverseCharge
	}
	switch tc.Ext.Get(ExtKeyNotSubject) {
	case "N1", "N2":
		tc.Key = tax.KeyOutsideScope
	}
	if tc.Key.IsEmpty() {
		tc.Key = tax.KeyStandard // "S1" fallback
	}
}

func extensionsRequireOneOf(keys ...cbc.Key) validation.Rule {
	return validation.By(func(val any) error {
		ext, _ := val.(tax.Extensions)
		if len(ext) == 0 {
			return nil
		}
		present := 0
		for _, k := range keys {
			if ext.Has(k) {
				present++
			}
		}
		if present > 1 {
			return fmt.Errorf("only one of %s is allowed", strings.Join(cbc.KeyStrings(keys), ", "))
		}
		if present == 0 {
			return fmt.Errorf("one of %s is required", strings.Join(cbc.KeyStrings(keys), ", "))
		}
		return nil
	})
}
