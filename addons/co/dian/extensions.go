package dian

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Extension keys used in Colombia.
const (
	ExtKeyMunicipality         cbc.Key = "co-dian-municipality"
	ExtKeyCreditCode           cbc.Key = "co-dian-credit-code"
	ExtKeyDebitCode            cbc.Key = "co-dian-debit-code"
	ExtKeyFiscalResponsibility cbc.Key = "co-dian-fiscal-responsibility"
)

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyMunicipality,
		Name: i18n.String{
			i18n.EN: "DIAN Municipality Code",
			i18n.ES: "Código de municipio DIAN",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The municipality code as defined by the DIAN.

				For further details on the list of possible codes, see:

				 * https://www.dian.gov.co/atencionciudadano/formulariosinstructivos/Formularios/2007/Codigos_municipios_2007.pdf
				 * https://github.com/ALAxHxC/MunicipiosDane
			`),
		},
		Pattern: `^\d{5}$`,
	},
	{
		Key: ExtKeyCreditCode,
		Name: i18n.String{
			i18n.EN: "Credit Code",
			i18n.ES: "Código de Crédito",
		},
		Desc: i18n.String{
			i18n.EN: "DIAN correction code for credit notes",
			i18n.ES: "Código de corrección DIAN para notas de crédito",
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Partial refund",
					i18n.ES: "Devolución parcial",
				},
				Desc: i18n.String{
					i18n.EN: "Partial refund of part of the goods or services.",
					i18n.ES: "Devolución de parte de los bienes; no aceptación de partes del servicio.",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Revoked",
					i18n.ES: "Anulación",
				},
				Desc: i18n.String{
					i18n.EN: "Previous document has been completely cancelled.",
					i18n.ES: "Anulación de la factura anterior.",
				},
			},
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "Discount",
					i18n.ES: "Descuento",
				},
				Desc: i18n.String{
					i18n.EN: "Partial or total discount.",
					i18n.ES: "Rebaja o descuento parcial o total.",
				},
			},
			{
				Code: "4",
				Name: i18n.String{
					i18n.EN: "Adjustment",
					i18n.ES: "Ajuste",
				},
				Desc: i18n.String{
					i18n.EN: "Price adjustment.",
					i18n.ES: "Ajuste de precio.",
				},
			},
			{
				Code: "5",
				Name: i18n.String{
					i18n.EN: "Other",
					i18n.ES: "Otros",
				},
			},
		},
	},
	{
		Key: ExtKeyDebitCode,
		Name: i18n.String{
			i18n.EN: "Debit Code",
			i18n.ES: "Código de Débito",
		},
		Desc: i18n.String{
			i18n.EN: "DIAN correction code for debit notes",
			i18n.ES: "Código de corrección DIAN para notas de débito",
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Interest",
					i18n.ES: "Intereses",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Pending charges",
					i18n.ES: "Gastos por cobrar",
				},
			},
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "Change in value",
					i18n.ES: "Cambio del valor",
				},
			},
			{
				Code: "4",
				Name: i18n.String{
					i18n.EN: "Other",
					i18n.ES: "Otros",
				},
			},
		},
	},
	{
		Key: ExtKeyFiscalResponsibility,
		Name: i18n.String{
			i18n.EN: "Fiscal Responsibility Code",
			i18n.ES: "Código de Responsabilidad Fiscal",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The fiscal responsibility code as defined by the DIAN for Colombian electronic invoicing.
				Maps to the UBL's "TaxLevelCode" field.

				For further details and the list of codes, see:

				  * https://www.dian.gov.co/impuestos/factura-electronica/Documents/Caja-de-herramientas-FE-V1-9.zip
				    (see Anexo Tecnico/Tablas Referenciadas, table 13.2.6.1)
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "O-13",
				Name: i18n.String{
					i18n.EN: "Major taxpayer",
					i18n.ES: "Gran contribuyente",
				},
			},
			{
				Code: "O-15",
				Name: i18n.String{
					i18n.EN: "Self-withholder",
					i18n.ES: "Autorretenedor",
				},
			},
			{
				Code: "O-23",
				Name: i18n.String{
					i18n.EN: "VAT withholding agent",
					i18n.ES: "Agente de retención IVA",
				},
			},
			{
				Code: "O-47",
				Name: i18n.String{
					i18n.EN: "Simple tax regime",
					i18n.ES: "Régimen simple de tributación",
				},
			},
			{
				Code: "R-99-PN",
				Name: i18n.String{
					i18n.EN: "Not applicable – Others",
					i18n.ES: "No aplica – Otros",
				},
				Desc: i18n.String{
					i18n.EN: "Used when the issuer/acquirer does not have any of the first 4 responsibilities. Applies to legal entities, individuals, or final consumers.",
					i18n.ES: "Se utiliza cuando el emisor/adquiriente no cuenta con las primeras 4 responsabilidades. Aplica para personas jurídicas, personas naturales o consumidor final.",
				},
			},
		},
	},
}
