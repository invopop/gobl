package au

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.CategoryDef{
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
					i18n.EN: "GST overview",
				},
				URL: "https://www.ato.gov.au/about-ato/research-and-statistics/in-detail/tax-gap/goods-and-services-tax-gap/overview",
			},
			{
				Title: i18n.String{
					i18n.EN: "GST-free sales",
				},
				URL: "https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes/gst/when-to-charge-gst-and-when-not-to/gst-free-sales",
			},
			{
				Title: i18n.String{
					i18n.EN: "Input taxed sales",
				},
				URL: "https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes/gst/when-to-charge-gst-and-when-not-to/input-taxed-sales",
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
					i18n.EN: "GST is 10% on most goods and services consumed in Australia, including imports and digital products.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2000, 7, 1),
						Percent: num.MakePercentage(10, 2),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyZero},
				Rate: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "GST-free",
				},
				Description: i18n.String{
					i18n.EN: "GST-free sales are not taxed. Common examples include basic food, some health and education services, and exports.",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyExempt},
				Rate: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Input-taxed",
				},
				Description: i18n.String{
					i18n.EN: "Input-taxed sales do not include GST in the price and no GST credits can be claimed.",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
		},
	},
}
