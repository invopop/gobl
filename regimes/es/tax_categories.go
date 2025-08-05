package es

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.CategoryDef{
	//
	// VAT
	//
	{
		Code:     tax.CategoryVAT,
		Retained: false,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.ES: "IVA",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.ES: "Impuesto sobre el Valor Añadido",
		},
		Description: &i18n.String{
			i18n.EN: here.Doc(`
				Known in Spanish as "Impuesto sobre el Valor Añadido" (IVA), is a consumption tax
				applied to the purchase of goods and services. It's a tax on the value added at
				each stage of production or distribution. Spain, as a member of the European Union,
				follows the EU's VAT Directive, but with specific rates and exemptions tailored
				to its local needs.
			`),
		},
		Keys: tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.ES: "Tipo General",
				},
				Values: []*tax.RateValueDef{
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
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral.With(TaxRateEquivalence),
				Name: i18n.String{
					i18n.EN: "Standard Rate + Equivalence Surcharge",
					i18n.ES: "Tipo General + Recargo de Equivalencia",
				},
				Values: []*tax.RateValueDef{
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
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.ES: "Tipo Reducido",
				},
				Values: []*tax.RateValueDef{
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
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced.With(TaxRateEquivalence),
				Name: i18n.String{
					i18n.EN: "Reduced Rate + Equivalence Surcharge",
					i18n.ES: "Tipo Reducido + Recargo de Equivalencia",
				},
				Values: []*tax.RateValueDef{
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
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Super-Reduced Rate",
					i18n.ES: "Tipo Superreducido",
				},
				Values: []*tax.RateValueDef{
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
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced.With(TaxRateEquivalence),
				Name: i18n.String{
					i18n.EN: "Super-Reduced Rate + Equivalence Surcharge",
					i18n.ES: "Tipo Superreducido + Recargo de Equivalencia",
				},
				Values: []*tax.RateValueDef{
					{
						Since:     cal.NewDate(1995, 1, 1),
						Percent:   num.MakePercentage(40, 3),
						Surcharge: num.NewPercentage(5, 3),
					},
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
		Title: i18n.String{
			i18n.EN: "Canary Island General Indirect Tax",
			i18n.ES: "Impuesto General Indirecto Canario",
		},
		// This is a subset of the possible rates.
		Rates: []*tax.RateDef{
			{
				Rate: tax.RateZero,
				Name: i18n.String{
					i18n.EN: "Zero Rate",
					i18n.ES: "Tipo Cero",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(0, 3),
					},
				},
			},
			{
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General Rate",
					i18n.ES: "Tipo General",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(70, 3),
					},
				},
			},
			{
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.ES: "Tipo Reducido",
				},
				Values: []*tax.RateValueDef{
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
		Title: i18n.String{
			i18n.EN: "Production, Services, and Import Tax",
			i18n.ES: "Impuesto sobre la Producción, los Servicios y la Importación",
		},
		// IPSI rates are complex and don't align well regular rates. Users are
		// recommended to include whatever percentage applies to their situation
		// directly in the invoice.
		Rates: []*tax.RateDef{},
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
		Title: i18n.String{
			i18n.EN: "Personal income tax.",
			i18n.ES: "Impuesto sobre la renta de las personas físicas.",
		},
		Rates: []*tax.RateDef{
			{
				Rate: TaxRatePro,
				Name: i18n.String{
					i18n.EN: "Professional Rate",
					i18n.ES: "Profesionales",
				},
				Values: []*tax.RateValueDef{
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
				Rate: TaxRateProStart,
				Name: i18n.String{
					i18n.EN: "Professional Starting Rate",
					i18n.ES: "Profesionales Inicio",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2007, 1, 1),
						Percent: num.MakePercentage(70, 3),
					},
				},
			},
			{
				Rate: TaxRateCapital,
				Name: i18n.String{
					i18n.EN: "Rental or Interest Capital",
					i18n.ES: "Alquileres o Intereses de Capital",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2007, 1, 1),
						Percent: num.MakePercentage(190, 3),
					},
				},
			},
			{
				Rate: TaxRateModules,
				Name: i18n.String{
					i18n.EN: "Modules Rate",
					i18n.ES: "Tipo Modulos",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2007, 1, 1),
						Percent: num.MakePercentage(10, 3),
					},
				},
			},
		},
	},
}
