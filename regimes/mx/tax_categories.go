package mx

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

// Tax rates specific to Mexico.
const (
	TaxRateExempt cbc.Key = "exempt"
)

var taxCategories = []*tax.Category{
	{
		Code: common.TaxCategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.ES: "IVA",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.ES: "Impuesto al Valor Agregado",
		},
		Retained: false,
		Rates: []*tax.Rate{
			{
				Key: common.TaxRateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.ES: "Tasa General",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(160, 3),
					},
				},
			},
			{
				Key: common.TaxRateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced (Border) Rate",
					i18n.ES: "Tasa Reducida (Fronteriza)",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(80, 3),
					},
				},
			},
			{
				Key: common.TaxRateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
					i18n.ES: "Tasa Cero",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Key: TaxRateExempt,
				Name: i18n.String{
					i18n.EN: "Exempt",
					i18n.ES: "Exenta",
				},
				Exempt: true,
			},
		},
	},
}
