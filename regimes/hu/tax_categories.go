package hu

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
			i18n.HU: "ÁFA",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.HU: "Általános forgalmi adó",
		},
		Keys: tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			// Standard rate
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.HU: "Általános kulcs",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(270, 3),
						Since:   cal.NewDate(2012, 1, 1),
					},
					{
						Percent: num.MakePercentage(250, 3),
						Since:   cal.NewDate(2009, 7, 1),
					},
					{
						Percent: num.MakePercentage(200, 3),
						Since:   cal.NewDate(2006, 1, 1),
					},
					{
						Percent: num.MakePercentage(250, 3),
						Since:   cal.NewDate(2004, 5, 1),
					},
				},
			},
			// Intermediate reduced rate
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.HU: "Kedvezményes kulcs",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(180, 3),
						Since:   cal.NewDate(2009, 7, 1),
					},
				},
			},
			// Lower reduced rate
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Super Reduced Rate",
					i18n.HU: "Kedvezményes alsó kulcs",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(50, 3),
						Since:   cal.NewDate(2004, 5, 1),
					},
				},
			},
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "NAV - ÁFA rates and tax-exempt activities",
				},
				URL: "https://nav.gov.hu/ugyfeliranytu/adokulcsok_jarulekmertekek/afakulcs_adomen",
			},
			{
				Title: i18n.String{
					i18n.EN: "Act CXXVII of 2007 on Value Added Tax",
				},
				URL: "https://njt.hu/jogszabaly/2007-127-00-00",
			},
		},
	},
}
