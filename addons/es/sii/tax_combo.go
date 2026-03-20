package sii

import (
	"fmt"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
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

func taxComboRules() *rules.Set {
	return rules.For(new(tax.Combo),
		rules.When(
			// Guard: only apply to VAT/IGIC combos processed by SII normalization
			// (which always sets ExtKeyRegime via SetIfEmpty).
			is.Func("sii vat/igic", taxComboForVATorIGIC),
			rules.Field("ext",
				// Code 01: regime is always required for VAT and IGIC
				rules.Assert("01", fmt.Sprintf("extension '%s' is required", ExtKeyRegime),
					tax.ExtensionsRequire(ExtKeyRegime),
				),
				// Code 03: outside scope and exempt are mutually exclusive
				rules.When(
					is.Func("has outside scope", taxComboExtHasOutsideScope),
					rules.Assert("03", fmt.Sprintf("extension '%s' must not be set when '%s' is set", ExtKeyExempt, ExtKeyOutsideScope),
						tax.ExtensionsExclude(ExtKeyExempt),
					),
				),
				// Code 04: E2 and E3 exempt codes not allowed with regime 01
				// https://sede.agenciatributaria.gob.es/static_files/Sede/Procedimiento_ayuda/G417/FicherosSuministros/V_1_1/Validaciones_ErroresSII_v1.1.pdf (Page 51, point 11)
				rules.When(
					tax.ExtensionsHasCodes(ExtKeyRegime, "01"),
					rules.Assert("04", fmt.Sprintf("exempt codes E2 and E3 not allowed with '%s' 01", ExtKeyRegime),
						tax.ExtensionsExcludeCodes(ExtKeyExempt, "E2", "E3"),
					),
				),
			),
			// Code 02: exempt must not be set when percent is set
			rules.When(
				is.Func("has percent", taxComboHasPercent),
				rules.Field("ext",
					rules.Assert("02", fmt.Sprintf("extension '%s' must not be set when percent is set", ExtKeyExempt),
						tax.ExtensionsExclude(ExtKeyExempt),
					),
				),
			),
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

func taxComboForVATorIGIC(val any) bool {
	tc, ok := val.(*tax.Combo)
	return ok && tc != nil && tc.Category.In(tax.CategoryVAT, es.TaxCategoryIGIC)
}

func taxComboHasPercent(val any) bool {
	tc, ok := val.(*tax.Combo)
	return ok && tc != nil && tc.Percent != nil
}

func taxComboExtHasOutsideScope(val any) bool {
	ext, ok := val.(tax.Extensions)
	return ok && ext.Has(ExtKeyOutsideScope)
}
