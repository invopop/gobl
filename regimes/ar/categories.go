package ar

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// TaxRateIncreased is the key for Argentina’s increased VAT rate.
const (
	TaxRateIncreased cbc.Key = "increased"
)

var categories = []*tax.CategoryDef{
	{
		Code:     tax.CategoryVAT,
		Retained: false,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.ES: "IVA",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.ES: "Impuesto al Valor Agregado",
		},
		Rates: []*tax.RateDef{
			{
				Key: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
					i18n.ES: "Alícuota Cero",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(0, 2),
					},
				},
			},
			{
				Key: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.ES: "Alícuota Reducida",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(105, 3),
					},
				},
			},
			{
				Key: tax.RateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.ES: "Alícuota General",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(21, 2),
					},
				},
			},
			{
				Key: TaxRateIncreased,
				Name: i18n.String{
					i18n.EN: "Increased Rate",
					i18n.ES: "Alícuota Incrementada",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(27, 2),
					},
				},
			},
			{
				Key: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Super Reduced Rate",
					i18n.ES: "Alícuota Super Reducida",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(5, 2),
					},
				},
			},
			{
				Key: tax.RateExempt,
				Name: i18n.String{
					i18n.EN: "Exempt",
					i18n.ES: "Exento",
				},
			},
			{
				Key: tax.RateSpecial,
				Name: i18n.String{
					i18n.EN: "Special Super Reduced Rate",
					i18n.ES: "Alícuota Especial Super Reducida",
				},
				Values: []*tax.RateValueDef{
					{Percent: num.MakePercentage(25, 3)},
				},
			},
		},
	},
}
