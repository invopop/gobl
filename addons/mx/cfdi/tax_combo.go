package cfdi

import (
	"github.com/invopop/gobl/regimes/mx/sat"
	"github.com/invopop/gobl/tax"
)

func normalizeTaxCombo(tc *tax.Combo) {
	var k tax.ExtValue
	switch tc.Category {
	case sat.TaxCategoryISR:
		k = "001"
	case tax.CategoryVAT, sat.TaxCategoryRVAT:
		k = "002"
	case sat.TaxCategoryIEPS, sat.TaxCategoryRIEPS:
		k = "003"
	default:
		return
	}
	if tc.Ext == nil {
		tc.Ext = make(tax.Extensions)
	}
	tc.Ext[ExtKeyTaxType] = k
}
