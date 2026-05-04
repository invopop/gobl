package en16931

import (
	"github.com/invopop/gobl/catalogues/cef"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Official subset of UNTDID 5305 category codes recognized by the EN 16931
const (
	TaxCategoryStandard       cbc.Code = "S"
	TaxCategoryZero           cbc.Code = "Z"
	TaxCategoryExempt         cbc.Code = "E"
	TaxCategoryReverseCharge  cbc.Code = "AE"
	TaxCategoryIntraCommunity cbc.Code = "K"
	TaxCategoryExport         cbc.Code = "G"
	TaxCategoryOutsideScope   cbc.Code = "O"
	TaxCategoryIGIC           cbc.Code = "L" // Canary Islands
	TaxCategoryIPSI           cbc.Code = "M" // Ceuta and Melilla
)

// exemptTaxCategories lists the UNTDID 5305 codes that require either
// a CEF VATEX code or an exemption reason note.
var exemptTaxCategories = []cbc.Code{
	TaxCategoryExempt, TaxCategoryReverseCharge, TaxCategoryIntraCommunity,
	TaxCategoryExport, TaxCategoryOutsideScope,
}

// VAT key mapping from GOBL tax keys to UNTDID 5305 codes.
var vatKeyMap = tax.ExtensionsOf(cbc.CodeMap{
	tax.KeyStandard:       TaxCategoryStandard,
	tax.KeyZero:           TaxCategoryZero,
	tax.KeyExempt:         TaxCategoryExempt,
	tax.KeyReverseCharge:  TaxCategoryReverseCharge,
	tax.KeyIntraCommunity: TaxCategoryIntraCommunity,
	tax.KeyExport:         TaxCategoryExport,
	tax.KeyOutsideScope:   TaxCategoryOutsideScope,
})

func normalizeTaxCombo(tc *tax.Combo) {
	switch tc.Category {
	case tax.CategoryVAT:
		if tc.Key.IsEmpty() {
			// Try doing a reverse map of the VAT category key
			k := vatKeyMap.Lookup(tc.Ext.Get(untdid.ExtKeyTaxCategory))
			if k.IsEmpty() {
				k = tax.KeyStandard
			}
			tc.Key = k
		}
		tc.Ext = tc.Ext.Set(untdid.ExtKeyTaxCategory, vatKeyMap.Get(tc.Key))
	case es.TaxCategoryIGIC:
		tc.Ext = tc.Ext.Set(untdid.ExtKeyTaxCategory, TaxCategoryIGIC)
	case es.TaxCategoryIPSI:
		tc.Ext = tc.Ext.Set(untdid.ExtKeyTaxCategory, TaxCategoryIPSI)
	default:
		// Assume any other tax is outside the scope.
		tc.Ext = tc.Ext.Set(untdid.ExtKeyTaxCategory, TaxCategoryOutsideScope)
	}
}

func taxComboRules() *rules.Set {
	return rules.For(new(tax.Combo),
		rules.Field("ext",
			rules.Assert("01", "tax category extension is required",
				tax.ExtensionsRequire(untdid.ExtKeyTaxCategory),
			),
		),
		rules.When(is.Func("is VAT", taxComboIsVAT),
			rules.Field("ext",
				rules.Assert("02", "VAT category code must be valid",
					tax.ExtensionsHasCodes(untdid.ExtKeyTaxCategory, vatKeyMap.Values()...),
				),
			),
		),
		rules.When(is.Func("is IGIC", taxComboIsIGIC),
			rules.Field("ext",
				rules.Assert("03", "must use IGIC category code",
					tax.ExtensionsHasCodes(untdid.ExtKeyTaxCategory, TaxCategoryIGIC),
				),
			),
		),
		rules.When(is.Func("is IPSI", taxComboIsIPSI),
			rules.Field("ext",
				rules.Assert("04", "must use IPSI category code",
					tax.ExtensionsHasCodes(untdid.ExtKeyTaxCategory, TaxCategoryIPSI),
				),
			),
		),
		rules.When(is.Func("is outside scope", taxComboIsOutsideScope),
			rules.Field("ext",
				rules.Assert("05", "must use outside scope category code",
					tax.ExtensionsHasCodes(untdid.ExtKeyTaxCategory, TaxCategoryOutsideScope),
				),
			),
		),
		rules.When(is.Func("is exempt", taxComboIsExempt),
			rules.Field("ext",
				// BR-E-10: VATEX extension required for exempt tax
				rules.Assert("06", "VATEX extension is required for exempt tax (BR-E-10)",
					tax.ExtensionsRequire(cef.ExtKeyVATEX),
				),
			),
		),
		// BR-S-10, BR-Z-10: standard, zero-rated, IGIC, and IPSI shall NOT have a VATEX code
		rules.When(is.Func("is non-exempt", taxComboIsNonExempt),
			rules.When(is.FuncContext("regime is not saudi arabia", regimeIsNotSA),
				rules.Field("ext",
					rules.Assert("07", "VATEX extension must not be set for standard, zero, IGIC, or IPSI categories (BR-S-10, BR-Z-10)",
						tax.ExtensionsExclude(cef.ExtKeyVATEX),
					),
				),
			),
		),
	)
}

func taxComboIsVAT(val any) bool {
	tc, ok := val.(*tax.Combo)
	return ok && tc != nil && tc.Category == tax.CategoryVAT
}

func taxComboIsIGIC(val any) bool {
	tc, ok := val.(*tax.Combo)
	return ok && tc != nil && tc.Category == es.TaxCategoryIGIC
}

func taxComboIsIPSI(val any) bool {
	tc, ok := val.(*tax.Combo)
	return ok && tc != nil && tc.Category == es.TaxCategoryIPSI
}

func taxComboIsOutsideScope(val any) bool {
	tc, ok := val.(*tax.Combo)
	return ok && tc != nil && !tc.Category.In(tax.CategoryVAT, es.TaxCategoryIGIC, es.TaxCategoryIPSI)
}

func taxComboIsExempt(val any) bool {
	tc, ok := val.(*tax.Combo)
	return ok && tc != nil && tc.Ext.Get(untdid.ExtKeyTaxCategory) == TaxCategoryExempt
}

func taxComboIsNonExempt(val any) bool {
	tc, ok := val.(*tax.Combo)
	return ok && tc != nil && tc.Ext.Get(untdid.ExtKeyTaxCategory).In(TaxCategoryStandard, TaxCategoryZero, TaxCategoryIGIC, TaxCategoryIPSI)
}

func regimeIsNotSA(ctx rules.Context, _ any) bool {
	isSA := false
	ctx.Each(func(v any) bool {
		if tax.RegimeIn(l10n.TaxCountryCode(l10n.SA)).Check(v) {
			isSA = true
			return false
		}
		return true
	})
	return !isSA
}
