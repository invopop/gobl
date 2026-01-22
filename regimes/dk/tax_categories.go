package dk

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
			i18n.DA: "moms",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.DA: "Meromsætningsafgift",
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General Rate",
					i18n.DA: "Standardsats",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(250, 3),
						Since:   cal.NewDate(1992, 1, 1),
					},
				},
			},
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Danish Tax Agency - VAT rates",
				},
				URL: "https://skat.dk/",
			},
		},
	},
}
