package sdi

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func normalizeTaxCombo(tc *tax.Combo) {
	if tc == nil {
		return
	}
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
		case tax.KeyStandard, tax.KeyZero:
			tc.Ext = tc.Ext.Delete(ExtKeyExempt)
		case tax.KeyOutsideScope:
			tc.Ext = tc.Ext.SetOneOf(ExtKeyExempt, "N1", "N2.1", "N2.2", "N7")
		case tax.KeyExport:
			tc.Ext = tc.Ext.SetOneOf(ExtKeyExempt, "N3.1", "N3.3", "N3.4")
		case tax.KeyIntraCommunity:
			tc.Ext = tc.Ext.Set(ExtKeyExempt, "N3.2")
		case tax.KeyExempt:
			tc.Ext = tc.Ext.SetOneOf(ExtKeyExempt, "N4", "N3.5", "N3.6", "N5")
		case tax.KeyReverseCharge:
			tc.Ext = tc.Ext.SetOneOf(ExtKeyExempt, "N6.9", "N6.1", "N6.2", "N6.3",
				"N6.4", "N6.5", "N6.6", "N6.7", "N6.8")
		}
	}
}

func normalizeTaxComboKey(tc *tax.Combo) {
	if tc.Key != "" {
		return
	}
	switch tc.Ext.Get(ExtKeyExempt) {
	case "N1", "N2.1", "N2.2", "N7":
		tc.Key = tax.KeyOutsideScope
	case "N3.1", "N3.3", "N3.4":
		tc.Key = tax.KeyExport
	case "N4", "N3.5", "N3.6", "N5":
		tc.Key = tax.KeyExempt
	case "N6.1", "N6.2", "N6.3",
		"N6.4", "N6.5", "N6.6",
		"N6.7", "N6.8", "N6.9":
		tc.Key = tax.KeyReverseCharge
	case "N3.2":
		tc.Key = tax.KeyIntraCommunity
	case cbc.CodeEmpty:
		// Assume standard, zero rate will have been normalized already
		tc.Key = tax.KeyStandard
	}
}

func validateTaxCombo(val any) error {
	c, ok := val.(*tax.Combo)
	if !ok {
		return nil
	}
	switch c.Category {
	case tax.CategoryVAT:
		return validation.ValidateStruct(c,
			validation.Field(&c.Ext,
				validation.When(
					c.Percent == nil,
					tax.ExtensionsRequire(ExtKeyExempt),
				),
				validation.Skip,
			),
		)
	// ensure retained taxes have the required extension
	case it.TaxCategoryIRPEF, it.TaxCategoryIRES, it.TaxCategoryINPS, it.TaxCategoryENPAM, it.TaxCategoryENASARCO:
		return validation.ValidateStruct(c,
			validation.Field(&c.Ext,
				tax.ExtensionsRequire(ExtKeyRetained),
				validation.Skip,
			),
		)
	}
	return nil
}
