package nz

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// TaxRateAccommodation is the rate key for long-term commercial accommodation.
const TaxRateAccommodation cbc.Key = "accommodation"

var taxCategories = []*tax.CategoryDef{
	{
		Code: tax.CategoryGST,
		Name: i18n.String{
			i18n.EN: "GST",
		},
		Title: i18n.String{
			i18n.EN: "Goods and Services Tax",
		},
		Keys: tax.GlobalGSTKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
				},
				Description: i18n.String{
					i18n.EN: "Standard GST rate applicable to most goods and services.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2010, 10, 1),
						Percent: num.MakePercentage(150, 3),
					},
					{
						Since:   cal.NewDate(1989, 7, 1),
						Percent: num.MakePercentage(125, 3),
					},
					{
						Since:   cal.NewDate(1986, 10, 1),
						Percent: num.MakePercentage(100, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyZero},
				Rate: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
				},
				Description: i18n.String{
					i18n.EN: "Zero-rated supplies including exports, international services, and certain land transactions between GST-registered parties.",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: TaxRateAccommodation,
				Name: i18n.String{
					i18n.EN: "Long-term Accommodation",
				},
				Description: i18n.String{
					i18n.EN: "Reduced rate for commercial accommodation provided for 28 or more consecutive days.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2024, 4, 1),
						Percent: num.MakePercentage(90, 3),
					},
				},
			},
		},
	},
}
