package it

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

var schemes = []*tax.Scheme{
	// Reverse Charge Scheme
	{
		Key: common.SchemeReverseCharge,
		Name: i18n.String{
			i18n.EN: "Reverse Charge",
			i18n.IT: "Inversione Contabile",
		},
		Categories: []cbc.Code{
			common.TaxCategoryVAT,
		},
	},
}
