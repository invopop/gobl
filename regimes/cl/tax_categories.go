package cl

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
	// VAT (IVA - Impuesto al Valor Agregado)
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
		Description: &i18n.String{
			i18n.EN: here.Doc(`
				Chile's IVA (Impuesto al Valor Agregado) is a consumption tax applied to the sale of goods and services. Chile has applied a single standard rate of 19% since October 1, 2003, making it one of the simpler VAT systems in Latin America with no reduced or super-reduced rates.

				The IVA applies to most goods and services unless specifically exempted. Common exemptions include certain financial services, educational services, and healthcare services. The tax is administered by the SII (Servicio de Impuestos Internos) and is a significant source of government revenue.
			`),
			i18n.ES: here.Doc(`
				El IVA (Impuesto al Valor Agregado) de Chile es un impuesto al consumo aplicado a la venta de bienes y servicios. Chile aplica una tasa estándar única del 19% desde el 1 de octubre de 2003, lo que lo convierte en uno de los sistemas de IVA más simples de América Latina, sin tasas reducidas o super-reducidas.

				El IVA se aplica a la mayoría de los bienes y servicios a menos que estén específicamente exentos. Las exenciones comunes incluyen ciertos servicios financieros, servicios educativos y servicios de salud. El impuesto es administrado por el SII (Servicio de Impuestos Internos) y es una fuente importante de ingresos del gobierno.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Decreto Ley Nº 825 - Law on VAT and Services",
					i18n.ES: "Decreto Ley Nº 825 - Ley sobre Impuesto a las Ventas y Servicios",
				},
				URL: "https://www.sii.cl/normativa_legislacion/sobreventasyservicios.pdf",
			},
			{
				Title: i18n.String{
					i18n.EN: "Ley 19888 - Law establishing 19% VAT rate",
					i18n.ES: "Ley 19888 - Ley que establece la tasa del IVA en 19%",
				},
				URL: "https://www.bcn.cl/leychile/Navegar?idNorma=213493",
			},
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
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2003, 10, 1),
						Percent: num.MakePercentage(19, 2),
					},
				},
			},
		},
	},
}
