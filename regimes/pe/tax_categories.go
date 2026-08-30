package pe

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.CategoryDef{
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "IGV",
			i18n.ES: "IGV",
		},
		Title: i18n.String{
			i18n.EN: "General Sales Tax",
			i18n.ES: "Impuesto General a las Ventas",
		},
		Description: &i18n.String{
			i18n.EN: here.Doc(`
				The IGV is Peru's value-added tax. The rate charged on invoices
				always combines the IGV itself with the Municipal Promotion Tax
				(IPM); both are levied together and shown as a single amount, so
				GOBL models them as one rate.
			`),
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "General Rate",
					i18n.ES: "Tasa General",
				},
				Description: i18n.String{
					i18n.EN: "Standard rate for most goods and services: 16% IGV plus 2% IPM (Ley 29666).",
					i18n.ES: "Tasa estándar para la mayoría de bienes y servicios: 16% de IGV más 2% de IPM (Ley 29666).",
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(180, 3),
						Since:   cal.NewDate(2011, 3, 1),
					},
					{
						Percent: num.MakePercentage(190, 3),
						Since:   cal.NewDate(2003, 8, 1),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate (tourism and restaurant MYPEs)",
					i18n.ES: "Tasa Reducida (MYPE de turismo y restaurantes)",
				},
				Description: i18n.String{
					i18n.EN: here.Doc(`
						Temporary rate for qualifying micro and small enterprises
						in the restaurant, hotel and tourist accommodation
						sectors (Ley 31556, extended by Ley 32219 with the IPM
						component raised by Ley 32387). Eligibility depends on
						the supplier's registration with SUNAT; the special
						regime is legislated to end on 31 December 2027.
					`),
					i18n.ES: here.Doc(`
						Tasa temporal para micro y pequeñas empresas calificadas
						de los rubros de restaurantes, hoteles y alojamientos
						turísticos (Ley 31556, prorrogada por la Ley 32219 con el
						componente del IPM elevado por la Ley 32387). La
						elegibilidad depende de la inscripción del emisor ante
						SUNAT; el régimen especial concluye el 31 de diciembre
						de 2027.
					`),
				},
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(150, 3),
						Since:   cal.NewDate(2027, 1, 1),
					},
					{
						Percent: num.MakePercentage(105, 3),
						Since:   cal.NewDate(2026, 1, 1),
					},
					{
						Percent: num.MakePercentage(100, 3),
						Since:   cal.NewDate(2022, 9, 1),
					},
				},
			},
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "SUNAT - IGV concept, rate and taxed operations",
					i18n.ES: "SUNAT - Concepto, tasa y operaciones gravadas - IGV",
				},
				URL: "https://orientacion.sunat.gob.pe/3053-concepto-tasa-y-operaciones-gravadas-igv-empresas",
				At:  cal.NewDateTime(2026, 8, 30, 0, 0, 0),
			},
			{
				Title: i18n.String{
					i18n.EN: "SUNAT - Reduced IGV for restaurant, hotel and tourist accommodation MYPEs",
					i18n.ES: "SUNAT - Reducción del IGV para restaurantes, hoteles y alojamientos turísticos",
				},
				URL: "https://orientacion.sunat.gob.pe/reduccion-del-igv-para-restaurantes-hoteles-y-alojamiento-turisticos-1",
				At:  cal.NewDateTime(2026, 8, 30, 0, 0, 0),
			},
			{
				Title: i18n.String{
					i18n.ES: "Ley 32219 - prórroga y modificación de la Ley 31556",
				},
				URL: "https://busquedas.elperuano.pe/dispositivo/NL/2357993-3",
			},
		},
	},
}
