package mydata

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/gr"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var taxComboRateMapVAT = tax.Extensions{
	tax.RateGeneral:                             "1",
	tax.RateReduced:                             "2",
	tax.RateSuperReduced:                        "3",
	tax.RateGeneral.With(gr.TaxRateIsland):      "4",
	tax.RateReduced.With(gr.TaxRateIsland):      "5",
	tax.RateSuperReduced.With(gr.TaxRateIsland): "6",
}

const taxVATRateExempt cbc.Code = "7"

func normalizeTaxCombo(tc *tax.Combo) {
	if tc == nil {
		return
	}

	// copy the SAF-T tax rate code to the line
	switch tc.Category {
	case tax.CategoryVAT:
		if tc.Country != "" && tc.Country != "EL" {
			// B2C sales outside Greece where local tax is applied.
			c := cbc.Code("29")
			if l10n.Union(l10n.EU).HasMember(tc.Country.Code()) {
				c = cbc.Code("30")
			}
			tc.Ext = tc.Ext.
				Set(ExtKeyExemption, c).
				Set(ExtKeyVATRate, taxVATRateExempt)
		}

		normalizeTaxComboKey(tc)

		switch tc.Key {
		case tax.KeyStandard, tax.KeyZero:
			if k, ok := taxComboRateMapVAT[tc.Rate]; ok {
				tc.Ext = tc.Ext.
					Set(ExtKeyVATRate, k).
					Set(ExtKeyVATRate, taxVATRateExempt)
			}
		case tax.KeyOutsideScope:
			tc.Ext = tc.Ext.
				SetOneOf(ExtKeyExemption, "1", "2", "4", "24", "29", "30", "31").
				Set(ExtKeyVATRate, taxVATRateExempt)
		case tax.KeyExempt:
			tc.Ext = tc.Ext.
				SetOneOf(ExtKeyExemption, "3", "5", "6",
					"7", "11", "12", "13", "15", "17",
					"18", "19", "20", "21", "22", "23", "25",
					"26", "27",
				).
				Set(ExtKeyVATRate, taxVATRateExempt)
		case tax.KeyExport:
			tc.Ext = tc.Ext.
				SetOneOf(ExtKeyExemption, "8", "9", "10", "28").
				Set(ExtKeyVATRate, taxVATRateExempt)
		case tax.KeyIntraCommunity:
			tc.Ext = tc.Ext.
				Set(ExtKeyExemption, "14").
				Set(ExtKeyVATRate, taxVATRateExempt)
		case tax.KeyReverseCharge:
			tc.Ext = tc.Ext.
				Set(ExtKeyExemption, "16").
				Set(ExtKeyVATRate, taxVATRateExempt)
		}
	}
}

func normalizeTaxComboKey(tc *tax.Combo) {
	if !tc.Key.IsEmpty() {
		return
	}
	switch tc.Ext.Get(ExtKeyExemption) {
	case "1", "2", "4", "24", "29", "30", "31":
		tc.Key = tax.KeyOutsideScope
	case "3", "5", "6", "7", "11", "12", "13", "15", "17",
		"18", "19", "20", "21", "22", "23", "25", "26", "27":
		tc.Key = tax.KeyExempt
	case "8", "9", "10", "28":
		tc.Key = tax.KeyExport
	case "14":
		tc.Key = tax.KeyIntraCommunity
	case "16":
		tc.Key = tax.KeyReverseCharge
	default:
		tc.Key = tax.KeyStandard
	}
}

func validateTaxCombo(tc *tax.Combo) error {
	if tc == nil {
		return nil
	}
	switch tc.Category {
	case tax.CategoryVAT:
		return validation.ValidateStruct(tc,
			validation.Field(&tc.Ext,
				tax.ExtensionsRequire(ExtKeyVATRate),
				validation.When(
					tc.Percent == nil,
					tax.ExtensionsRequire(ExtKeyExemption),
				),
				validation.When(
					// MyDATA uses income category and type for accounting purposes
					// and for them to be grouped with taxes. We ensure they're present
					// here so that the
					tc.Ext.Has(ExtKeyIncomeCat) || tc.Ext.Has(ExtKeyIncomeType),
					tax.ExtensionsRequire(ExtKeyIncomeCat, ExtKeyIncomeType),
				),
				validation.Skip,
			),
		)
	}
	return nil
}
