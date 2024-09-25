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
		Sources: []*tax.Source{
			{
				Title: i18n.String{
					i18n.EN: "How GST Works",
				},
				URL: "https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes/gst/how-gst-works",
			},
		},
		Retained: false,
		Rates: []*tax.RateDef{
			{
				Key: tax.RateExempt,
				Name: i18n.String{
					i18n.EN: "Exempt",
				},
				Exempt: true,
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
					i18n.EN: "Zero",
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
	{
		Code: TaxCategoryWET,
		Name: i18n.String{
			i18n.EN: "WET",
		},
		Title: i18n.String{
			i18n.EN: "Wine Equalisation Tax",
		},
		Sources: []*tax.Source{
			{
				Title: i18n.String{
					i18n.EN: "Wine Equalisation Tax",
				},
				URL: "https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes/wine-equalisation-tax",
			},
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
		Sources: []*tax.Source{
			{
				Title: i18n.String{
					i18n.EN: "Luxury Car Tax",
				},
				URL: "https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes/luxury-car-tax",
			},
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
