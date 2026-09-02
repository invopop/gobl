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
					i18n.EN: here.Doc(`
						Standard rate for most goods and services, combining the
						IGV and the IPM. The internal split between the two is
						set by law year to year (16% + 2% until 2025, recomposed
						annually by Ley 32387 from 2026); the 18% total charged
						on invoices does not change.
					`),
					i18n.ES: here.Doc(`
						Tasa estándar para la mayoría de bienes y servicios, que
						combina el IGV y el IPM. La composición interna entre
						ambos la fija la ley año a año (16% + 2% hasta 2025,
						recompuesta anualmente por la Ley 32387 desde 2026); el
						18% total del comprobante no varía.
					`),
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
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("TUO de la Ley del IGV e ISC (D.S. 055-99-EF)"),
				URL:   "https://www.sunat.gob.pe/legislacion/igv/ley/",
			},
			{
				Title: i18n.String{
					i18n.EN: "SUNAT - IGV concept, rate and taxed operations",
					i18n.ES: "SUNAT - Concepto, tasa y operaciones gravadas - IGV",
				},
				URL: "https://orientacion.sunat.gob.pe/3053-concepto-tasa-y-operaciones-gravadas-igv-empresas",
				At:  cal.NewDateTime(2026, 8, 30, 0, 0, 0),
			},
		},
	},
}
