package ro

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.CategoryDef{
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.RO: "TVA",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.RO: "Taxa pe Valoarea Adăugată",
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.RO: "Cota Standard",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(210, 3),
						Since:   cal.NewDate(2025, 8, 1),
					},
					{
						Percent: num.MakePercentage(190, 3),
						Since:   cal.NewDate(2017, 1, 1),
					},
					{
						Percent: num.MakePercentage(240, 3),
						Since:   cal.NewDate(2010, 7, 1),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "First Reduced Rate",
					i18n.RO: "Prima Cota Redusă",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(110, 3), 
						Since:   cal.NewDate(2025, 8, 1),
					},
					{
						Percent: num.MakePercentage(90, 3), 
						Since:   cal.NewDate(2017, 1, 1),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Second Reduced Rate",
					i18n.RO: "A Doua Cota Redusă",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(110, 3), 
						Since:   cal.NewDate(2025, 8, 1),
					},
					{
						Percent: num.MakePercentage(50, 3),
						Since:   cal.NewDate(2011, 1, 1),
					},					
				},
			},			
		},
	},
}
