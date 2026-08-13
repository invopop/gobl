package ie

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.CategoryDef{
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.GA: "CBL",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.GA: "Cáin Bhreisluacha",
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.GA: "Ráta Caighdeánach",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(230, 3),
						Since:   cal.NewDate(2012, 1, 1),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "First Reduced Rate",
					i18n.GA: "An Chéad Ráta Laghdaithe",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(135, 3),
						Since:   cal.NewDate(2003, 1, 1),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Second Reduced Rate",
					i18n.GA: "An Dara Ráta Laghdaithe",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(90, 3),
						Since:   cal.NewDate(2011, 1, 1),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSpecial,
				Name: i18n.String{
					i18n.EN: "Livestock Rate",
					i18n.GA: "Ráta Beostoic",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(48, 3),
						Since:   cal.NewDate(2008, 1, 1),
					},
				},
			},
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Revenue - Current VAT Rates",
				},
				URL: "https://www.revenue.ie/en/vat/vat-rates/search-vat-rates/current-vat-rates.aspx",
				At:  cal.NewDateTime(2025, 1, 29, 0, 0, 0),
			},
			{
				Title: i18n.String{
					i18n.EN: "Citizens Information - Value Added Tax",
				},
				URL: "https://www.citizensinformation.ie/en/money-and-tax/tax/duties-and-vat/value-added-tax/",
				At:  cal.NewDateTime(2025, 1, 29, 0, 0, 0),
			},
		},
	},
}
