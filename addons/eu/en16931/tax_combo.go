package en16931

import (
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var taxCategoryMap = tax.Extensions{
	tax.RateStandard: "S",
	tax.RateReduced:  "S", // Same as standard
	tax.RateZero:     "Z",
	tax.RateExempt:   "E",
	tax.RateExempt.With(tax.TagReverseCharge):           "AE",
	tax.RateExempt.With(tax.TagExport).With(tax.TagEEA): "K",
	tax.RateExempt.With(tax.TagExport):                  "G",
}

// acceptedTaxCategories as defined by the EN 16931 code list values data.
var acceptedTaxCategories = []tax.ExtValue{
	"S", "Z", "E", "AE", "K", "G", "O", "L", "M",
}

func normalizeTaxCombo(tc *tax.Combo) {
	switch tc.Category {
	case tax.CategoryVAT:
		if tc.Rate.IsEmpty() {
			return
		}
		k, ok := taxCategoryMap[tc.Rate]
		if !ok {
			return
		}
		tc.Ext = tc.Ext.Merge(
			tax.Extensions{untdid.ExtKeyTaxCategory: k},
		)
	case es.TaxCategoryIGIC:
		tc.Ext = tc.Ext.Merge(
			tax.Extensions{untdid.ExtKeyTaxCategory: "L"},
		)
	case es.TaxCategoryIPSI:
		tc.Ext = tc.Ext.Merge(
			tax.Extensions{untdid.ExtKeyTaxCategory: "M"},
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
			tax.ExtensionsHasValues(untdid.ExtKeyTaxCategory, acceptedTaxCategories...),
			validation.Skip,
		),
	)
}
