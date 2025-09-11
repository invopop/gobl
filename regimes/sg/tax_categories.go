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
					URL: "https://www.iras.gov.sg/taxes/goods-services-tax-(gst)/basics-of-gst/current-gst-rates",
				},
			},
			Retained: false,
			Keys:     tax.GlobalGSTKeys(),
			Rates: []*tax.RateDef{
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateGeneral,
					Name: i18n.String{
						i18n.EN: "General rate",
					},
					Description: i18n.String{
						i18n.EN: "GST-registered businesses are required to charge and account for GST at 9% on all sales of goods and services in Singapore unless the sale can be zero-rated or exempted under the GST law.",
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
