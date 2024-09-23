package cfdi

import (
	"github.com/invopop/gobl/regimes/mx"
	"github.com/invopop/gobl/tax"
)

func normalizeTaxCombo(tc *tax.Combo) {
	var k tax.ExtValue
	switch tc.Category {
	case mx.TaxCategoryISR:
		k = "001"
	case tax.CategoryVAT, mx.TaxCategoryRVAT:
		k = "002"
	case mx.TaxCategoryIEPS, mx.TaxCategoryRIEPS:
		k = "003"
	default:
		return
	}
	if tc.Ext == nil {
		tc.Ext = make(tax.Extensions)
	}
	tc.Ext[ExtKeyTaxType] = k
}
