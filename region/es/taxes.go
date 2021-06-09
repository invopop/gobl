package es

import (
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/region/eu"
	"github.com/invopop/gobl/tax"
)

// Local tax category definitions which are not considered standard.
const (
	TaxCategoryVATEqualizationSurcharge tax.Code = "VATEQS"
	TaxCategoryIRPF                     tax.Code = "IRPF"
	TaxCategoryIGIC                     tax.Code = "IGIC"
	TaxCategoryIPSI                     tax.Code = "IPSI"
)

// Specific tax rate codes.
const (
	// VAT non-standard Rates
	TaxRateVATTobacco tax.Code = "tobacco"

	// IRPF "Autonomo" Rates
	TaxRateIRPFStandard     tax.Code = "STD"  // Professional or artistic
	TaxRateIRPFFirst        tax.Code = "1ST"  // First 2 years
	TaxRateIRPFModules      tax.Code = "MOD"  // Module system
	TaxRateIRPFAgriculture  tax.Code = "AGR"  // Agricultural
	TaxRateIRPFAgriculture2 tax.Code = "AGR2" // Agricultural special
)

var taxRegion = tax.Region{
	Code: "es",
	Name: i18n.String{
		i18n.EN: "Spain",
		i18n.ES: "España",
	},
	Categories: []tax.Category{
		//
		// VAT
		//
		{
			Code: eu.TaxCategoryVAT,
			Name: i18n.String{
				i18n.EN: "VAT",
				i18n.ES: "IVA",
			},
			Desc: i18n.String{
				i18n.EN: "Value Added Tax",
				i18n.ES: "Impuesto sobre el Valor Añadido",
			},
			Retained: false,
			Defs: []tax.Def{
				{
					Code: eu.TaxRateVATZero,
					Name: i18n.String{
						i18n.EN: "VAT Zero Rate",
						i18n.ES: "IVA Tipo Zero",
					},
					Values: []tax.Value{
						{
							Percent: num.MakePercentage(0, 3),
						},
					},
				},
				{
					Code: eu.TaxRateVATStandard,
					Name: i18n.String{
						i18n.EN: "VAT Standard Rate",
						i18n.ES: "IVA Tipo General",
					},
					Values: []tax.Value{
						{
							Since:   org.MakeDate(2012, 9, 1),
							Percent: num.MakePercentage(210, 3),
						},
						{
							Since:   org.MakeDate(2010, 7, 1),
							Percent: num.MakePercentage(180, 3),
						},
						{
							Since:   org.MakeDate(1995, 1, 1),
							Percent: num.MakePercentage(160, 3),
						},
						{
							Since:   org.MakeDate(1993, 1, 1),
							Percent: num.MakePercentage(150, 3),
						},
					},
				},
				{
					Code: eu.TaxRateVATReduced,
					Name: i18n.String{
						i18n.EN: "VAT Reduced Rate",
						i18n.ES: "IVA Tipo Reducido",
					},
					Values: []tax.Value{
						{
							Since:   org.MakeDate(2012, 9, 1),
							Percent: num.MakePercentage(100, 3),
						},
						{
							Since:   org.MakeDate(2010, 7, 1),
							Percent: num.MakePercentage(80, 3),
						},
						{
							Since:   org.MakeDate(1995, 1, 1),
							Percent: num.MakePercentage(70, 3),
						},
						{
							Since:   org.MakeDate(1993, 1, 1),
							Percent: num.MakePercentage(60, 3),
						},
					},
				},
				{
					Code: eu.TaxRateVATSuperReduced,
					Name: i18n.String{
						i18n.EN: "VAT Super-Reduced Rate",
						i18n.ES: "IVA Tipo Superreducido",
					},
					Values: []tax.Value{
						{
							Since:   org.MakeDate(1995, 1, 1),
							Percent: num.MakePercentage(40, 3),
						},
						{
							Since:   org.MakeDate(1993, 1, 1),
							Percent: num.MakePercentage(30, 3),
						},
					},
				},
			},
		},
		//
		// VAT Equalization Surcharge (Recargo de equivalencia)
		//
		{
			Code: TaxCategoryVATEqualizationSurcharge,
			Name: i18n.String{
				i18n.EN: "VAT Equalization Surcharge",
				i18n.ES: "IVA Recargo de Equivalencia",
			},
			Retained: false,
			Defs: []tax.Def{
				{
					Code: eu.TaxRateVATZero,
					Name: i18n.String{
						i18n.EN: "VAT Zero Rate",
						i18n.ES: "IVA Tipo Exento",
					},
					Values: []tax.Value{
						{
							Percent: num.MakePercentage(0, 3),
						},
					},
				},
				{
					Code: eu.TaxRateVATStandard,
					Name: i18n.String{
						i18n.EN: "VAT Standard Rate",
						i18n.ES: "IVA Tipo General",
					},
					Values: []tax.Value{
						{
							Since:   org.MakeDate(2012, 9, 1),
							Percent: num.MakePercentage(520, 4),
						},
						{
							Since:   org.MakeDate(1993, 1, 1),
							Percent: num.MakePercentage(400, 4),
						},
					},
				},
				{
					Code: eu.TaxRateVATReduced,
					Name: i18n.String{
						i18n.EN: "VAT Reduced Rate",
						i18n.ES: "IVA Tipo Reducido",
					},
					Values: []tax.Value{
						{
							Since:   org.MakeDate(2012, 9, 1),
							Percent: num.MakePercentage(140, 4),
						},
						{
							Since:   org.MakeDate(1993, 1, 1),
							Percent: num.MakePercentage(100, 4),
						},
					},
				},
				{
					Code: eu.TaxRateVATSuperReduced,
					Name: i18n.String{
						i18n.EN: "VAT Super-Reduced Rate",
						i18n.ES: "IVA Tipo Superreducido",
					},
					Values: []tax.Value{
						{
							Since:   org.MakeDate(1993, 1, 1),
							Percent: num.MakePercentage(50, 4),
						},
					},
				},
				{
					Code: TaxRateVATTobacco,
					Name: i18n.String{
						i18n.EN: "VAT Tobacco Rate",
						i18n.ES: "IVA Tipo Tobaco",
					},
					Values: []tax.Value{
						{
							Since:   org.MakeDate(2007, 1, 1),
							Percent: num.MakePercentage(75, 4),
						},
					},
				},
			},
		},
		//
		// IRPF
		//
		{
			Code:     TaxCategoryIRPF,
			Retained: true,
			Defs: []tax.Def{
				{
					Code: TaxRateIRPFStandard,
					Name: i18n.String{
						i18n.EN: "IRPF Standard Rate",
						i18n.ES: "IRPF Tipo General",
					},
					Values: []tax.Value{
						{
							Since:   org.MakeDate(2015, 7, 12),
							Percent: num.MakePercentage(150, 3),
						},
						{
							Since:   org.MakeDate(2015, 1, 1),
							Percent: num.MakePercentage(190, 3),
						},
						{
							Since:   org.MakeDate(2012, 9, 1),
							Percent: num.MakePercentage(210, 3),
						},
						{
							Since:   org.MakeDate(2007, 1, 1),
							Percent: num.MakePercentage(150, 3),
						},
					},
				},
				{
					Code: TaxRateIRPFFirst,
					Name: i18n.String{
						i18n.EN: "IRPF Starting Rate",
						i18n.ES: "IRPF Tipo Inicial",
					},
					Values: []tax.Value{
						{
							Since:   org.MakeDate(2007, 1, 1),
							Percent: num.MakePercentage(70, 3),
						},
					},
				},
				{
					Code: TaxRateIRPFModules,
					Name: i18n.String{
						i18n.EN: "IRPF Modules Rate",
						i18n.ES: "IRPF Tipo Modulos",
					},
					Values: []tax.Value{
						{
							Since:   org.MakeDate(2007, 1, 1),
							Percent: num.MakePercentage(10, 3),
						},
					},
				},
			},
		},
	},
}
