package at

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.Category{
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
		Sources: []*tax.Source{
			{
				Title: i18n.String{
					i18n.EN: "Business Service Portal - Rates of VAT",
				},
				URL: "https://www.usp.gv.at/en/steuern-finanzen/umsatzsteuer/steuersaetze-der-umsatzsteuer.html",
			},
		},
		Retained: false,
		Rates: []*tax.Rate{
			{
				Key: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Key: tax.RateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
				},
				Values: []*tax.RateValue{
					{
						Since:   cal.NewDate(2011, 1, 4),
						Percent: num.MakePercentage(200, 3),
					},
				},
			},
			{
				Key: tax.RateIntermediate,
				Name: i18n.String{
					i18n.EN: "Intermediate Rate",
				},
				Values: []*tax.RateValue{
					{
						Since:   cal.NewDate(2011, 1, 1),
						Percent: num.MakePercentage(130, 3),
					},
				},
			},
			{
				Key: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
				},
				Values: []*tax.RateValue{
					{
						Since:   cal.NewDate(2011, 1, 4),
						Percent: num.MakePercentage(100, 3),
					},
				},
			},
		},
	},
}
