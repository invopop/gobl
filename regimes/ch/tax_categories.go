package ch

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
					i18n.EN: "Current Swiss VAT rates | FTA",
				},
				URL: "https://www.estv.admin.ch/estv/en/home/value-added-tax/vat-rates-switzerland.html",
			},
		},
		Retained: false,
		Rates: []*tax.RateDef{
			{
				Key: tax.RateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
				},
				Description: i18n.String{
					i18n.EN: "Applies to most goods and services.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2024, 1, 1),
						Percent: num.MakePercentage(81, 3),
					},
				},
			},
			{
				Key: tax.RateIntermediate,
				Name: i18n.String{
					i18n.EN: "Intermediate Rate",
				},
				Description: i18n.String{
					i18n.EN: "Applies to accommodation services.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2024, 1, 1),
						Percent: num.MakePercentage(38, 3),
					},
				},
			},
			{
				Key: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
				},
				Description: i18n.String{
					i18n.EN: "Applies to food, books, newspapers, and public transport.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2024, 1, 1),
						Percent: num.MakePercentage(26, 3),
					},
				},
			},
		},
	},
}
