package ee

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.CategoryDef{
	//
	// VAT (Käibemaks)
	//
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.ET: "Käibemaks",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.ET: "Käibemaks",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Estonian Tax and Customs Board - VAT rates",
				},
				URL: "https://www.emta.ee/en/business-client/taxes-and-payment/value-added-tax/vat-rates-and-supply-exempt-tax/standard-vat-rate",
			},
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Rate: tax.RateGeneral,
				Keys: []cbc.Key{tax.KeyStandard},
				Name: i18n.String{
					i18n.EN: "Standard Rate",
				},
				// Values are kept in descending date order, as required by GOBL.
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2025, 7, 1),
						Percent: num.MakePercentage(240, 3),
					},
					{
						Since:   cal.NewDate(2024, 1, 1),
						Percent: num.MakePercentage(220, 3),
					},
					{
						Since:   cal.NewDate(2009, 7, 1),
						Percent: num.MakePercentage(200, 3),
					},
				},
			},
			{
				// Higher reduced rate, currently applied to accommodation services.
				Rate: tax.RateReduced,
				Keys: []cbc.Key{tax.KeyStandard},
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(130, 3),
					},
				},
			},
			{
				// Lower reduced rate for books, press publications and certain medicines.
				Rate: tax.RateSuperReduced,
				Keys: []cbc.Key{tax.KeyStandard},
				Name: i18n.String{
					i18n.EN: "Super-Reduced Rate",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(90, 3),
					},
				},
			},
			{
				Rate: tax.RateZero,
				Keys: []cbc.Key{tax.KeyZero},
				Name: i18n.String{
					i18n.EN: "Zero Rate",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
		},
	},
}
