package sg

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

func taxCategories() []*tax.CategoryDef {
	return []*tax.CategoryDef{
		// GST
		{
			Code: tax.CategoryGST,
			Name: i18n.String{
				i18n.EN: "GST",
			},
			Title: i18n.String{
				i18n.EN: "Goods and Services Tax",
			},
			Sources: []*cbc.Source{
				{
					Title: i18n.String{
						i18n.EN: "Goods and Services Tax (GST)",
					},
					URL: "https://www.iras.gov.sg/taxes/goods-services-tax-(gst)/",
				},
			},
			Retained: false,
			Rates: []*tax.RateDef{
				{
					Rate: tax.RateZero,
					Name: i18n.String{
						i18n.EN: "Zero Rate",
					},
					Description: i18n.String{
						i18n.EN: "Zero-rated supplies are goods and services that are taxable at 0%: this referes to international services and export of goods.",
					},

					Values: []*tax.RateValueDef{
						{
							Percent: num.MakePercentage(0, 3),
						},
					},
				},
				{
					Rate: tax.RateGeneral,
					Name: i18n.String{
						i18n.EN: "Standard rate",
					},
					Description: i18n.String{
						i18n.EN: "For the majority of sales of goods and services: it applies to all products or services for which no other rate is expressly provided.",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2024, 1, 1),
							Percent: num.MakePercentage(9, 2),
						},
						{
							Since:   cal.NewDate(2023, 1, 1),
							Percent: num.MakePercentage(8, 2),
						},
						{
							Since:   cal.NewDate(2007, 7, 1),
							Percent: num.MakePercentage(7, 2),
						},
						{
							Since:   cal.NewDate(2004, 1, 1),
							Percent: num.MakePercentage(5, 2),
						},
						{
							Since:   cal.NewDate(2003, 1, 1),
							Percent: num.MakePercentage(4, 2),
						},
						{
							Since:   cal.NewDate(1994, 4, 1),
							Percent: num.MakePercentage(3, 2),
						},
					},
				},
			},
		},
	}
}
