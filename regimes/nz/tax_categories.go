package nz

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

func taxCategories() []*tax.CategoryDef {
	return []*tax.CategoryDef{
		//
		// GST
		//
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
						i18n.EN: "GST - Inland Revenue New Zealand",
					},
					URL: "https://www.ird.govt.nz/gst",
				},
			},
			Retained: false,
			Keys:     tax.GlobalGSTKeys(),
			Rates: []*tax.RateDef{
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateGeneral,
					Name: i18n.String{
						i18n.EN: "Standard Rate",
					},
					Description: i18n.String{
						i18n.EN: "Applies to most goods and services in New Zealand.",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2010, 10, 1),
							Percent: num.MakePercentage(15, 2),
						},
						{
							Since:   cal.NewDate(1989, 7, 1),
							Percent: num.MakePercentage(125, 3),
						},
						{
							Since:   cal.NewDate(1986, 10, 1),
							Percent: num.MakePercentage(10, 2),
						},
					},
				},
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateReduced,
					Name: i18n.String{
						i18n.EN: "Accommodation Rate",
					},
					Description: i18n.String{
						i18n.EN: "Applies to long-term accommodation of 28 or more consecutive days. Equivalent to 60% of the standard rate applied to the full charge.",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2010, 10, 1),
							Percent: num.MakePercentage(9, 2),
						},
						{
							Since:   cal.NewDate(1989, 7, 1),
							Percent: num.MakePercentage(75, 3),
						},
						{
							Since:   cal.NewDate(1986, 10, 1),
							Percent: num.MakePercentage(6, 2),
						},
					},
				},
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateZero,
					Name: i18n.String{
						i18n.EN: "Zero Rate",
					},
					Description: i18n.String{
						i18n.EN: "Zero-rated supplies including exported goods, international transport, and certain land transactions between GST-registered parties.",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(1986, 10, 1),
							Percent: num.MakePercentage(0, 2),
						},
					},
				},
			},
		},
	}
}
