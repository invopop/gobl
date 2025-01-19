package sg

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.CategoryDef{
	// GST
	{
		Code: tax.CategoryGST,
		Name: i18n.String{
			i18n.EN: "GST",
		},
		Title: i18n.String{
			i18n.EN: "Goods and Services Tax",
		},
		Sources: []*tax.Source{
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
				Key: tax.RateZero,
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
				Key: tax.RateStandard,
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
				},
			},
			{
				Key: tax.RateExempt,
				Name: i18n.String{
					i18n.EN: "Exempt",
				},
				Exempt: true,
				Description: i18n.String{
					i18n.EN: "Certain goods and services are exempt from GST: this includes financial services, sale and lease of residential properties, digital payment tokens, and the import of investment precious metals.",
				},
			},
		},
	},
}
