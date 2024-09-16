package co

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// Local tax categories.
const (
	TaxCategoryIC        cbc.Code = "IC"  // Impuesto Consumo
	TaxCategoryICA       cbc.Code = "ICA" // Impuesto de Industria y Comercio
	TaxCategoryINC       cbc.Code = "INC"
	TaxCategoryReteIVA   cbc.Code = "RVAT" // ReteIVA
	TaxCategoryReteRenta cbc.Code = "RR"   // ReteRenta
	TaxCategoryReteICA   cbc.Code = "RICA" // ReteICA
)

var taxCategories = []*tax.CategoryDef{
	//
	// VAT
	//
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.ES: "IVA",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.ES: "Impuesto al Valor Agregado",
		},
		Retained: false,
		Rates: []*tax.RateDef{
			{
				Key: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
					i18n.ES: "Cero",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Key: tax.RateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.ES: "Estándar",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2017, 1, 1),
						Percent: num.MakePercentage(190, 3),
					},
					{
						Since:   cal.NewDate(2006, 1, 1),
						Percent: num.MakePercentage(160, 3),
					},
				},
			},
			{
				Key: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.ES: "Reducido",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2006, 1, 1),
						Percent: num.MakePercentage(50, 3),
					},
				},
			},
		},
	},
	//
	// IC - national
	//
	{
		Code: TaxCategoryIC,
		Name: i18n.String{
			i18n.ES: "IC",
		},
		Title: i18n.String{
			i18n.EN: "Consumption Tax",
			i18n.ES: "Impuesto sobre Consumo",
		},
		Retained: false,
		Rates:    []*tax.RateDef{},
	},
	//
	// ICA - local taxes
	//
	{
		Code: TaxCategoryICA,
		Name: i18n.String{
			i18n.ES: "ICA",
		},
		Title: i18n.String{
			i18n.EN: "Industry and Commerce Tax",
			i18n.ES: "Impuesto de Industria y Comercio",
		},
		Retained: false,
		Rates:    []*tax.RateDef{},
	},
	//
	// ReteIVA
	//
	{
		Code: TaxCategoryReteIVA,
		Name: i18n.String{
			i18n.ES: "ReteIVA",
		},
		Title: i18n.String{
			i18n.ES: "Retención en la fuente por el Impuesto al Valor Agregado",
		},
		Retained: true,
		Rates:    []*tax.RateDef{},
	},
	//
	// ReteICA
	//
	{
		Code: TaxCategoryReteICA,
		Name: i18n.String{
			i18n.ES: "ReteICA",
		},
		Title: i18n.String{
			i18n.ES: "Retención en la fuente por el Impuesto de Industria y Comercio",
		},
		Retained: true,
		Rates:    []*tax.RateDef{},
	},
	//
	// ReteRenta
	//
	{
		Code: TaxCategoryReteRenta,
		Name: i18n.String{
			i18n.ES: "Retefuente",
		},
		Title: i18n.String{
			i18n.ES: "Retención en la fuente por el Impuesto de la Renta",
		},
		Retained: true,
		Rates:    []*tax.RateDef{},
	},
}
