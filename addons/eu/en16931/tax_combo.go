package en16931

import (
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Official subset of UNTDID 5305 category codes recognized by the EN 16931
const (
	TaxCategoryStandard      cbc.Code = "S"
	TaxCategoryZero          cbc.Code = "Z"
	TaxCategoryExempt        cbc.Code = "E"
	TaxCategoryReverseCharge cbc.Code = "AE"
	TaxCategoryExportEEA     cbc.Code = "K"
	TaxCategoryExport        cbc.Code = "G"
	TaxCategoryOutsideScope  cbc.Code = "O"
	TaxCategoryIGIC          cbc.Code = "L" // Canary Islands
	TaxCategoryIPSI          cbc.Code = "M" // Ceuta and Melilla
)

var vatRateCategoryMap = tax.Extensions{
	tax.RateStandard: TaxCategoryStandard,
	tax.RateReduced:  TaxCategoryStandard, // Same as standard
	tax.RateZero:     TaxCategoryZero,
	tax.RateExempt:   TaxCategoryExempt,
	tax.RateExempt.With(tax.TagReverseCharge):           TaxCategoryReverseCharge,
	tax.RateExempt.With(tax.TagExport).With(tax.TagEEA): TaxCategoryExportEEA,
	tax.RateExempt.With(tax.TagExport):                  TaxCategoryExport,
}

// acceptedTaxCategories as defined by the EN 16931 code list values data.
var vatAppliesTaxCategories = []cbc.Code{
	TaxCategoryStandard,
	TaxCategoryZero,
}

var vatExemptTaxCategories = []cbc.Code{
	TaxCategoryExempt,
	TaxCategoryReverseCharge,
	TaxCategoryExportEEA,
	TaxCategoryExport,
}

func normalizeTaxCombo(tc *tax.Combo) {
	switch tc.Category {
	case tax.CategoryVAT:
		if tc.Rate.IsEmpty() {
			return
		}
		k, ok := vatRateCategoryMap[tc.Rate]
		if !ok {
			return
		}
		tc.Ext = tc.Ext.Merge(
			tax.Extensions{untdid.ExtKeyTaxCategory: k},
		)
	case es.TaxCategoryIGIC:
		tc.Ext = tc.Ext.Merge(
			tax.Extensions{untdid.ExtKeyTaxCategory: TaxCategoryIGIC},
		)
	case es.TaxCategoryIPSI:
		tc.Ext = tc.Ext.Merge(
			tax.Extensions{untdid.ExtKeyTaxCategory: TaxCategoryIPSI},
		)
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
				tc.Category == tax.CategoryVAT && tc.Percent != nil,
				tax.ExtensionsHasCodes(untdid.ExtKeyTaxCategory, vatAppliesTaxCategories...),
			),
			validation.When(
				tc.Category == tax.CategoryVAT && tc.Percent == nil,
				tax.ExtensionsHasCodes(untdid.ExtKeyTaxCategory, vatExemptTaxCategories...),
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
			validation.Skip,
		),
	)
}
