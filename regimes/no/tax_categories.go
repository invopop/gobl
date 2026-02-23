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
						Percent: num.MakePercentage(250, 3),
						Since:   cal.NewDate(2005, 1, 1),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.NB: "Redusert sats",
				},
				Description: i18n.String{
					i18n.EN: "Food, beverages, water and wastewater services",
					i18n.NB: "Næringsmidler, vann og avløpstjenester",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(150, 3),
						Since:   cal.NewDate(2005, 1, 1),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Super-reduced Rate",
					i18n.NB: "Lav sats",
				},
				Description: i18n.String{
					i18n.EN: "Transport, accommodation, cinema, broadcasting",
					i18n.NB: "Transport, overnatting, kino, kringkasting",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(120, 3),
						Since:   cal.NewDate(2005, 1, 1),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSpecial,
				Name: i18n.String{
					i18n.EN: "Special Rate",
					i18n.NB: "Særskilt sats",
				},
				Description: i18n.String{
					i18n.EN: "Raw fish (wild marine resources via fiskesalgslag)",
					i18n.NB: "Råfisk (viltlevende marine ressurser via fiskesalgslag)",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(1111, 4),
						// Exact historical start date is not traceable in current
						// Skatteetaten documentation; 1970-01-01 is a conservative
						// lower bound.
						Since: cal.NewDate(1970, 1, 1),
					},
				},
			},
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Skatteetaten - VAT rates",
					i18n.NB: "Skatteetaten - Satser for merverdiavgift",
				},
				URL: "https://www.skatteetaten.no/en/rates/value-added-tax/",
				At:  cal.NewDateTime(2025, 6, 15, 0, 0, 0),
			},
			{
				Title: i18n.String{
					i18n.NB: "Lovdata - merverdiavgiftsloven",
				},
				URL: "https://lovdata.no/dokument/NL/lov/2009-06-19-58",
				At:  cal.NewDateTime(2025, 6, 15, 0, 0, 0),
			},
		},
	},
}
