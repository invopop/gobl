package fi

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
			i18n.FI: "ALV",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.FI: "Arvonlisävero",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Verohallinto - VAT Rates"),
				URL:   "https://www.vero.fi/en/businesses-and-corporations/taxes-and-charges/vat/rates-of-vat/",
			},
			{
				Title: i18n.NewString("Finlex - Arvonlisäverolaki 1501/1993"),
				URL:   "https://www.finlex.fi/fi/laki/ajantasa/1993/19931501",
			},
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.FI: "Yleinen verokanta",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2024, 9, 1),
						Percent: num.MakePercentage(255, 3), // 25.5%
					},
					{
						Since:   cal.NewDate(2013, 1, 1),
						Percent: num.MakePercentage(240, 3), // 24.0%
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.FI: "Ensimmäinen alennettu verokanta",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2026, 1, 1),
						Percent: num.MakePercentage(135, 3), // 13.5%
					},
					{
						Since:   cal.NewDate(2013, 1, 1),
						Percent: num.MakePercentage(140, 3), // 14.0%
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Super-Reduced Rate",
					i18n.FI: "Toinen alennettu verokanta",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2013, 1, 1),
						Percent: num.MakePercentage(100, 3), // 10.0%
					},
				},
			},
		},
	},
}
