package zatca

import (
	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/catalogues/cef"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// SA VATEX exemption reason codes (subset of CEF VATEX list, per ZATCA spec).
const (
	VatexFinancialServices        cbc.Code = "VATEX-SA-29"
	VatexLifeInsurance            cbc.Code = "VATEX-SA-29-7"
	VatexRealEstate               cbc.Code = "VATEX-SA-30"
	VatexExportGoods              cbc.Code = "VATEX-SA-32"
	VatexExportServices           cbc.Code = "VATEX-SA-33"
	VatexIntlTransportGoods       cbc.Code = "VATEX-SA-34-1"
	VatexIntlTransportPassengers  cbc.Code = "VATEX-SA-34-2"
	VatexIntlTransportRelated     cbc.Code = "VATEX-SA-34-3"
	VatexQualifyingTransportMeans cbc.Code = "VATEX-SA-34-4"
	VatexTransportRelated         cbc.Code = "VATEX-SA-34-5"
	VatexMedicines                cbc.Code = "VATEX-SA-35"
	VatexQualifyingMetals         cbc.Code = "VATEX-SA-36"
	VatexPrivateEducation         cbc.Code = "VATEX-SA-EDU"
	VatexPrivateHealthcare        cbc.Code = "VATEX-SA-HEA"
	VatexMilitaryGoods            cbc.Code = "VATEX-SA-MLTRY"
	VatexOutOfScope               cbc.Code = "VATEX-SA-OOS"
)

var vatexValidCodes = map[cbc.Code][]cbc.Code{
	en16931.TaxCategoryExempt: {
		VatexFinancialServices,
		VatexLifeInsurance,
		VatexRealEstate,
	},
	en16931.TaxCategoryStandard: {},
	en16931.TaxCategoryZero: {
		VatexExportGoods,
		VatexExportServices,
		VatexIntlTransportGoods,
		VatexIntlTransportPassengers,
		VatexIntlTransportRelated,
		VatexQualifyingTransportMeans,
		VatexTransportRelated,
		VatexMedicines,
		VatexQualifyingMetals,
		VatexPrivateEducation,
		VatexPrivateHealthcare,
		VatexMilitaryGoods,
	},
	en16931.TaxCategoryOutsideScope: {
		VatexOutOfScope,
	},
}

var vatKeyMap = tax.Extensions{
	tax.KeyStandard:     en16931.TaxCategoryStandard,
	tax.KeyZero:         en16931.TaxCategoryZero,
	tax.KeyExempt:       en16931.TaxCategoryExempt,
	tax.KeyOutsideScope: en16931.TaxCategoryOutsideScope,
}

func normalizeTaxCombo(tc *tax.Combo) {
	if tc == nil || tc.Category != tax.CategoryVAT {
		return
	}
	if tc.Key.IsEmpty() {
		k := vatKeyMap.Lookup(tc.Ext.Get(untdid.ExtKeyTaxCategory))
		if k.IsEmpty() {
			k = tax.KeyStandard
		}
		tc.Key = k
	}
	if cat := vatKeyMap.Get(tc.Key); cat != "" {
		tc.Ext = tc.Ext.Set(untdid.ExtKeyTaxCategory, cat)
	}
}

func taxComboRules() *rules.Set {
	return rules.For(new(tax.Combo),

		rules.Field("ext",
			rules.Assert("01", "VATEX exemption code must be present and valid for Z/E/O categories, and must not be set for Standard (BR-KSA-CL-04)",
				is.Func("valid SA VATEX code", taxComboHasValidVATEX),
			),
		),

		// Extensions
		rules.Field("ext",
			rules.Assert("02", "VAT category code must contain one of the values (S, Z, E, O) (BR-KSA-18)",
				tax.ExtensionsHasCodes(untdid.ExtKeyTaxCategory,
					en16931.TaxCategoryStandard,
					en16931.TaxCategoryZero,
					en16931.TaxCategoryExempt,
					en16931.TaxCategoryOutsideScope,
				),
			),
		),

		// Category
		rules.Field("cat",
			rules.Assert("04", "tax schema id must be 'VAT' (BR-KSA-54)", is.In(tax.CategoryVAT)),
		),
	)
}

func taxComboHasValidVATEX(val any) bool {
	ext, ok := val.(tax.Extensions)
	if !ok {
		return false
	}
	category := ext.Get(untdid.ExtKeyTaxCategory)
	vatex := ext.Get(cef.ExtKeyVATEX)

	switch category {
	case en16931.TaxCategoryStandard:
		return vatex == cbc.CodeEmpty
	case en16931.TaxCategoryExempt,
		en16931.TaxCategoryZero,
		en16931.TaxCategoryOutsideScope:
		allowed, ok := vatexValidCodes[category]
		return ok && vatex.In(allowed...)
	case cbc.CodeEmpty:
		return true
	default:
		return false
	}
}
