package es

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

// Tax Rates which may be used by TicketBAI in the Basque Country.
const (
	TaxRateExempt    cbc.Key = "exempt"
	TaxRateArticle20 cbc.Key = "article-20"
	TaxRateArticle21 cbc.Key = "article-21"
	TaxRateArticle22 cbc.Key = "article-22"
	TaxRateArticle23 cbc.Key = "article-23"
	TaxRateArticle25 cbc.Key = "article-25"
	TaxRateOther     cbc.Key = "other"
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
		Codes: cbc.CodeSet{
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
						Percent: num.MakePercentage(21, 2),
					},
					{
						Since:   cal.NewDate(2010, 7, 1),
						Percent: num.MakePercentage(18, 2),
					},
					{
						Since:   cal.NewDate(1995, 1, 1),
						Percent: num.MakePercentage(16, 2),
					},
					{
						Since:   cal.NewDate(1993, 1, 1),
						Percent: num.MakePercentage(15, 2),
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
						Percent:   num.MakePercentage(21, 2),
						Surcharge: num.NewPercentage(52, 3),
					},
					{
						Since:     cal.NewDate(2010, 7, 1),
						Percent:   num.MakePercentage(18, 2),
						Surcharge: num.NewPercentage(4, 2),
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
						Percent: num.MakePercentage(10, 2),
					},
					{
						Since:   cal.NewDate(2010, 7, 1),
						Percent: num.MakePercentage(8, 2),
					},
					{
						Since:   cal.NewDate(1995, 1, 1),
						Percent: num.MakePercentage(7, 2),
					},
					{
						Since:   cal.NewDate(1993, 1, 1),
						Percent: num.MakePercentage(6, 2),
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
						Percent:   num.MakePercentage(10, 2),
						Surcharge: num.NewPercentage(14, 3),
					},
					{
						Since:     cal.NewDate(2010, 7, 1),
						Percent:   num.MakePercentage(8, 2),
						Surcharge: num.NewPercentage(1, 2),
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
						Percent: num.MakePercentage(4, 2),
					},
					{
						Since:   cal.NewDate(1993, 1, 1),
						Percent: num.MakePercentage(3, 2),
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
						Percent:   num.MakePercentage(4, 2),
						Surcharge: num.NewPercentage(5, 3),
					},
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateArticle20),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to Article 20 of the Foral VAT Law",
					i18n.ES: "Exenta por el artículo 20 de la Norma Foral del IVA",
				},
				Codes: cbc.CodeSet{
					KeyTicketBAICausaExencion: "E1",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateArticle21),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to Article 21 of the Foral VAT Law",
					i18n.ES: "Exenta por el artículo 21 de la Norma Foral del IVA",
				},
				Codes: cbc.CodeSet{
					KeyTicketBAICausaExencion: "E2",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateArticle22),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to Article 22 of the Foral VAT Law",
					i18n.ES: "Exenta por el artículo 22 de la Norma Foral del IVA",
				},
				Codes: cbc.CodeSet{
					KeyTicketBAICausaExencion: "E3",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateArticle23),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to Articles 23 and 24 of the Foral VAT Law",
					i18n.ES: "Exenta por el artículos 23 y 24 de la Norma Foral del IVA",
				},
				Codes: cbc.CodeSet{
					KeyTicketBAICausaExencion: "E4",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateArticle25),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to Article 25 of the Foral VAT law",
					i18n.ES: "Exenta por el artículo 25 de la Norma Foral del IVA",
				},
				Codes: cbc.CodeSet{
					KeyTicketBAICausaExencion: "E5",
				},
			},
			{
				Key:    TaxRateExempt.With(TaxRateOther),
				Exempt: true,
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to other reasons",
					i18n.ES: "Exenta por otra causa",
				},
				Codes: cbc.CodeSet{
					KeyTicketBAICausaExencion: "E6",
				},
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
		Codes: cbc.CodeSet{
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
		Codes: cbc.CodeSet{
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
		Codes: cbc.CodeSet{
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
						Percent: num.MakePercentage(15, 2),
					},
					{
						Since:   cal.NewDate(2015, 1, 1),
						Percent: num.MakePercentage(19, 2),
					},
					{
						Since:   cal.NewDate(2012, 9, 1),
						Percent: num.MakePercentage(21, 2),
					},
					{
						Since:   cal.NewDate(2007, 1, 1),
						Percent: num.MakePercentage(15, 2),
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
						Percent: num.MakePercentage(7, 2),
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
						Percent: num.MakePercentage(19, 2),
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
						Percent: num.MakePercentage(1, 2),
					},
				},
			},
		},
	},
}
