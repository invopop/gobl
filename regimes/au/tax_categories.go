package au

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// Tax categories specific for Australia.
const (
	TaxCategoryWET cbc.Code = "WET" // Wine Equalisation Tax
	TaxCategoryLCT cbc.Code = "LCT" // Luxury Car Tax
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
					i18n.EN: "Australian Taxation Office - GST",
				},
				URL: "https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes/gst",
			},
			{
				Title: i18n.String{
					i18n.EN: "Australian Taxation Office - GST-free sales",
				},
				URL: "https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes/gst/when-to-charge-gst-and-when-not-to/gst-free-sales",
			},
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Rate: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
				},
				Description: i18n.String{
					i18n.EN: "Some supplies are zero-rated under the GST, mainly: most basic food, some education courses and materials, some medical and health services, menstrual products, medical aids and medicines, some childcare and religious services, water and sewerage services, precious metals, exports, sales of businesses as going concerns, cars for people with disabilities (when requirements are met), farmland, international transport, and eligible emissions units.",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General Rate",
				},
				Description: i18n.String{
					i18n.EN: "For the majority of sales of goods and services: it applies to all products or services for which no other rate is expressly provided.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2000, 7, 1),
						Percent: num.MakePercentage(100, 2),
					},
				},
			},
		},
	},

	// Wine Equalisation Tax (WET)
	{
		Code: TaxCategoryWET,
		Name: i18n.String{
			i18n.EN: "WET",
		},
		Title: i18n.String{
			i18n.EN: "Wine Equalisation Tax",
		},
		Description: &i18n.String{
			i18n.EN: "A tax of 29% of the wholesale value of wine. Generally only payable if you are registered or required to be registered for GST. Designed to be paid on the last wholesale sale of wine, usually between wholesaler and retailer.",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Australian Taxation Office - Wine Equalisation Tax",
				},
				URL: "https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes/wine-equalisation-tax",
			},
		},
		Retained: false,
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General WET Rate",
				},
				Description: i18n.String{
					i18n.EN: "Applied to the wholesale value of wine.",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(290, 3),
					},
				},
			},
		},
	},

	// Luxury Car Tax (LCT)
	{
		Code: TaxCategoryLCT,
		Name: i18n.String{
			i18n.EN: "LCT",
		},
		Title: i18n.String{
			i18n.EN: "Luxury Car Tax",
		},
		Description: &i18n.String{
			i18n.EN: "A tax on cars that have a GST-inclusive value above the LCT threshold. Imposed at 33% on the amount above the luxury car threshold.",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Australian Taxation Office - Luxury Car Tax",
				},
				URL: "https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes/luxury-car-tax",
			},
		},
		Retained: false,
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General LCT Rate",
				},
				Description: i18n.String{
					i18n.EN: "Applied to the amount above the luxury car threshold.",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(330, 3),
					},
				},
			},
		},
	},
}
