package nz

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// TaxRateAccommodation is the key for the long-term accommodation rate.
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
					i18n.EN: "Applies to all taxable supplies of goods and services that are not zero-rated or exempt.",
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
					i18n.EN: "Applies to exported goods and services, international transport, land transactions between GST-registered parties, going concern sales, and duty-free goods.",
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
					i18n.EN: "Applies to commercial accommodation provided for 28 or more consecutive days. Effective from 1 April 2024.",
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
