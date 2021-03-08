package be

import (
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// Tax Category definitions
const (
	TaxCategoryVAT tax.Category = "VAT"
)

var taxDefs = tax.Defs{
	/*
	 * VAT
	 */
	{
		Category: TaxCategoryVAT,
		Name:     "VAT Zero Rate",
		Code:     "00",
		Rates: []tax.Rate{
			{
				Value: num.MakePercentage(0, 3),
			},
		},
	},
}
