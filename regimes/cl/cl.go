// Package cl provides the tax regime definition for Chile.
package cl

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the tax country code for Chile.
const CountryCode = "CL"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("cl", rules.GOBL.Add(CountryCode),
		taxIdentityRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(normalizeTaxIdentity)),
	)
}

// New provides the tax regime definition for Chile.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.CLP,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Chile",
			i18n.ES: "Chile",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Chile's Impuesto al Valor Agregado (IVA) is a value-added tax applied
				to most sales of goods and services. It is administered by the
				Servicio de Impuestos Internos (SII).

				Businesses and individuals are identified by their RUT (Rol Único
				Tributario), a numeric identifier of 7 or 8 digits followed by a
				check digit (0-9 or the letter "K"), validated with a modulus-11
				checksum.
			`),
			i18n.ES: here.Doc(`
		 	 	El Impuesto al Valor Agregado (IVA) en Chile se aplica a la mayoría
		 	 	de las ventas de bienes y servicios. Es administrado por el Servicio
		 	 	de Impuestos Internos (SII).

		 	 	Las empresas y personas se identifican con su RUT (Rol Único
		 	 	Tributario), un identificador numérico de 7 u 8 dígitos seguido de
		 		 un dígito verificador (0-9 o la letra "K"), validado con un
		 	 	algoritmo de módulo 11.
		 	`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "SII - Circular No. 47 (2003), VAT rate increase to 19%",
					i18n.ES: "SII - Circular N° 47 (2003), aumento de la tasa de IVA a 19%",
				},
				URL: "https://www.sii.cl/documentos/circulares/2003/circu47.htm",
			},
			{
				Title: i18n.String{
					i18n.EN: "Ley sobre Impuesto a las Ventas y Servicios (D.L. 825)",
					i18n.ES: "Ley sobre Impuesto a las Ventas y Servicios (D.L. 825)",
				},
				URL: "https://www.bcn.cl/leychile/navegar?idNorma=6369&idParte=10549143&idVersion=2025-08-01",
			},
		},
		TimeZone:   "America/Santiago",
		Categories: taxCategories(),
	}
}
