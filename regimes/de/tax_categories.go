package de

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
			i18n.DE: "USt",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.DE: "Mehrwertsteuer",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Federal Ministry of Finance: 100 Years of VAT in Germany",
					i18n.DE: "Bundesfinanzministerium: 100 Jahre Umsatzsteuer in Deutschland",
				},
				URL: "https://www.bundesfinanzministerium.de/Monatsberichte/2019/12/Inhalte/Kapitel-3-Analysen/3-4-100-jahre-umsatzsteuer.html",
			},
			{
				Title: i18n.String{
					i18n.EN: "German Bundestag Research Service: Role of VAT and its Contribution to Financing Social Security Systems",
					i18n.DE: "Wissenschaftliche Dienste des Deutschen Bundestages: Rolle der Umsatzsteuer und ihr Beitrag zur Finanzierung der sozialen Sicherungssysteme",
				},
				URL: "https://www.bundestag.de/resource/blob/410468/8ceeef0b94cdfa0b39a9925dece2aa75/WD-4-040-10-pdf.pdf",
			},
			{
				Title: i18n.String{
					i18n.EN: "Federal Statistical Office: Effect of the VAT Reduction",
					i18n.DE: "Statistisches Bundesamt: Auswirkungen der Mehrwertsteuersenkung",
				},
				URL: "https://www.destatis.de/EN/Themes/Society-Environment/Income-Consumption-Living-Conditions/Consumption-Expenditure/consumption-1-VAT.html",
			},
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General rate",
					i18n.DE: "Allgemeiner Steuersatz",
				},
				Description: i18n.String{
					i18n.EN: "For the majority of sales of goods and services: it applies to all products or services for which no other rate is expressly provided.",
					i18n.DE: "Für den Großteil der Verkäufe von Waren und Dienstleistungen gilt: Dies gilt für alle Produkte oder Dienstleistungen, für die ausdrücklich kein anderer Satz festgelegt ist.",
				},

				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2021, 1, 1),
						Percent: num.MakePercentage(19, 2),
					},
					{
						Since:   cal.NewDate(2020, 7, 1), // COVID temporary measures
						Percent: num.MakePercentage(16, 2),
					},
					{
						Since:   cal.NewDate(2007, 1, 1),
						Percent: num.MakePercentage(19, 2),
					},
					{
						Since:   cal.NewDate(1998, 4, 1),
						Percent: num.MakePercentage(16, 2),
					},
					{
						Since:   cal.NewDate(1993, 1, 1),
						Percent: num.MakePercentage(15, 2),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced rate",
					i18n.DE: "Verminderter Steuersatz",
				},
				Description: i18n.String{
					i18n.EN: "Applicable in particular to basic foodstuffs, books and magazines, cultural events, hotel accommodations, public transportation, medical products, or home renovation.",
					i18n.DE: "Insbesondere anwendbar auf Grundnahrungsmittel, Bücher und Zeitschriften, kulturelle Veranstaltungen, Hotelunterkünfte, öffentliche Verkehrsmittel, medizinische Produkte oder Hausrenovierung.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2021, 1, 1),
						Percent: num.MakePercentage(7, 2),
					},
					{
						Since:   cal.NewDate(2020, 7, 1), // COVID temporary measures
						Percent: num.MakePercentage(5, 2),
					},
					{
						Since:   cal.NewDate(1983, 7, 1),
						Percent: num.MakePercentage(7, 2),
					},
				},
			},
		},
	},
}
