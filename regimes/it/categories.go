package it

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

// Local tax category definitions which are not considered standard.
const (
	TaxCategoryRA   cbc.Code = "RA" // https://www.agenziaentrate.gov.it/portale/imposta-sul-reddito-delle-persone-fisiche-irpef-/aliquote-e-calcolo-dell-irpef
	TaxCategoryINPS cbc.Code = "INPS"
)

var categories = []*tax.Category{
	{
		Code:     common.TaxCategoryVAT,
		Retained: false,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.IT: "IVA",
		},
		Desc: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.IT: "Imposta sul Valore Aggiunto",
		},
		Rates: []*tax.Rate{
			{
				Key: common.TaxRateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
					i18n.IT: "Aliquota Zero",
				},
				Tags: vatZeroTaxTags,
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Key: common.TaxRateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Minimum Rate",
					i18n.IT: "Aliquota Minima",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(40, 3),
					},
				},
			},
			{
				Key: common.TaxRateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.IT: "Aliquota Ridotta",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(50, 3),
					},
				},
			},
			{
				Key: common.TaxRateIntermediate,
				Name: i18n.String{
					i18n.EN: "Intermediate Rate",
					i18n.IT: "Aliquota Intermedia",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(100, 3),
					},
				},
			},
			{
				Key: common.TaxRateStandard,
				Name: i18n.String{
					i18n.EN: "Ordinary Rate",
					i18n.IT: "Aliquota Ordinaria",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(220, 3),
					},
				},
			},
		},
	},
	{
		Code:     TaxCategoryRA,
		Retained: true,
		Name: i18n.String{
			i18n.EN: "RA",
			i18n.IT: "RA",
		},
		Desc: i18n.String{
			i18n.EN: "Withholding Tax",
			i18n.IT: "Ritenuta d'Acconto",
		},
		Rates: []*tax.Rate{
			{
				Key: common.TaxRateZero,
				Name: i18n.String{
					i18n.EN: "Ordinary Rate",
					i18n.IT: "Aliquota Ordinaria",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(200, 3),
					},
				},
			},
		},
	},
}
