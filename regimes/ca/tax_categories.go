package ca

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.CategoryDef{
	//
	// General Sales Tax (GST)
	//
	{
		Code: tax.CategoryGST,
		Name: i18n.String{
			i18n.EN: "GST",
		},
		Title: i18n.String{
			i18n.EN: "General Sales Tax",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "GST/HST provincial rates table",
				},
				URL: "https://www.canada.ca/en/revenue-agency/services/tax/businesses/topics/gst-hst-businesses/charge-collect-which-rate/calculator.html",
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
					i18n.EN: "Some supplies are zero-rated under the GST, mainly: basic groceries, agricultural products, farm livestock, most fishery products such, prescription drugs and drug-dispensing services, certain medical devices, feminine hygiene products, exports, many transportation services where the origin or destination is outside Canada",
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
					i18n.EN: "General rate",
				},
				Description: i18n.String{
					i18n.EN: "For the majority of sales of goods and services: it applies to all products or services for which no other rate is expressly provided.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2022, 1, 1),
						Percent: num.MakePercentage(5, 2),
					},
				},
			},
		},
	},
	//
	// Harmonized Sales Tax (HST)
	//
	{
		Code: TaxCategoryHST,
		Name: i18n.String{
			i18n.EN: "HST",
		},
		Title: i18n.String{
			i18n.EN: "Harmonized Sales Tax",
		},
		// TODO: determine local rates
		Rates: []*tax.RateDef{},
	},
	//
	// Provincial Sales Tax (PST)
	//
	{
		Code: TaxCategoryPST,
		Name: i18n.String{
			i18n.EN: "PST",
		},
		Title: i18n.String{
			i18n.EN: "Provincial Sales Tax",
		},
		// TODO: determine local rates
		Rates: []*tax.RateDef{},
	},
}
