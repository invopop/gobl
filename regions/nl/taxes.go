package es

import (
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regions/common"
	"github.com/invopop/gobl/tax"
)

// Local tax category definitions which are not considered standard.
const (
	TaxCategoryVATEqualizationSurcharge tax.Code = "VATEQS"
	TaxCategoryIRPF                     tax.Code = "IRPF"
	TaxCategoryIGIC                     tax.Code = "IGIC"
	TaxCategoryIPSI                     tax.Code = "IPSI"
)

var taxRegion = tax.Region{
	Code: "ES",
	Name: i18n.String{
		i18n.EN: "The Netherlands",
		i18n.NL: "Nederland",
	},
	Categories: []tax.Category{
		//
		// VAT
		//
		{
			Code: common.TaxCategoryVAT,
			Name: i18n.String{
				i18n.EN: "VAT",
				i18n.NL: "BTW",
			},
			Desc: i18n.String{
				i18n.EN: "Value Added Tax",
				i18n.NL: "Belasting Toegevoegde Waarde",
			},
			Retained: false,
			Defs: []tax.Def{
				{
					Code: common.TaxRateVATZero,
					Name: i18n.String{
						i18n.EN: "VAT Zero Rate",
						i18n.NL: `BTW 0%-tarief`,
					},
					Values: []tax.Value{
						{
							Percent: num.MakePercentage(0, 3),
						},
					},
				},
				{
					Code: common.TaxRateVATStandard,
					Name: i18n.String{
						i18n.EN: "VAT Standard Rate",
						i18n.NL: "BTW Standaardtarief",
					},
					Values: []tax.Value{
						{
							Since:   org.NewDate(2012, 9, 1),
							Percent: num.MakePercentage(210, 3),
						},
						{
							Since:   org.NewDate(2010, 7, 1),
							Percent: num.MakePercentage(180, 3),
						},
						{
							Since:   org.NewDate(1995, 1, 1),
							Percent: num.MakePercentage(160, 3),
						},
						{
							Since:   org.NewDate(1993, 1, 1),
							Percent: num.MakePercentage(150, 3),
						},
					},
				},
				{
					Code: common.TaxRateVATReduced,
					Name: i18n.String{
						i18n.EN: "VAT Reduced Rate",
						i18n.NL: "BTW Gereduceerd Tarief",
					},
					Values: []tax.Value{
						{
							Since:   org.NewDate(1900, 1, 1),
							Percent: num.MakePercentage(100, 9),
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
					Code: common.TaxRateVATZero,
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
					Code: common.TaxRateVATStandard,
					Name: i18n.String{
						i18n.EN: "VAT Standard Rate",
						i18n.ES: "IVA Tipo General",
					},
					Values: []tax.Value{
						{
							Since:   org.NewDate(2012, 9, 1),
							Percent: num.MakePercentage(520, 4),
						},
						{
							Since:   org.NewDate(1993, 1, 1),
							Percent: num.MakePercentage(400, 4),
						},
					},
				},
				{
					Code: common.TaxRateVATReduced,
					Name: i18n.String{
						i18n.EN: "VAT Reduced Rate",
						i18n.ES: "IVA Tipo Reducido",
					},
					Values: []tax.Value{
						{
							Since:   org.NewDate(2012, 9, 1),
							Percent: num.MakePercentage(140, 4),
						},
						{
							Since:   org.NewDate(1993, 1, 1),
							Percent: num.MakePercentage(100, 4),
						},
					},
				},
				{
					Code: common.TaxRateVATSuperReduced,
					Name: i18n.String{
						i18n.EN: "VAT Super-Reduced Rate",
						i18n.ES: "IVA Tipo Superreducido",
					},
					Values: []tax.Value{
						{
							Since:   org.NewDate(1993, 1, 1),
							Percent: num.MakePercentage(50, 4),
						},
					},
				},
			},
		},
	},
}
