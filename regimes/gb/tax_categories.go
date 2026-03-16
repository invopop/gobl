package gb

import (
	"github.com/invopop/gobl/cal"
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
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Value Added Tax/Goods and Services Tax (VAT/GST) (1976-2024)",
				},
				URL: "https://www.oecd.org/tax/tax-policy/tax-database/",
			},
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Standard rate",
				},
				Description: i18n.String{
					i18n.EN: "Applies to most goods and services unless specifically exempted or subject to a reduced or zero rate.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2011, 1, 4),
						Percent: num.MakePercentage(200, 3),
					},
					{
						Since:   cal.NewDate(2010, 1, 1),
						Percent: num.MakePercentage(175, 3),
					},
					{
						Since:   cal.NewDate(2008, 12, 1),
						Percent: num.MakePercentage(150, 3),
					},
					{
						Since:   cal.NewDate(1991, 4, 1),
						Percent: num.MakePercentage(175, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced rate",
				},
				Description: i18n.String{
					i18n.EN: "Applies to certain goods and services including domestic fuel and power, children's car seats, and some energy-saving materials.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(1997, 9, 1),
						Percent: num.MakePercentage(50, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyZero},
				Rate: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero rate",
				},
				Description: i18n.String{
					i18n.EN: "Applies to most food, children's clothing, books, newspapers, and public transport.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(1973, 4, 1),
						Percent: num.MakePercentage(0, 2),
					},
				},
			},
		},
	},
}
