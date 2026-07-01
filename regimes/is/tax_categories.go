package is

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
			i18n.IS: "VSK",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.IS: "Virðisaukaskattur",
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.IS: "Almennt skatthlutfall",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(240, 3),
						Since:   cal.NewDate(2015, 1, 1),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.IS: "Lægra skatthlutfall",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(110, 3),
						Since:   cal.NewDate(2015, 1, 1),
					},
				},
			},
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Skatturinn - Value Added Tax (VSK)",
				},
				URL: "https://www.skatturinn.is/english/companies/value-added-tax/",
				At:  cal.NewDateTime(2026, 6, 30, 0, 0, 0),
			},
			{
				Title: i18n.String{
					i18n.EN: "KPMG - Icelandic Tax Facts 2025",
				},
				URL: "https://assets.kpmg.com/content/dam/kpmg/is/pdf/2025/01/Icelandic-Tax-Facts-2025.pdf",
				At:  cal.NewDateTime(2026, 6, 30, 0, 0, 0),
			},
		},
	},
}
