package es

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func taxCategories() []*tax.CategoryDef {
	return []*tax.CategoryDef{
		//
		// VAT
		//
		{
			Code:     tax.CategoryVAT,
			Retained: false,
			Name: i18n.String{
				i18n.EN: "VAT",
				i18n.ES: "IVA",
				i18n.GL: "IVE",
				i18n.EU: "BEZ",
				i18n.CA: "IVA",
			},
			Title: i18n.String{
				i18n.EN: "Value Added Tax",
				i18n.ES: "Impuesto sobre el Valor Añadido",
				i18n.GL: "Imposto sobre o Valor Engadido",
				i18n.EU: "Balio Erantsiaren Zerga",
				i18n.CA: "Impost sobre el Valor Afegit",
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
						i18n.GL: "Tipo Xeral",
						i18n.EU: "Tasa Orokorra",
						i18n.CA: "Tipus general",
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
						i18n.GL: "Tipo Xeral + Recargo de Equivalencia",
						i18n.EU: "Tasa Orokorra + Baliokidetasun Errekargua",
						i18n.CA: "Tipus general + Recàrrec d'equivalència",
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
						i18n.GL: "Tipo Reducido",
						i18n.EU: "Tasa Murriztua",
						i18n.CA: "Tipus reduït",
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
						i18n.GL: "Tipo Reducido + Recargo de Equivalencia",
						i18n.EU: "Tasa Murriztua + Baliokidetasun Errekargua",
						i18n.CA: "Tipus reduït + Recàrrec d'equivalència",
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
						i18n.GL: "Tipo Superreducido",
						i18n.EU: "Tasa Oso Murriztua",
						i18n.CA: "Tipus superreduït",
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
						i18n.GL: "Tipo Superreducido + Recargo de Equivalencia",
						i18n.EU: "Tasa Oso Murriztua + Baliokidetasun Errekargua",
						i18n.CA: "Tipus superreduït + Recàrrec d'equivalència",
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
				i18n.GL: "IGIC",
				i18n.EU: "IGIC",
				i18n.CA: "IGIC",
			},
			Title: i18n.String{
				i18n.EN: "Canary Island General Indirect Tax",
				i18n.ES: "Impuesto General Indirecto Canario",
				i18n.GL: "Imposto Xeral Indirecto Canario",
				i18n.EU: "Kanariar Uharteetako Zerga Orokor Zeharkakoa",
				i18n.CA: "Impost General Indirecte Canari",
			},
			// Use the same global VAT keys as IGIC is effectively a local VAT, unlike
			// IPSI which has more in common with a sales tax.
			Keys: tax.GlobalVATKeys(),
			// This is a subset of the possible rates, notably the "increased" rates applied on luxury
			// items and some professional services are not included here. Users are recommended to include whatever
			// percentage applies to their situation directly in the invoice.
			Rates: []*tax.RateDef{
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateGeneral,
					Name: i18n.String{
						i18n.EN: "General Rate",
						i18n.ES: "Tipo General",
						i18n.GL: "Tipo Xeral",
						i18n.EU: "Tasa Orokorra",
						i18n.CA: "Tipus general",
					},
					Values: []*tax.RateValueDef{
						{
							Percent: num.MakePercentage(70, 3),
						},
					},
				},
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateReduced,
					Name: i18n.String{
						i18n.EN: "Reduced Rate",
						i18n.ES: "Tipo Reducido",
						i18n.GL: "Tipo Reducido",
						i18n.EU: "Tasa Murriztua",
						i18n.CA: "Tipus reduït",
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
				i18n.GL: "IPSI",
				i18n.EU: "IPSI",
				i18n.CA: "IPSI",
			},
			Title: i18n.String{
				i18n.EN: "Production, Services, and Import Tax",
				i18n.ES: "Impuesto sobre la Producción, los Servicios y la Importación",
				i18n.GL: "Imposto sobre a Produción, os Servizos e a Importación",
				i18n.EU: "Ekoizpen, Zerbitzu eta Inportazio Zerga",
				i18n.CA: "Impost sobre la Producció, els Serveis i la Importació",
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
				i18n.GL: "IRPF",
				i18n.EU: "IRPF",
				i18n.CA: "IRPF",
			},
			Title: i18n.String{
				i18n.EN: "Personal income tax.",
				i18n.ES: "Impuesto sobre la renta de las personas físicas.",
				i18n.GL: "Imposto sobre a renda das persoas físicas.",
				i18n.EU: "Pertsona fisikoen errentaren gaineko zerga.",
				i18n.CA: "Impost sobre la renda de les persones físiques.",
			},
			Rates: []*tax.RateDef{
				{
					Rate: TaxRatePro,
					Name: i18n.String{
						i18n.EN: "Professional Rate",
						i18n.ES: "Profesionales",
						i18n.GL: "Profesionais",
						i18n.EU: "Profesionalak",
						i18n.CA: "Professionals",
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
						i18n.GL: "Inicio Profesionais",
						i18n.EU: "Hasiera Profesionalak",
						i18n.CA: "Inici Professionals",
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
						i18n.GL: "Alugueres ou Intereses de Capital",
						i18n.EU: "Errentamenduak edo Kapitalaren Interesak",
						i18n.CA: "Lloguers o Interessos de Capital",
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
						i18n.GL: "Tipo Módulos",
						i18n.EU: "Modulu Tasa",
						i18n.CA: "Tipus Mòduls",
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
}
