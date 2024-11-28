package hu

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.CategoryDef{
	// VAT
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.HU: "ÁFA",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.HU: "Általános forgalmi adó",
		},
		Sources: []*tax.Source{
			{
				Title: i18n.String{
					i18n.EN: "Value Added Tax/Goods and Services Tax (VAT/GST) (1976-2023)",
					i18n.HU: "Hozzáadott forgalmi adó/áru- és szolgáltatásadó (ÁFA/GST) (1976-2023)",
				},
				URL: "https://www.oecd.org/tax/tax-policy/tax-database/",
			},
		},
		Retained: false,
		Rates: []*tax.RateDef{
			{
				Key: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
					i18n.HU: "ÁFA-mentes",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2024, 1, 1),
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Key: tax.RateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.HU: "Általános",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(27, 2),
						Since:   cal.NewDate(2012, 1, 1),
					},
					{
						Percent: num.MakePercentage(25, 2),
						Since:   cal.NewDate(2009, 1, 1),
					},
					{
						Percent: num.MakePercentage(20, 2),
						Since:   cal.NewDate(2006, 1, 1),
					},
					{
						Percent: num.MakePercentage(25, 2),
						Since:   cal.NewDate(1988, 1, 1),
					},
				},
			},
			{
				Key: tax.RateIntermediate,
				Name: i18n.String{
					i18n.EN: "Intermediate Rate",
					i18n.HU: "Köztes",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(18, 2),
						Since:   cal.NewDate(2009, 7, 1),
					},
					{
						Percent: num.MakePercentage(5, 2),
						Since:   cal.NewDate(2006, 9, 1),
					},
					{
						Percent: num.MakePercentage(15, 2),
						Since:   cal.NewDate(2004, 1, 1),
					},
					{
						Percent: num.MakePercentage(12, 2),
						Since:   cal.NewDate(1995, 1, 1),
					},
					{
						Percent: num.MakePercentage(10, 2),
						Since:   cal.NewDate(1993, 8, 1),
					},
					{
						Percent: num.MakePercentage(6, 2),
						Since:   cal.NewDate(1993, 1, 1),
					},
					{
						Percent: num.MakePercentage(15, 2),
						Since:   cal.NewDate(1988, 1, 1),
					},
				},
			},
			{
				Key: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.HU: "Csökkentett",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(5, 2),
						Since:   cal.NewDate(2004, 1, 1),
					},
					{
						Percent: num.MakePercentage(0, 2),
						Since:   cal.NewDate(1995, 1, 1),
					},
					{
						Percent: num.MakePercentage(10, 2),
						Since:   cal.NewDate(1993, 8, 1),
					},
					{
						Percent: num.MakePercentage(0, 2),
						Since:   cal.NewDate(1988, 1, 1),
					},
				},
			},
		},
	},
}
