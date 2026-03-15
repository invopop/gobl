package verifactu

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
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

func taxComboRules() *rules.Set {
	return rules.For(new(tax.Combo),
		rules.When(
			// Guard: only apply to VAT/IGIC combos that have been processed by verifactu
			// normalization (which always sets ExtKeyRegime via SetIfEmpty).
			rules.By("verifactu vat/igic", comboIsVerifactuVATorIGIC),
			rules.Field("ext",
				rules.Assert("01", "regime is required", rules.By("has regime", extHasRegime)),
			),
			rules.Assert("10", "cannot use both op_class and exempt at the same time",
				rules.By("not both op_class and exempt", comboNotBothOpClassAndExempt),
			),
			rules.When(
				rules.By("has percent", comboHasPercent),
				rules.Field("ext",
					rules.Assert("02", "op_class is required for taxed operations",
						rules.By("has op_class", extHasOpClass),
					),
				),
			),
			// https://www.agenciatributaria.es/static_files/AEAT_Desarrolladores/EEDD/IVA/VERI-FACTU/Validaciones_Errores_Veri-Factu.pdf (Page 10, section 15.5)
			rules.When(
				rules.By("regime 01", comboRegimeIs01),
				rules.Field("ext",
					rules.Assert("11", "exempt codes E2 and E3 not allowed with regime 01",
						rules.By("no E2/E3", extExcludesE2E3),
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

func comboIsVerifactuVATorIGIC(val any) bool {
	tc, ok := val.(*tax.Combo)
	return ok && tc != nil && tc.Category.In(tax.CategoryVAT, es.TaxCategoryIGIC) && tc.Ext.Has(ExtKeyRegime)
}

func extHasRegime(val any) bool {
	ext, ok := val.(tax.Extensions)
	return ok && ext.Has(ExtKeyRegime)
}

func comboHasPercent(val any) bool {
	tc, ok := val.(*tax.Combo)
	return ok && tc != nil && tc.Percent != nil
}

func extHasOpClass(val any) bool {
	ext, ok := val.(tax.Extensions)
	return ok && ext.Has(ExtKeyOpClass)
}

func comboNotBothOpClassAndExempt(val any) bool {
	tc, ok := val.(*tax.Combo)
	if !ok || tc == nil {
		return true
	}
	return !(tc.Ext.Has(ExtKeyOpClass) && tc.Ext.Has(ExtKeyExempt))
}

func comboRegimeIs01(val any) bool {
	tc, ok := val.(*tax.Combo)
	return ok && tc != nil && tc.Ext.Get(ExtKeyRegime).In("01")
}

func extExcludesE2E3(val any) bool {
	ext, ok := val.(tax.Extensions)
	return ok && !ext.Get(ExtKeyExempt).In("E2", "E3")
}
