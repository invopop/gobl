package ro

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// Romanian VAT rate information based on:
// - Law 141/2025 (VAT increase and consolidation): https://legislatie.just.ro/Public/DetaliiDocumentAfis/284000
// - Law 227/2015 (Fiscal Code): https://legislatie.just.ro/Public/DetaliiDocument/171282
var taxCategories = []*tax.CategoryDef{
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.RO: "TVA",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.RO: "Taxa pe Valoarea Adăugată",
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
				Description: i18n.String{
					i18n.EN: "Applies to most goods and services unless specified otherwise.",
					i18n.RO: "Se aplică majorității bunurilor și serviciilor, cu excepția cazului în care se specifică altfel.",
				},
				Values: []*tax.RateValueDef{
					{
						// Increased to 21% by Law 141/2025
						Since:   cal.NewDate(2025, 8, 1),
						Percent: num.MakePercentage(210, 3),
					},
					{
						// Standard rate was 19% (2017 - July 2025)
						Since:   cal.NewDate(2017, 1, 1),
						Percent: num.MakePercentage(190, 3),
					},
					{
						// Previous standard rate was 20% (2016)
						Since:   cal.NewDate(2016, 1, 1),
						Percent: num.MakePercentage(200, 3),
					},
					{
						// Previous standard rate was 24% (2010-2015)
						Since:   cal.NewDate(2010, 7, 1),
						Percent: num.MakePercentage(240, 3),
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
					i18n.EN: "Applies to food, pharmaceuticals, books, hotels, restaurants, and water supplies.",
					i18n.RO: "Se aplică la produse alimentare, medicamente, cărți, hoteluri, restaurante și apă.",
				},
				Values: []*tax.RateValueDef{
					{
						// Consolidated reduced rate of 11% (Law 141/2025)
						// Replaces the previous 9% and 5% rates
						Since:   cal.NewDate(2025, 8, 1),
						Percent: num.MakePercentage(110, 3),
					},
					{
						// Previous reduced rate of 9% (2017 - July 2025)
						Since:   cal.NewDate(2017, 1, 1),
						Percent: num.MakePercentage(90, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Super-Reduced Rate",
					i18n.RO: "Cota redusă suplimentară",
				},
				Description: i18n.String{
					i18n.EN: "Legacy rate for textbooks, social housing, etc. (Merged into 11% as of Aug 2025).",
					i18n.RO: "Cotă veche pentru manuale, locuințe sociale etc. (Înglobată în cota de 11% din august 2025).",
				},
				Values: []*tax.RateValueDef{
					{
						// Note: This rate was effectively abolished for new transactions
						// starting Aug 1, 2025, as items moved to the 11% bracket.
						// We keep this entry for historical validation of documents before that date.
						Since:   cal.NewDate(2017, 1, 1),
						Percent: num.MakePercentage(50, 3),
					},
				},
			},
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Law 141/2025 (Fiscal Measures)",
					i18n.RO: "Legea nr. 141/2025 privind măsuri fiscal-bugetare",
				},
				// Official link to the law published in July 2025
				URL: "https://legislatie.just.ro/Public/DetaliiDocumentAfis/284000",
				At:  cal.NewDateTime(2025, 7, 25, 0, 0, 0),
			},
			{
				Title: i18n.String{
					i18n.EN: "Law 227/2015 (Fiscal Code)",
					i18n.RO: "Legea 227/2015 privind Codul fiscal",
				},
				URL: "https://legislatie.just.ro/Public/DetaliiDocument/171282",
				At:  cal.NewDateTime(2024, 12, 17, 0, 0, 0),
			},
		},
	},
}
