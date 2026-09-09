package ee

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
			i18n.ET: "KM",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.ET: "Käibemaks",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Estonian Tax and Customs Board - Value Added Tax"),
				URL:   "https://www.emta.ee/en/business-client/taxes-and-payment/value-added-tax",
			},
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.ET: "Tavakäibemaks",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2025, 7, 1),
						Percent: num.MakePercentage(240, 3), // 24.0%
					},
					{
						Since:   cal.NewDate(2024, 1, 1),
						Percent: num.MakePercentage(220, 3), // 22.0%
					},
					{
						Since:   cal.NewDate(2009, 7, 1),
						Percent: num.MakePercentage(200, 3), // 20.0%
					},
					{
						Since:   cal.NewDate(1993, 1, 1),
						Percent: num.MakePercentage(180, 3), // 18.0%
					},
					{
						Since:   cal.NewDate(1991, 1, 1),
						Percent: num.MakePercentage(100, 3), // 10.0%
					},
				},
			},
			{ // intermediate rate for "accommodation or accommodation with breakfast"
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateIntermediate,
				Name: i18n.String{
					i18n.EN: "Intermediate Rate",
					i18n.ET: "Keskmine käibemaks",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2025, 1, 1),    // previous to 1 Jan 2025 this was included in 9% reduced rate
						Percent: num.MakePercentage(130, 3), // 13.0%
					},
				},
			},
			{ // reduced rate for books, medical equipment, press publications
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.ET: "Alandatud käibemaks",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2013, 1, 1),
						Percent: num.MakePercentage(90, 3), // 9.0%
					},
				},
			},
			{ // zero rate for exports, intra-community supplies, certain other goods & services
				Keys: []cbc.Key{tax.KeyZero},
				Rate: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
					i18n.ET: "Nullkäibemaks",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(0, 3), // 0.0%
					},
				},
			},
		},
	},
}
