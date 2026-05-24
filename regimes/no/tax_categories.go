package no

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
			i18n.NB: "MVA",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.NB: "Merverdiavgift",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Skatteetaten - VAT rates",
				},
				URL: "https://www.skatteetaten.no/en/rates/value-added-tax/",
			},
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General Rate",
					i18n.NB: "Alminnelig sats",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(1995, 1, 1),
						Percent: num.MakePercentage(250, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate (Food, Water Supply)",
					i18n.NB: "Redusert sats (mat, vannforsyning)",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2005, 1, 1),
						Percent: num.MakePercentage(150, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Super-Reduced Rate (Transport, Accommodation)",
					i18n.NB: "Lav sats (transport, overnatting)",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2018, 1, 1),
						Percent: num.MakePercentage(120, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyZero},
				Rate: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
					i18n.NB: "Nullsats",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(1970, 1, 1),
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
		},
	},
}
