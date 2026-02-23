package hr

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
			i18n.HR: "PDV",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.HR: "Porez na dodanu vrijednost",
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General Rate",
					i18n.HR: "Opća stopa",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(250, 3),
						Since:   cal.NewDate(2013, 7, 17),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.HR: "Snižena stopa",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(130, 3),
						Since:   cal.NewDate(2013, 7, 17),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Super Reduced Rate",
					i18n.HR: "Super snižena stopa",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(50, 3),
						Since:   cal.NewDate(2013, 7, 17),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
					i18n.HR: "Nulta stopa",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(0, 3),
						Since:   cal.NewDate(2013, 7, 17),
					},
				},
			},
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Croatian Tax Administration - VAT",
					i18n.HR: "Porezna uprava - Porez na dodanu vrijednost",
				},
				URL: "https://porezna-uprava.gov.hr/en/vat/7362",
			},
		},
	},
}
