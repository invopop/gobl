package es

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.Category{
	//
	// VAT
	//
	{
		Code:     common.TaxCategoryVAT,
		Retained: false,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.ES: "IVA",
		},
		Desc: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.ES: "Impuesto sobre el Valor Añadido",
		},
		Map: cbc.CodeMap{
			KeyFacturaETaxTypeCode: "01",
		},
		Rates: []*tax.Rate{
			{
				Key: common.TaxRateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
					i18n.ES: "Tipo Cero",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Key: common.TaxRateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.ES: "Tipo General",
				},
				Values: []*tax.RateValue{
					{
						Since:   cal.NewDate(2012, 9, 1),
						Percent: num.MakePercentage(210, 3),
					},
					{
						Since:   cal.NewDate(2010, 7, 1),
						Percent: num.MakePercentage(180, 3),
					},
					{
						Since:   cal.NewDate(1995, 1, 1),
						Percent: num.MakePercentage(160, 3),
					},
					{
						Since:   cal.NewDate(1993, 1, 1),
						Percent: num.MakePercentage(150, 3),
					},
				},
			},
			{
				Key: common.TaxRateStandard.With(TaxRateEquivalence),
				Name: i18n.String{
					i18n.EN: "Standard Rate + Equivalence Surcharge",
					i18n.ES: "Tipo General + Recargo de Equivalencia",
				},
				Values: []*tax.RateValue{
					{
						Since:     cal.NewDate(2012, 9, 1),
						Percent:   num.MakePercentage(210, 3),
						Surcharge: num.NewPercentage(52, 3),
					},
					{
						Since:     cal.NewDate(2010, 7, 1),
						Percent:   num.MakePercentage(180, 3),
						Surcharge: num.NewPercentage(40, 3),
					},
				},
			},
			{
				Key: common.TaxRateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.ES: "Tipo Reducido",
				},
				Values: []*tax.RateValue{
					{
						Since:   cal.NewDate(2012, 9, 1),
						Percent: num.MakePercentage(100, 3),
					},
					{
						Since:   cal.NewDate(2010, 7, 1),
						Percent: num.MakePercentage(80, 3),
					},
					{
						Since:   cal.NewDate(1995, 1, 1),
						Percent: num.MakePercentage(70, 3),
					},
					{
						Since:   cal.NewDate(1993, 1, 1),
						Percent: num.MakePercentage(60, 3),
					},
				},
			},
			{
				Key: common.TaxRateReduced.With(TaxRateEquivalence),
				Name: i18n.String{
					i18n.EN: "Reduced Rate + Equivalence Surcharge",
					i18n.ES: "Tipo Reducido + Recargo de Equivalencia",
				},
				Values: []*tax.RateValue{
					{
						Since:     cal.NewDate(2012, 9, 1),
						Percent:   num.MakePercentage(100, 3),
						Surcharge: num.NewPercentage(14, 3),
					},
					{
						Since:     cal.NewDate(2010, 7, 1),
						Percent:   num.MakePercentage(80, 3),
						Surcharge: num.NewPercentage(10, 3),
					},
				},
			},
			{
				Key: common.TaxRateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Super-Reduced Rate",
					i18n.ES: "Tipo Superreducido",
				},
				Values: []*tax.RateValue{
					{
						Since:   cal.NewDate(1995, 1, 1),
						Percent: num.MakePercentage(40, 3),
					},
					{
						Since:   cal.NewDate(1993, 1, 1),
						Percent: num.MakePercentage(30, 3),
					},
				},
			},
			{
				Key: common.TaxRateSuperReduced.With(TaxRateEquivalence),
				Name: i18n.String{
					i18n.EN: "Super-Reduced Rate + Equivalence Surcharge",
					i18n.ES: "Tipo Superreducido + Recargo de Equivalencia",
				},
				Values: []*tax.RateValue{
					{
						Since:     cal.NewDate(1995, 1, 1),
						Percent:   num.MakePercentage(40, 3),
						Surcharge: num.NewPercentage(5, 3),
					},
				},
			},
			{
				Key:    common.TaxRateExempt,
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Exempt",
					i18n.ES: "Exenta",
				},
				Extensions: []cbc.Key{ExtKeyTBAIExemption},
			},
		},
	},

	//
	// IGIC
	//
	{
		Code:     TaxCategoryIGIC,
		Retained: false,
		Name: i18n.String{
			i18n.EN: "IGIC",
			i18n.ES: "IGIC",
		},
		Map: cbc.CodeMap{
			KeyFacturaETaxTypeCode: "03",
		},
		Desc: i18n.String{
			i18n.EN: "Canary Island General Indirect Tax",
			i18n.ES: "Impuesto General Indirecto Canario",
		},
		// This is a subset of the possible rates.
		Rates: []*tax.Rate{
			{
				Key: common.TaxRateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
					i18n.ES: "Tipo Cero",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Key: common.TaxRateStandard,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.ES: "Tipo General",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(70, 3),
					},
				},
			},
			{
				Key: common.TaxRateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.ES: "Tipo Reducido",
				},
				Values: []*tax.RateValue{
					{
						Percent: num.MakePercentage(30, 3),
					},
				},
			},
		},
	},

	//
	// IPSI
	//
	{
		Code:     TaxCategoryIPSI,
		Retained: false,
		Name: i18n.String{
			i18n.EN: "IPSI",
			i18n.ES: "IPSI",
		},
		Map: cbc.CodeMap{
			KeyFacturaETaxTypeCode: "02",
		},
		Desc: i18n.String{
			i18n.EN: "Production, Services, and Import Tax",
			i18n.ES: "Impuesto sobre la Producción, los Servicios y la Importación",
		},
		// IPSI rates are complex and don't align well regular rates. Users are
		// recommended to include whatever percentage applies to their situation
		// directly in the invoice.
		Rates: []*tax.Rate{},
	},

	//
	// IRPF
	//
	{
		Code:     TaxCategoryIRPF,
		Retained: true,
		Name: i18n.String{
			i18n.EN: "IRPF",
			i18n.ES: "IRPF",
		},
		Map: cbc.CodeMap{
			KeyFacturaETaxTypeCode: "04",
		},
		Desc: i18n.String{
			i18n.EN: "Personal income tax.",
			i18n.ES: "Impuesto sobre la renta de las personas físicas.",
		},
		Rates: []*tax.Rate{
			{
				Key: TaxRatePro,
				Name: i18n.String{
					i18n.EN: "Professional Rate",
					i18n.ES: "Profesionales",
				},
				Values: []*tax.RateValue{
					{
						Since:   cal.NewDate(2015, 7, 12),
						Percent: num.MakePercentage(150, 3),
					},
					{
						Since:   cal.NewDate(2015, 1, 1),
						Percent: num.MakePercentage(190, 3),
					},
					{
						Since:   cal.NewDate(2012, 9, 1),
						Percent: num.MakePercentage(210, 3),
					},
					{
						Since:   cal.NewDate(2007, 1, 1),
						Percent: num.MakePercentage(150, 3),
					},
				},
			},
			{
				Key: TaxRateProStart,
				Name: i18n.String{
					i18n.EN: "Professional Starting Rate",
					i18n.ES: "Profesionales Inicio",
				},
				Values: []*tax.RateValue{
					{
						Since:   cal.NewDate(2007, 1, 1),
						Percent: num.MakePercentage(70, 3),
					},
				},
			},
			{
				Key: TaxRateCapital,
				Name: i18n.String{
					i18n.EN: "Rental or Interest Capital",
					i18n.ES: "Alquileres o Intereses de Capital",
				},
				Values: []*tax.RateValue{
					{
						Since:   cal.NewDate(2007, 1, 1),
						Percent: num.MakePercentage(190, 3),
					},
				},
			},
			{
				Key: TaxRateModules,
				Name: i18n.String{
					i18n.EN: "Modules Rate",
					i18n.ES: "Tipo Modulos",
				},
				Values: []*tax.RateValue{
					{
						Since:   cal.NewDate(2007, 1, 1),
						Percent: num.MakePercentage(10, 3),
					},
				},
			},
		},
	},
}
