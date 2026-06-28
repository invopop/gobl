package lv

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
			i18n.LV: "PVN",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.LV: "Pievienotās Vertības Nodoklis",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "State Revenue Service - VAT overview",
				},
				URL: "https://www.fm.gov.lv/lv/tax-rates",
			},
			{
				Title: i18n.String{
					i18n.EN: "Value Added Tax Law",
				},
				URL: "https://likumi.lv/ta/en/en/id/253451-value-added-tax-law",
			},
			{
				Title: i18n.String{
					i18n.EN: "EU Council Directive on the common system of value added tax",
				},
				URL: "https://eur-lex.europa.eu/legal-content/EN/TXT/HTML/?uri=CELEX:32006L0112",
			},
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			// STANDARD TAX RATE - GENERAL (Minimum 15%)
			//Applied to most goods and services in Latvia.
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General Rate",
					i18n.LV: "Nodokļa Standartlikmi",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2013, 1, 1),
						Percent: num.MakePercentage(210, 3),
					},
				},
			},
			//STANDARD TAX RATE - HIGHEST REDUCED RATE
			//Applied to essencial goods and services (medicines, food suplies, etc).
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateIntermediate,
				Name: i18n.String{
					i18n.EN: "First Reduced Rate",
					i18n.LV: "Pirmā samazinātā likme",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2013, 1, 1),
						Percent: num.MakePercentage(120, 3),
					},
				},
			},
			//STANDARD TAX RATE - LOWEST REDUCED RATE
			//Applied to book suplies, press and media publications and related anexes.
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Second Reduced Rate",
					i18n.LV: "Otrā samazinātā likme",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2018, 1, 1),
						Percent: num.MakePercentage(50, 3),
					},
				},
			},
			//NO SUPER-REDUCED TAX RATE (Maximum 5%)
			//ZERO TAX RATE FOR EXPORT OR INTRACOMMUNITY.
		},
	},
}
