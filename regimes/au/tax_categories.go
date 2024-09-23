package au

import (
	"time"

	"github.com/invopop/gobl/cal"
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
		Retained: false,
		Rates: []*tax.RateDef{
			{
				Key: tax.RateExempt,
				Name: i18n.String{
					i18n.EN: "GST-free",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2000, time.July, 1),
						Percent: num.MakePercentage(0, 2),
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
						Since:   cal.NewDate(2000, time.July, 1),
						Percent: num.MakePercentage(10, 2),
					},
				},
			},
			{
				Key: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Input-Taxed",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2000, time.July, 1),
						Percent: num.MakePercentage(0, 2),
					},
				},
			},
		},
	},
	//Source: https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes/wine-equalisation-tax
	{
		Code: TaxCategoryWET,
		Name: i18n.String{
			i18n.EN: "WET",
		},
		Title: i18n.String{
			i18n.EN: "Wine Equalisation Tax",
		},
		Retained: false,
		Rates: []*tax.RateDef{
			{
				Key: tax.RateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2000, time.July, 1),
						Percent: num.MakePercentage(29, 2),
					},
				},
			},
		},
	},
	{
		Code: TaxCategoryLCT,
		Name: i18n.String{
			i18n.EN: "LCT",
		},
		Title: i18n.String{
			i18n.EN: "Luxury Car Tax",
		},
		Retained: false,
		Rates: []*tax.RateDef{
			{
				Key: tax.RateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2000, time.July, 1),
						Percent: num.MakePercentage(33, 2),
					},
				},
			},
		},
	},
}
