package hu

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.Category{
	// VAT
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
		Extensions: []cbc.Key{
			ExtKeyExemptionCode,
		},
		Rates: []*tax.Rate{
			{
				Key: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
					i18n.HU: "ÁFA-mentes",
				},
				Values: []*tax.RateValue{
					{
						Since:   cal.NewDate(2024, 1, 1),
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Key: tax.RateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.HU: "Általános",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(27, 3),
					},
				},
			},
			{
				Key: tax.RateIntermediate,
				Name: i18n.String{
					i18n.EN: "Intermediate Rate",
					i18n.HU: "Köztes",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(18, 3),
					},
				},
			},
			{
				Key: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.HU: "Csökkentett",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(5, 3),
					},
				},
			},
		},
	},
}
