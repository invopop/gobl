// Package il defines VAT tax categories specific to Israel.
package il

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
			i18n.HE: "מע\"מ",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.HE: "מס ערך מוסף",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Israel Tax Authority - VAT Law 1975",
					i18n.HE: "רשות המסים בישראל - חוק מס ערך מוסף תשל\"ו-1975",
				},
				URL: "https://www.gov.il/en/departments/israel_tax_authority",
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
					i18n.HE: "שיעור כללי",
				},
				Description: i18n.String{
					i18n.EN: "Applies to most goods and services unless specified otherwise.",
					i18n.HE: "חל על רוב הסחורות והשירותים אלא אם צוין אחרת.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2025, 1, 1),
						Percent: num.MakePercentage(18, 2),
					},
					{
						Since:   cal.NewDate(2015, 10, 1),
						Percent: num.MakePercentage(17, 2),
					},
					{
						Since:   cal.NewDate(2013, 6, 2),
						Percent: num.MakePercentage(18, 2),
					},
					{
						Since:   cal.NewDate(2012, 9, 1),
						Percent: num.MakePercentage(17, 2),
					},
					{
						Since:   cal.NewDate(2010, 1, 1),
						Percent: num.MakePercentage(16, 2),
					},
					{
						Since:   cal.NewDate(2009, 7, 1),
						Percent: num.MakePercentage(165, 3),
					},
					{
						Since:   cal.NewDate(2006, 7, 1),
						Percent: num.MakePercentage(155, 3),
					},
					{
						Since:   cal.NewDate(2005, 9, 1),
						Percent: num.MakePercentage(165, 3),
					},
					{
						Since:   cal.NewDate(2004, 3, 1),
						Percent: num.MakePercentage(17, 2),
					},
				},
			},
		},
	},
}
