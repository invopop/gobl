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
			i18n.NO: "MVA",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.NO: "Merverdiavgift",
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Normal Rate",
					i18n.NO: "Normalsats",
				},
				Values: []*tax.RateValueDef{
					// Current rate snapshot (historical changes out of scope)
					{Percent: num.MakePercentage(250, 3)},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate (Foodstuffs)",
					i18n.NO: "Redusert sats (mat)",
				},
				Values: []*tax.RateValueDef{
					{Percent: num.MakePercentage(150, 3)},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate (Passenger transport / accommodation / cinema)",
					i18n.NO: "Redusert sats (transport/overnatting)",
				},
				Values: []*tax.RateValueDef{
					{Percent: num.MakePercentage(120, 3)},
				},
			},
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "The Norwegian Tax Administration - VAT rates",
				},
				URL: "https://www.skatteetaten.no/en/rates/value-added-tax/",
				At:  cal.NewDateTime(2026, 2, 24, 0, 0, 0),
			},
			{
				Title: i18n.String{
					i18n.EN: "Altinn - Value added tax (VAT) rates",
				},
				URL: "https://info.altinn.no/en/start-and-run-business/direct-and-indirect-taxes/indirect-taxes/value-added-tax/",
				At:  cal.NewDateTime(2026, 2, 24, 0, 0, 0),
			},
		},
	},
}
