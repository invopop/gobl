package sdi

import (
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func normalizeTaxCombo(tc *tax.Combo) {
	if tc == nil {
		return
	}

	// Note: code mappings take from the following URL, Appendix 5.1, with
	// adjustments for using the outside-scope key.
	// https://www.fatturapa.gov.it/export/documenti/Technical-Rules-for-European-Invoicing-v2.5.pdf

	switch tc.Category {
	case tax.CategoryVAT:
		if tc.Country != "" && tc.Country != "IT" {
			if l10n.Union(l10n.EU).HasMember(tc.Country.Code()) {
				tc.Ext = tc.Ext.
					Set(ExtKeyExempt, "N7")
			} else {
				tc.Ext = tc.Ext.
					Set(ExtKeyExempt, "N2.1")
			}
			return
		}
		normalizeTaxComboKey(tc)
		switch tc.Key {
		case tax.KeyStandard:
			tc.Ext = tc.Ext.Delete(ExtKeyExempt)
		case tax.KeyZero:
			tc.Ext = tc.Ext.Set(ExtKeyExempt, "N1")
		case tax.KeyOutsideScope:
			tc.Ext = tc.Ext.SetOneOf(ExtKeyExempt, "N2.1", "N2.2", "N7")
		case tax.KeyExport:
			tc.Ext = tc.Ext.SetOneOf(ExtKeyExempt, "N3.1", "N3.3", "N3.4", "N3.5")
		case tax.KeyIntraCommunity:
			tc.Ext = tc.Ext.SetOneOf(ExtKeyExempt, "N3.2", "N3.6")
		case tax.KeyExempt:
			tc.Ext = tc.Ext.SetOneOf(ExtKeyExempt, "N4", "N5")
		case tax.KeyReverseCharge:
			tc.Ext = tc.Ext.SetOneOf(ExtKeyExempt, "N6.9", "N6.1", "N6.2", "N6.3",
				"N6.4", "N6.5", "N6.6", "N6.7", "N6.8")
		}
	}
}

func normalizeTaxComboKey(tc *tax.Combo) {
	if tc.Key != "" && tc.Key != tax.KeyExempt {
		return
	}
	switch tc.Ext.Get(ExtKeyExempt) {
	case "N1":
		tc.Percent = &num.PercentageZero
		tc.Key = tax.KeyZero
	case "N2.1", "N2.2", "N7":
		tc.Key = tax.KeyOutsideScope
	case "N3.1", "N3.3", "N3.4", "N3.5":
		tc.Key = tax.KeyExport
	case "N3.2", "N3.6":
		tc.Key = tax.KeyIntraCommunity
	case "N4", "N5":
		tc.Key = tax.KeyExempt
	case "N6.1", "N6.2", "N6.3",
		"N6.4", "N6.5", "N6.6",
		"N6.7", "N6.8", "N6.9":
		tc.Key = tax.KeyReverseCharge
	case cbc.CodeEmpty:
		if tc.Key == cbc.KeyEmpty {
			// Assume standard
			tc.Key = tax.KeyStandard
		}
	}
}

func taxComboRules() *rules.Set {
	return rules.For(new(tax.Combo),
		// VAT: exempt extension required when no percent
		rules.When(is.Func("VAT without percent", taxComboIsVATWithoutPercent),
			rules.Field("ext",
				rules.Assert("01",
					fmt.Sprintf("VAT tax combo without percent requires '%s' extension", ExtKeyExempt),
					tax.ExtensionsRequire(ExtKeyExempt),
				),
			),
		),
		// Retained taxes require the retained extension
		rules.When(is.Func("is retained tax", taxComboIsRetained),
			rules.Field("ext",
				rules.Assert("02",
					fmt.Sprintf("retained tax combo requires '%s' extension", ExtKeyRetained),
					tax.ExtensionsRequire(ExtKeyRetained),
				),
			),
		),
	)
}

func taxComboIsVATWithoutPercent(val any) bool {
	c, ok := val.(*tax.Combo)
	if !ok || c == nil {
		return false
	}
	return c.Category == tax.CategoryVAT && c.Percent == nil
}

func taxComboIsRetained(val any) bool {
	c, ok := val.(*tax.Combo)
	if !ok || c == nil {
		return false
	}
	switch c.Category {
	case it.TaxCategoryIRPEF, it.TaxCategoryIRES, it.TaxCategoryINPS, it.TaxCategoryENPAM, it.TaxCategoryENASARCO, it.TaxCategoryCP:
		return true
	}
	return false
}
