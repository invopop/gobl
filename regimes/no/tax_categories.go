// Package no defines VAT tax categories specific to Norway
package no

import (
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
			i18n.NO: "MVA",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.NO: "Merverdiavgift",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Norwegian Tax Administration - VAT Guidelines",
					i18n.NO: "Skatteetaten – Merverdiavgiftsveiledning",
				},
				URL: "https://www.skatteetaten.no",
			},
		},
		Retained: false,
		Rates: []*tax.RateDef{
			{
				Key: tax.RateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.NO: "Standard sats",
				},
				Description: i18n.String{
					i18n.EN: "Applies to most goods and services.",
					i18n.NO: "Gjelder for de fleste varer og tjenester.",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(25, 2),
					},
				},
			},
			{
				Key: "reduced-15",
				Name: i18n.String{
					i18n.EN: "Reduced Rate (15%)",
					i18n.NO: "Redusert sats (15%)",
				},
				Description: i18n.String{
					i18n.EN: "Applies primarily to food and non-alcoholic beverages.",
					i18n.NO: "Gjelder mat og alkoholfrie drikkevarer.",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(15, 2),
					},
				},
			},
			{
				Key: "reduced-12",
				Name: i18n.String{
					i18n.EN: "Reduced Rate (12%)",
					i18n.NO: "Redusert sats (12%)",
				},
				Description: i18n.String{
					i18n.EN: "Applies to passenger transport, hotel accommodation, and cultural events.",
					i18n.NO: "Gjelder persontransport, hotellovernatting og kulturelle arrangementer.",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(12, 2),
					},
				},
			},
			{
				Key: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
					i18n.NO: "Nullsats",
				},
				Description: i18n.String{
					i18n.EN: "Applies to exports and certain international services.",
					i18n.NO: "Gjelder eksport og visse internasjonale tjenester.",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(0, 2),
					},
				},
			},
			{
				Key: tax.RateZero.With(TagBooks),
				Name: i18n.String{
					i18n.EN: "Zero Rate - Books and Periodicals",
					i18n.NO: "Nullsats - Bøker og Tidsskrifter",
				},
				Description: i18n.String{
					i18n.EN: "Zero rate for books (including electronic books and parallel audio-book editions) and newspapers/magazines in the last retail sale (§ 6-4 MVAL).",
					i18n.NO: "Nullsats for bøker (inkludert elektroniske bøker og parallelle lydbokutgaver) og aviser/tidsskrifter i siste detaljsalg (§ 6-4 MVAL).",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(0, 2),
					},
				},
			},
			{
				Key: tax.RateExempt,
				Name: i18n.String{
					i18n.EN: "Exempt",
					i18n.NO: "Fritatt",
				},
				Description: i18n.String{
					i18n.EN: "Exempt from VAT. Applies to financial services, insurance, healthcare, and education.",
					i18n.NO: "Fritatt for merverdiavgift. Gjelder finansielle tjenester, forsikring, helsetjenester og utdanning.",
				},
				Exempt: true,
			},
		},
	},
}
