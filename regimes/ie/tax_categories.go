package ie

import (
	"github.com/invopop/gobl/cal"
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
		Retained: false,
		Rates: []*tax.RateDef{
			{
				Key: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Key: tax.RateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2021, 2, 28),
						Percent: num.MakePercentage(23, 2),
					},
					{
						// Due to Covid
						Since:   cal.NewDate(2020, 9, 1),
						Percent: num.MakePercentage(21, 2),
					},
					{
						Since:   cal.NewDate(2012, 1, 1),
						Percent: num.MakePercentage(23, 2),
					},
				},
			},
			{
				Key: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2003, 1, 1),
						Percent: num.MakePercentage(135, 3),
					},
				},
			},
			{
				Key: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2011, 7, 1),
						Percent: num.MakePercentage(9, 2),
					},
				},
			},
			{
				Key: tax.RateSpecial,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2005, 1, 1),
						Percent: num.MakePercentage(48, 3),
					},
				},
			},
		},
	},
}
