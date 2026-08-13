package nl

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.CategoryDef{
	//
	// VAT
	//
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.NL: "BTW",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.NL: "Belasting Toegevoegde Waarde",
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General Rate",
					i18n.NL: "Algemeen Tarief",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(210, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.NL: "Gereduceerd Tarief",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(90, 3),
					},
				},
			},
		},
	},
}
