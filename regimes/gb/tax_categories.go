package gb

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
				Key: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
				},
				Values: []*tax.RateValue{
					{
						Since:   cal.NewDate(2011, 1, 4),
						Percent: num.MakePercentage(50, 3),
					},
				},
			},
		},
	},
}
