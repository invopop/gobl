package nz

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

// TaxRateAccommodation is the key for the long-term accommodation rate.
const TaxRateAccommodation cbc.Key = "accommodation"

var taxCategories = []*tax.CategoryDef{
	{
		Code: tax.CategoryGST,
		Name: i18n.String{
			i18n.EN: "GST",
		},
		Title: i18n.String{
			i18n.EN: "Goods and Services Tax",
		},
		Description: &i18n.String{
			i18n.EN: here.Doc(`
				New Zealand's Goods and Services Tax (GST) is a consumption tax applied to
				most goods and services sold in New Zealand. Introduced in 1986, it operates
				as a value-added tax collected at each stage of the supply chain. Businesses
				registered for GST charge the tax on their sales and can claim back GST paid
				on business purchases. The current standard rate is 15%.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Goods and Services Tax Act 1985",
				},
				URL: "https://www.legislation.govt.nz/act/public/1985/0141/latest/DLM81035.html",
			},
			{
				Title: i18n.String{
					i18n.EN: "IRD - Charging GST",
				},
				URL: "https://www.ird.govt.nz/gst/charging-gst",
			},
		},
		Keys: tax.GlobalGSTKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
				},
				Description: i18n.String{
					i18n.EN: "Applies to all taxable supplies of goods and services that are not zero-rated or exempt.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2010, 10, 1),
						Percent: num.MakePercentage(150, 3),
					},
					{
						Since:   cal.NewDate(1989, 7, 1),
						Percent: num.MakePercentage(125, 3),
					},
					{
						Since:   cal.NewDate(1986, 10, 1),
						Percent: num.MakePercentage(100, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyZero},
				Rate: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
				},
				Description: i18n.String{
					i18n.EN: "Applies to exported goods and services, international transport, land transactions between GST-registered parties, going concern sales, and duty-free goods.",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: TaxRateAccommodation,
				Name: i18n.String{
					i18n.EN: "Long-term Accommodation",
				},
				Description: i18n.String{
					i18n.EN: "Applies to commercial accommodation provided for 28 or more consecutive days. Effective from 1 April 2024.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2024, 4, 1),
						Percent: num.MakePercentage(90, 3),
					},
				},
			},
		},
	},
}
