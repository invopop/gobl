package de

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.Category{
	//
	// VAT
	//
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.DE: "MwSt",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.DE: "Mehrwertsteuer",
		},
		Sources: []*tax.Source{
			{
				Title: i18n.String{
					i18n.EN: "Value Added Tax/Goods and Services Tax (VAT/GST) (1976-2023)",
					i18n.DE: "Umsatzsteuer/Güter - und Dienstleistungssteuer (USt/GST) (1976-2023)",
				},
				URL: "https://www.oecd.org/tax/tax-policy/tax-database/",
			},
		},
		Retained: false,
		Rates: []*tax.Rate{
			{
				Key: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
					i18n.DE: "Nullsatz",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Key: tax.RateStandard,
				Name: i18n.String{
					i18n.EN: "Standard rate",
					i18n.DE: "Standardsteuersatz",
				},
				Description: i18n.String{
					i18n.EN: "For the majority of sales of goods and services: it applies to all products or services for which no other rate is expressly provided.",
					i18n.DE: "Für den Großteil der Verkäufe von Waren und Dienstleistungen gilt: Dies gilt für alle Produkte oder Dienstleistungen, für die ausdrücklich kein anderer Satz festgelegt ist.",
				},

				Values: []*tax.RateValue{
					{
						Since:   cal.NewDate(2022, 1, 1),
						Percent: num.MakePercentage(19, 2),
					},
					{
						Since:   cal.NewDate(2020, 7, 1), // COVID temporary measures
						Percent: num.MakePercentage(16, 2),
					},
					{
						Since:   cal.NewDate(2007, 7, 1),
						Percent: num.MakePercentage(19, 2),
					},
					{
						Since:   cal.NewDate(1993, 1, 1),
						Percent: num.MakePercentage(16, 2),
					},
				},
			},
			{
				Key: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced rate",
					i18n.DE: "Verminderter Steuersatz",
				},
				Description: i18n.String{
					i18n.EN: "Applicable in particular to basic foodstuffs, books and magazines, cultural events, hotel accommodations, public transportation, medical products, or home renovation.",
					i18n.DE: "Insbesondere anwendbar auf Grundnahrungsmittel, Bücher und Zeitschriften, kulturelle Veranstaltungen, Hotelunterkünfte, öffentliche Verkehrsmittel, medizinische Produkte oder Hausrenovierung.",
				},
				Values: []*tax.RateValue{
					{
						Since:   cal.NewDate(2022, 1, 1),
						Percent: num.MakePercentage(7, 2),
					},
					{
						Since:   cal.NewDate(2020, 7, 1), // COVID temporary measures
						Percent: num.MakePercentage(5, 2),
					},
					{
						Since:   cal.NewDate(2007, 7, 1),
						Percent: num.MakePercentage(7, 2),
					},
					{
						Since:   cal.NewDate(1993, 1, 1),
						Percent: num.MakePercentage(5, 2),
					},
				},
			},
		},
	},
}
