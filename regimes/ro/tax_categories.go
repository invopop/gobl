package ro

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.CategoryDef{
	//
	// VAT
	//
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.RO: "TVA",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.RO: "Taxa pe valoarea adăugată",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "ANAF - Tax rates",
					i18n.RO: "ANAF - Cote de impozitare",
				},
				URL: "https://www.anaf.ro/",
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
					i18n.RO: "Cota standard",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2025, 8, 1),
						Percent: num.MakePercentage(21, 2),
					},
					{
						Since:   cal.NewDate(2017, 1, 1),
						Percent: num.MakePercentage(19, 2),
					},
					{
						Since:   cal.NewDate(2016, 1, 1),
						Percent: num.MakePercentage(20, 2),
					},
					{
						Since:   cal.NewDate(2010, 7, 1),
						Percent: num.MakePercentage(24, 2),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.RO: "Cota redusă",
				},
				Description: i18n.String{
					i18n.EN: "Applicable to food, non-alcoholic beverages, hotel accommodation, restaurants, books, newspapers, medical products, and other specified goods and services.",
					i18n.RO: "Aplicabilă pentru alimente, băuturi nealcoolice, cazare hotelieră, restaurante, cărți, ziare, produse medicale și alte bunuri și servicii specificate.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2025, 8, 1),
						Percent: num.MakePercentage(11, 2),
					},
					{
						Since:   cal.NewDate(2017, 1, 1),
						Percent: num.MakePercentage(9, 2),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Super-Reduced Rate",
					i18n.RO: "Cota super-redusă",
				},
				Description: i18n.String{
					i18n.EN: "Applicable to social housing, museums, zoos, botanical gardens, historic monuments, and other specified goods and services.",
					i18n.RO: "Aplicabilă pentru locuințe sociale, muzee, grădini zoologice, grădini botanice, monumente istorice și alte bunuri și servicii specificate.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2025, 8, 1),
						Percent: num.MakePercentage(11, 2),
					},
					{
						Since:   cal.NewDate(2017, 1, 1),
						Percent: num.MakePercentage(5, 2),
					},
				},
			},
		},
	},
}
