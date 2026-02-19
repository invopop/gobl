package en16931

import (
	"github.com/invopop/gobl/catalogues/cef"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
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

// VAT key mapping from GOBL tax keys to UNTDID 5305 codes.
var vatKeyMap = tax.Extensions{
	tax.KeyStandard:       TaxCategoryStandard,
	tax.KeyZero:           TaxCategoryZero,
	tax.KeyExempt:         TaxCategoryExempt,
	tax.KeyReverseCharge:  TaxCategoryReverseCharge,
	tax.KeyIntraCommunity: TaxCategoryIntraCommunity,
	tax.KeyExport:         TaxCategoryExport,
	tax.KeyOutsideScope:   TaxCategoryOutsideScope,
}

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

func validateTaxCombo(tc *tax.Combo) error {
	if tc == nil {
		return nil
	}
	return validation.ValidateStruct(tc,
		validation.Field(&tc.Ext,
			tax.ExtensionsRequire(untdid.ExtKeyTaxCategory),
			validation.When(
				tc.Category == tax.CategoryVAT,
				tax.ExtensionsHasCodes(untdid.ExtKeyTaxCategory, vatKeyMap.Values()...),
			),
			validation.When(
				tc.Category == es.TaxCategoryIGIC,
				tax.ExtensionsHasCodes(untdid.ExtKeyTaxCategory, TaxCategoryIGIC),
			),
			validation.When(
				tc.Category == es.TaxCategoryIPSI,
				tax.ExtensionsHasCodes(untdid.ExtKeyTaxCategory, TaxCategoryIPSI),
			),
			validation.When(
				!tc.Category.In(tax.CategoryVAT, es.TaxCategoryIGIC, es.TaxCategoryIPSI),
				tax.ExtensionsHasCodes(untdid.ExtKeyTaxCategory, TaxCategoryOutsideScope),
			),
			validation.When(
				// BR-S-10, BR-Z-10: standard and zero-rated shall NOT have an exemption reason
				tc.Ext.Get(untdid.ExtKeyTaxCategory).In(TaxCategoryStandard, TaxCategoryZero),
				tax.ExtensionsExclude(cef.ExtKeyVATEX, ExtKeyExemptionReason),
			),
			validation.When(
				// BR-E-10, BR-AE-10, BR-IC-10, BR-G-10, BR-O-10: all other non-standard
				// categories shall have either a CEF VATEX code or a free-text reason.
				tc.Ext.Get(untdid.ExtKeyTaxCategory).In(
					TaxCategoryExempt, TaxCategoryReverseCharge, TaxCategoryIntraCommunity,
					TaxCategoryExport, TaxCategoryOutsideScope,
				),

				validation.When(
					tc.Ext.Get(cef.ExtKeyVATEX).IsEmpty(),
					tax.ExtensionsRequire(ExtKeyExemptionReason),
				),
				tax.ExtensionsAllowOneOf(cef.ExtKeyVATEX, ExtKeyExemptionReason),
			),
			validation.Skip,
		),
	)
}
