package saft

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var taxRateMap = tax.Extensions{
	tax.RateReduced:      TaxRateReduced,
	tax.RateIntermediate: TaxRateIntermediate,
	tax.RateGeneral:      TaxRateNormal,
	tax.RateOther:        TaxRateOther, // set when empty
}

func normalizeTaxCombo(tc *tax.Combo) {
	if tc == nil {
		return
	}

	// copy the SAF-T tax rate code to the line
	switch tc.Category {
	case tax.CategoryVAT:
		if tc.Country != "" && tc.Country != l10n.PT.Tax() {
			tc.Ext = tc.Ext.
				Set(ExtKeyTaxRate, TaxRateOther)
			return
		}

		prepareTaxComboKey(tc)

		switch tc.Key {
		case tax.KeyStandard:
			c, ok := taxRateMap[tc.Rate]
			if ok {
				tc.Ext = tc.Ext.
					Delete(ExtKeyExemption).
					Set(ExtKeyTaxRate, c)
			}
		case tax.KeyReverseCharge:
			tc.Ext = tc.Ext.
				Set(ExtKeyTaxRate, TaxRateExempt).
				SetOneOf(ExtKeyExemption, "M40", // assume cross-border is default
					"M30", "M31", "M32", "M33", "M41", "M42", "M43",
				)
		case tax.KeyOutsideScope:
			tc.Ext = tc.Ext.
				Set(ExtKeyTaxRate, TaxRateExempt).
				Set(ExtKeyExemption, "M99")
		case tax.KeyIntraCommunity:
			tc.Ext = tc.Ext.
				Set(ExtKeyTaxRate, TaxRateExempt).
				Set(ExtKeyExemption, "M16")
		case tax.KeyExport:
			tc.Ext = tc.Ext.
				Set(ExtKeyTaxRate, TaxRateExempt).
				SetOneOf(ExtKeyExemption, "M05", "M04")
		case tax.KeyExempt, tax.KeyZero: // no difference in PT
			tc.Ext = tc.Ext.
				Set(ExtKeyTaxRate, TaxRateExempt).
				SetOneOf(ExtKeyExemption, "M07", // health, education, etc.
					"M01", "M02", "M03", "M06", "M09", "M10", "M11",
					"M12", "M13", "M14", "M15", "M19", "M20", "M21",
					"M25", "M26",
				)
		}
	}
}

func prepareTaxComboKey(tc *tax.Combo) {
	// We need to do reverse mappings for the exempt key in order to cope
	// with earlier usage of the "exempt" rate which was too generic.
	if !tc.Key.IsEmpty() && tc.Key != tax.KeyExempt {
		return
	}
	switch tc.Ext.Get(ExtKeyExemption) {
	case "M30", "M31", "M32", "M33", "M40", "M41", "M42", "M43":
		tc.Key = tax.KeyReverseCharge
	case "M05", "M04":
		tc.Key = tax.KeyExport
	case "M16":
		tc.Key = tax.KeyIntraCommunity
	case "M99":
		tc.Key = tax.KeyOutsideScope
	case "M01", "M02", "M03", "M06", "M07", "M09", "M10", "M11",
		"M12", "M13", "M14", "M15", "M19", "M20", "M21", "M25",
		"M26":
		tc.Key = tax.KeyExempt
	default:
		if tc.Key.IsEmpty() {
			tc.Key = tax.KeyStandard
		}
	}
}

func validateTaxCombo(val any) error {
	c, ok := val.(*tax.Combo)
	if !ok {
		return nil
	}
	switch c.Category {
	case tax.CategoryVAT:
		return validation.ValidateStruct(c, validateVATExt(&c.Ext))
	}
	return nil
}

func validateVATExt(ext *tax.Extensions) *validation.FieldRules {
	return validation.Field(ext,
		tax.ExtensionsRequire(pt.ExtKeyRegion, ExtKeyTaxRate),
		validation.When(
			(*ext)[ExtKeyTaxRate] == TaxRateExempt,
			tax.ExtensionsRequire(ExtKeyExemption),
		),
		validation.Skip,
	)
}
