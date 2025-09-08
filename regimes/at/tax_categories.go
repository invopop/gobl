package at

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
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Business Service Portal - Rates of VAT",
				},
				URL: "https://www.usp.gv.at/en/steuern-finanzen/umsatzsteuer/steuersaetze-der-umsatzsteuer.html",
			},
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General Rate",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2011, 1, 4),
						Percent: num.MakePercentage(200, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateIntermediate,
				Name: i18n.String{
					i18n.EN: "Intermediate Rate",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(130, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2011, 1, 4),
						Percent: num.MakePercentage(100, 3),
					},
				},
			},
		},
	},
}
