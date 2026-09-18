package verifactu

import (
	"fmt"

	"github.com/invopop/gobl/cbc"
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

		normalizeTaxComboClassification(tc)

	case es.TaxCategoryIPSI:
		// IPSI (Ceuta and Melilla) shares the operation class and exemption
		// code lists with VAT and IGIC. Since revision 1.1.6 of the AEAT
		// validation rules (November 2025) a regime code is also expected,
		// drawn from the IPSI-specific subset in ipsiRegimeCodes.
		prepareTaxComboKey(tc)
		normalizeTaxComboClassification(tc)
		if tc.Ext.Get(ExtKeyExempt) == "E1" {
			// Domestic exemptions under article 7 of Ley 8/1991 map to
			// "19 - Operaciones interiores exentas" for IPSI.
			tc.Ext = tc.Ext.SetIfEmpty(ExtKeyRegime, "19")
		}
		tc.Ext = tc.Ext.SetIfEmpty(ExtKeyRegime, "01")
	}
}

// ipsiRegimeCodes lists the ClaveRegimen values the AEAT accepts when
// Impuesto is 02 (IPSI). See section 15.6 of the VERI*FACTU validation rules.
var ipsiRegimeCodes = []cbc.Code{"01", "08", "11", "18", "19", "20"}

// normalizeTaxComboClassification deterministically sets the operation class
// or the exemption code according to the tax combo key, ensuring the two
// extensions are never present at the same time.
func normalizeTaxComboClassification(tc *tax.Combo) {
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

func taxComboRules() *rules.Set {
	return rules.For(new(tax.Combo),
		rules.When(
			// Guard: the E2/E3 restriction on regime 01 only applies to VAT and
			// IGIC (AEAT validation 15.5).
			is.Func("verifactu vat/igic", taxComboForVATorIGIC),
			rules.Field("ext",
				rules.When(
					tax.ExtensionsHasCodes(ExtKeyRegime, "01"),
					rules.Assert("02", fmt.Sprintf("exempt codes E2 and E3 not allowed with '%s' 01", ExtKeyRegime),
						tax.ExtensionsExcludeCodes(ExtKeyExempt, "E2", "E3"),
					),
				),
			),
		),
		rules.When(
			// Guard: IPSI only accepts a subset of the regime codes (AEAT validation 15.6).
			is.Func("verifactu ipsi", taxComboForIPSI),
			rules.Field("ext",
				rules.Assert("05", fmt.Sprintf("extension '%s' for IPSI must be one of %v", ExtKeyRegime, ipsiRegimeCodes),
					tax.ExtensionsHasCodes(ExtKeyRegime, ipsiRegimeCodes...),
				),
			),
		),
		rules.When(
			// Guard: regime, operation class and exemption codes apply to VAT, IGIC,
			// and IPSI. Verifactu normalization always sets the regime for these
			// categories via SetIfEmpty.
			is.Func("verifactu vat/igic/ipsi", taxComboForVATorIGICorIPSI),
			rules.Field("ext",
				rules.Assert("01", fmt.Sprintf("extension '%s' is required", ExtKeyRegime),
					tax.ExtensionsRequire(ExtKeyRegime),
				),
				rules.Assert("03", fmt.Sprintf("cannot use both '%s' and '%s' at the same time", ExtKeyOpClass, ExtKeyExempt),
					tax.ExtensionsAllowOneOf(ExtKeyOpClass, ExtKeyExempt),
				),
			),
			rules.When(
				is.Func("has percent", taxComboHasPercent),
				rules.Field("ext",
					rules.Assert("04", fmt.Sprintf("extension '%s' is required for taxed operations", ExtKeyOpClass),
						tax.ExtensionsRequire(ExtKeyOpClass),
					),
				),
			),
			// https://www.agenciatributaria.es/static_files/AEAT_Desarrolladores/EEDD/IVA/VERI-FACTU/Validaciones_Errores_Veri-Factu.pdf (Page 10, section 15.5)
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

func taxComboForVATorIGIC(val any) bool {
	tc, ok := val.(*tax.Combo)
	return ok && tc != nil && tc.Category.In(tax.CategoryVAT, es.TaxCategoryIGIC)
}

func taxComboForIPSI(val any) bool {
	tc, ok := val.(*tax.Combo)
	return ok && tc != nil && tc.Category == es.TaxCategoryIPSI
}

func taxComboForVATorIGICorIPSI(val any) bool {
	tc, ok := val.(*tax.Combo)
	return ok && tc != nil && tc.Category.In(tax.CategoryVAT, es.TaxCategoryIGIC, es.TaxCategoryIPSI)
}

func taxComboHasPercent(val any) bool {
	tc, ok := val.(*tax.Combo)
	return ok && tc != nil && tc.Percent != nil
}
