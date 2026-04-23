package mydata

import (
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/gr"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

var taxComboRateMapVAT = map[cbc.Key]cbc.Code{
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
			tc.Key = tax.KeyOutsideScope
			return
		}

		normalizeTaxComboKey(tc)

		switch tc.Key {
		case tax.KeyStandard, tax.KeyZero:
			if k, ok := taxComboRateMapVAT[tc.Rate]; ok {
				tc.Ext = tc.Ext.
					Set(ExtKeyVATRate, k)
			}
			tc.Ext = tc.Ext.Delete(ExtKeyExemption)
		case tax.KeyOutsideScope:
			tc.Ext = tc.Ext.
				SetOneOf(ExtKeyExemption, "1", "2", "3", "4", "24", "29", "30", "31").
				Set(ExtKeyVATRate, taxVATRateExempt)
		case tax.KeyExempt:
			tc.Ext = tc.Ext.
				SetOneOf(ExtKeyExemption, "7", "5", "6",
					"9", "10", "11", "12", "13", "15", "17",
					"18", "19", "20", "21", "22", "23", "25",
					"26", "27",
				).
				Set(ExtKeyVATRate, taxVATRateExempt)
		case tax.KeyExport:
			tc.Ext = tc.Ext.
				SetOneOf(ExtKeyExemption, "8", "28").
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
	case "1", "2", "3", "4", "24", "29", "30", "31":
		tc.Key = tax.KeyOutsideScope
	case "5", "6", "7", "9", "10", "11", "12", "13", "15", "17",
		"18", "19", "20", "21", "22", "23", "25", "26", "27":
		tc.Key = tax.KeyExempt
	case "8", "28":
		tc.Key = tax.KeyExport
	case "14":
		tc.Key = tax.KeyIntraCommunity
	case "16":
		tc.Key = tax.KeyReverseCharge
	default:
		tc.Key = tax.KeyStandard
	}
}

func taxComboRules() *rules.Set {
	return rules.For(new(tax.Combo),
		rules.When(is.Func("category is VAT", taxComboIsVAT),
			rules.Field("ext",
				rules.Assert("01",
					fmt.Sprintf("VAT combo requires '%s' extension", ExtKeyVATRate),
					tax.ExtensionsRequire(ExtKeyVATRate),
				),
			),
		),
		rules.When(is.Func("VAT with no percent", taxComboVATNoPercent),
			rules.Field("ext",
				rules.Assert("02",
					fmt.Sprintf("exempt VAT combo requires '%s' extension", ExtKeyExemption),
					tax.ExtensionsRequire(ExtKeyExemption),
				),
			),
		),
		rules.When(is.Func("VAT with income ext", taxComboVATHasIncomeExt),
			rules.Field("ext",
				rules.Assert("03",
					fmt.Sprintf("income extensions '%s' and '%s' must both be present",
						ExtKeyIncomeCat, ExtKeyIncomeType),
					tax.ExtensionsRequire(ExtKeyIncomeCat, ExtKeyIncomeType),
				),
			),
		),
	)
}

func taxComboIsVAT(val any) bool {
	tc, ok := val.(*tax.Combo)
	return ok && tc != nil && tc.Category == tax.CategoryVAT
}

func taxComboVATNoPercent(val any) bool {
	tc, ok := val.(*tax.Combo)
	return ok && tc != nil && tc.Category == tax.CategoryVAT && tc.Percent == nil
}

func taxComboVATHasIncomeExt(val any) bool {
	tc, ok := val.(*tax.Combo)
	if !ok || tc == nil || tc.Category != tax.CategoryVAT {
		return false
	}
	return tc.Ext.Has(ExtKeyIncomeCat) || tc.Ext.Has(ExtKeyIncomeType)
}
