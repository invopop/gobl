package cz

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
			i18n.CS: "DPH",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.CS: "Daň z přidané hodnoty",
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			// Standard rate
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.CS: "Základní sazba",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(210, 3),
						Since:   cal.NewDate(2013, 1, 1),
					},
					{
						Percent: num.MakePercentage(200, 3),
						Since:   cal.NewDate(2010, 1, 1),
					},
					{
						Percent: num.MakePercentage(190, 3),
						Since:   cal.NewDate(2004, 5, 1),
					},
				},
			},
			// Reduced rate (merged from two reduced rates in 2024)
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.CS: "Snížená sazba",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(120, 3),
						Since:   cal.NewDate(2024, 1, 1),
					},
					{
						Percent: num.MakePercentage(150, 3),
						Since:   cal.NewDate(2013, 1, 1),
					},
					{
						Percent: num.MakePercentage(140, 3),
						Since:   cal.NewDate(2012, 1, 1),
					},
					{
						Percent: num.MakePercentage(100, 3),
						Since:   cal.NewDate(2010, 1, 1),
					},
					{
						Percent: num.MakePercentage(90, 3),
						Since:   cal.NewDate(2008, 1, 1),
					},
					{
						Percent: num.MakePercentage(50, 3),
						Since:   cal.NewDate(2004, 5, 1),
					},
				},
			},
			// Second reduced rate (abolished 2024-01-01, merged into single 12% reduced rate)
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Second Reduced Rate",
					i18n.CS: "Druhá snížená sazba",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(100, 3),
						Since:   cal.NewDate(2015, 1, 1),
					},
				},
			},
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Czech Republic - General rules and VAT rates",
				},
				URL: "https://portal.gov.cz/en/informace/general-rules-and-vat-rates-INF-205",
			},
			{
				Title: i18n.String{
					i18n.EN: "Registering for VAT in the Czech Republic",
				},
				URL: "https://portal.gov.cz/en/informace/registering-for-vat-INF-204",
			},
		},
	},
}
