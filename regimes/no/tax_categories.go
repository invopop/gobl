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
				Description: i18n.String{
					i18n.EN: "Standard rate for most goods and services.",
					i18n.NB: "Alminnelig sats for de fleste varer og tjenester.",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(250, 3),
						Since:   cal.NewDate(2005, 1, 1),
					},
					{
						Percent: num.MakePercentage(240, 3),
						Since:   cal.NewDate(2001, 1, 1),
					},
					{
						Percent: num.MakePercentage(230, 3),
						Since:   cal.NewDate(1995, 1, 1),
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
					i18n.EN: "Food, beverages, water and wastewater services.",
					i18n.NB: "Næringsmidler, vann og avløpstjenester.",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(150, 3),
						Since:   cal.NewDate(2012, 1, 1),
					},
					{
						Percent: num.MakePercentage(140, 3),
						Since:   cal.NewDate(2007, 1, 1),
					},
					{
						Percent: num.MakePercentage(130, 3),
						Since:   cal.NewDate(2006, 1, 1),
					},
					{
						Percent: num.MakePercentage(110, 3),
						Since:   cal.NewDate(2005, 1, 1),
					},
					{
						Percent: num.MakePercentage(120, 3),
						Since:   cal.NewDate(2001, 7, 1),
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
					i18n.EN: "Transport, accommodation, cinema, broadcasting.",
					i18n.NB: "Transport, overnatting, kino, kringkasting.",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(120, 3),
						Since:   cal.NewDate(2021, 10, 1),
					},
					{
						Percent: num.MakePercentage(60, 3), // COVID temporary reduction
						Since:   cal.NewDate(2020, 4, 1),
					},
					{
						Percent: num.MakePercentage(120, 3),
						Since:   cal.NewDate(2018, 1, 1),
					},
					{
						Percent: num.MakePercentage(100, 3),
						Since:   cal.NewDate(2016, 1, 1),
					},
					{
						Percent: num.MakePercentage(80, 3),
						Since:   cal.NewDate(2006, 1, 1),
					},
					{
						Percent: num.MakePercentage(70, 3),
						Since:   cal.NewDate(2005, 1, 1),
					},
					{
						Percent: num.MakePercentage(60, 3),
						Since:   cal.NewDate(2004, 3, 1),
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
					i18n.EN: "Raw fish (wild marine resources via fiskesalgslag).",
					i18n.NB: "Råfisk (viltlevende marine ressurser via fiskesalgslag).",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(1111, 4),
						// §5-8 merverdiavgiftsloven (2009 Act)
						Since: cal.NewDate(2009, 1, 1),
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
