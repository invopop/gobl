package tbai

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
)

func normalizeTaxCombo(tc *tax.Combo) {
	switch tc.Category {
	case tax.CategoryVAT, es.TaxCategoryIGIC:
		if tc.Country != "" && tc.Country != l10n.ES.Tax() {
			// Assume this is a not subject to VAT
			tc.Ext = tc.Ext.
				SetOneOf(ExtKeyExempt, "RL", "IE")
			return
		}

		prepareTaxComboKey(tc)

		// Deterministically set the exemption code.
		switch tc.Key {
		case tax.KeyStandard, tax.KeyZero: // Default
			tc.Ext = tc.Ext.
				Delete(ExtKeyExempt)
		case tax.KeyReverseCharge:
			tc.Ext = tc.Ext.
				Set(ExtKeyExempt, "S2")
		case tax.KeyOutsideScope:
			tc.Ext = tc.Ext.
				SetOneOf(ExtKeyExempt, "OT", "RL", "VT", "IE")
		case tax.KeyExempt:
			tc.Ext = tc.Ext.
				SetOneOf(ExtKeyExempt, "E1", "E6")
		case tax.KeyExport:
			tc.Ext = tc.Ext.
				SetOneOf(ExtKeyExempt, "E2", "E3", "E4")
		case tax.KeyIntraCommunity:
			tc.Ext = tc.Ext.
				Set(ExtKeyExempt, "E5")
		}
	}
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
	case "S1":
		tc.Key = tax.KeyStandard
		tc.Ext = tc.Ext.Delete(ExtKeyExempt)
	case "S2":
		tc.Key = tax.KeyReverseCharge
	case "OT", "RL", "VT", "IE":
		tc.Key = tax.KeyOutsideScope
	}
	if tc.Key.IsEmpty() {
		tc.Key = tax.KeyStandard // "S1" fallback
	}
}
