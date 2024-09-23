package au

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"

	"time"
)

var taxCategories = []*tax.CategoryDef{
	{
		Code: tax.CategoryGST,
		Name: i18n.String{
			i18n.EN: "GST",
		},
		Title: i18n.String{
			i18n.EN: "Goods and Services Tax",
		},
		Retained: false,
		Rates: []*tax.RateDef{
			{
				Key: tax.RateExempt,
				Name: i18n.String{
					i18n.EN: "GST-free",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2000, time.July, 1),
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Key: tax.RateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2000, time.July, 1),
						Percent: num.MakePercentage(100, 3),
					},
				},
			},
			{
				Key: tax.RateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2000, time.July, 1),
						Percent: num.MakePercentage(100, 3),
					},
				},
			},
		},
	},
}
