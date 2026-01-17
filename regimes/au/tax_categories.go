package au

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// Sources:
//  - https://www.ato.gov.au/business/gst
//  - https://www.ato.gov.au/business/gst/when-to-charge-gst

func taxCategories() []*tax.CategoryDef {
	return []*tax.CategoryDef{
		// GST - Goods and Services Tax
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
						i18n.EN: "Australian Taxation Office - GST",
					},
					URL: "https://www.ato.gov.au/business/gst",
				},
			},
			Retained: false,
			Keys:     tax.GlobalGSTKeys(),
			Rates: []*tax.RateDef{
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateGeneral,
					Name: i18n.String{
						i18n.EN: "Standard rate",
					},
					Description: i18n.String{
						i18n.EN: "Standard GST rate applicable to most goods and services in Australia. GST is charged at 10% on most goods, services and other items sold or consumed in Australia.",
					},
					Values: []*tax.RateValueDef{
						{
							// GST introduced on 1 July 2000 at 10%
							Since:   cal.NewDate(2000, 7, 1),
							Percent: num.MakePercentage(10, 2),
						},
					},
				},
				{
					Keys: []cbc.Key{tax.KeyZero},
					Rate: tax.RateZero,
					Name: i18n.String{
						i18n.EN: "Zero-rated",
					},
					Description: i18n.String{
						i18n.EN: "GST-free (zero-rated) supplies including basic food, most health and medical services, educational courses, childcare, exports, and certain other supplies. While these are taxable supplies, GST is charged at 0%.",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2000, 7, 1),
							Percent: num.MakePercentage(0, 3),
						},
					},
				},
			},
		},
	}
}
