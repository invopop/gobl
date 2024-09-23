package mx

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// Local tax categories
const (
	TaxCategoryRVAT  cbc.Code = "RVAT"  // IVA (Retenido)
	TaxCategoryIEPS  cbc.Code = "IEPS"  // Impuesto Especial sobre Producción y Servicios
	TaxCategoryRIEPS cbc.Code = "RIEPS" // Impuesto Especial sobre Producción y Servicios (Retenido)
	TaxCategoryISR   cbc.Code = "ISR"   // Impuesto Sobre la Renta
)

// Local tax rates
const (
	TaxRateExempt cbc.Key = "exempt"
)

var taxCategories = []*tax.CategoryDef{
	//
	// IVA
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
				Key: tax.RateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.ES: "Tasa General",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(160, 3),
					},
				},
			},
			{
				Key: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced (Border) Rate",
					i18n.ES: "Tasa Reducida (Fronteriza)",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(80, 3),
					},
				},
			},
			{
				Key: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
					i18n.ES: "Tasa Cero",
				},
				Values: []*tax.RateValueDef{
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

	//
	// IVA (Retenido)
	//
	{
		Code: TaxCategoryRVAT,
		Name: i18n.String{
			i18n.EN: "Retained VAT",
			i18n.ES: "IVA Retenido",
		},
		Title: i18n.String{
			i18n.EN: "Retained Value Added Tax",
			i18n.ES: "Impuesto al Valor Agregado Retenido",
		},
		Retained: true,
		Rates:    []*tax.RateDef{},
	},

	//
	// IEPS
	//
	{
		Code: TaxCategoryIEPS,
		Name: i18n.String{
			i18n.EN: "IEPS",
			i18n.ES: "IEPS",
		},
		Title: i18n.String{
			i18n.EN: "Special Tax on Production and Services",
			i18n.ES: "Impuesto Especial sobre Producción y Servicios",
		},
		Retained: false,
		Rates:    []*tax.RateDef{},
	},

	//
	// IEPS (Retenido)
	//
	{
		Code: TaxCategoryRIEPS,
		Name: i18n.String{
			i18n.EN: "Retained IEPS",
			i18n.ES: "IEPS Retenido",
		},
		Title: i18n.String{
			i18n.EN: "Retained Special Tax on Production and Services",
			i18n.ES: "Impuesto Especial sobre Producción y Servicios Retenido",
		},
		Retained: true,
		Rates:    []*tax.RateDef{},
	},

	//
	// ISR
	//
	{
		Code: TaxCategoryISR,
		Name: i18n.String{
			i18n.EN: "ISR",
			i18n.ES: "ISR",
		},
		Title: i18n.String{
			i18n.EN: "Income Tax",
			i18n.ES: "Impuesto Sobre la Renta",
		},
		Retained: true,
		Rates:    []*tax.RateDef{},
	},
}
