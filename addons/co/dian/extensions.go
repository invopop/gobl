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
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "DIAN Municipality Codes",
				},
				URL:         "https://www.dian.gov.co/atencionciudadano/formulariosinstructivos/Formularios/2007/Codigos_municipios_2007.pdf",
				ContentType: "application/pdf",
			},
			{
				Title: i18n.String{
					i18n.EN: "Municipalities of Colombia - Github",
				},
				URL: "https://github.com/ALAxHxC/MunicipiosDane",
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The municipality code as defined by the DIAN.

				Set the 5-digit code for the municipality where the issuer is located in both
				the supplier and customer:

				~~~js
				"supplier": {
					"name": "EXAMPLE SUPPLIER S.A.S.",
					"tax_id": {
						"country": "CO",
						"code": "9014514812"
					},
					"ext": {
						"co-dian-municipality": "11001" // Bogotá, D.C.
					},
					// [...]
				},
				"customer": {
					"name": "EXAMPLE CUSTOMER S.A.S.",
					"tax_id": {
						"country": "CO",
						"code": "9014514805"
					},
					"ext": {
						"co-dian-municipality": "05001" // Medellín
					},
					// [...]
				},
				~~~
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
			i18n.EN: here.Doc(`
				The DIAN correction code is required when issuing credit notes in Colombia
				and is not automatically assigned by GOBL. It must be be included inside the
				~preceding~ document references.

				The extension will be offered as an option in the invoice correction process.

				Usage example:

				~~~js
				"preceding": [
					{
						"uuid": "0190e063-7676-7000-8c58-2db7172a4e58",
						"type": "standard",
						"series": "SETT",
						"code": "1010006",
						"issue_date": "2024-07-23",
						"reason": "Reason",
						"stamps": [
							{
								"prv": "dian-cude",
								"val": "57601dd1ab69213ccf8cfd5894f2e9fbfe23643f3a24e2f2526a5bb88d058a0842fffcb339694b6704dc105a9d813327"
							}
						],
						"ext": {
							"co-dian-credit-code": "3"
						}
					}
				],
				~~~
			`),
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
			i18n.EN: here.Doc(`
				The DIAN correction code is required when issuing debit notes in Colombia
				and is not automatically assigned by GOBL.

				The extension will be offered as an option in the invoice correction process.
			`),
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
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "DIAN Fiscal Responsibility Codes, see Anexo Tecnico/Tablas Referenciadas, table 13.2.6.1",
				},
				URL:         "https://www.dian.gov.co/impuestos/factura-electronica/Documents/Caja-de-herramientas-FE-V1-9.zip",
				ContentType: "application/zip",
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The fiscal responsibility code as defined by the DIAN for Colombian electronic invoicing.
				Maps to the UBL's ~TaxLevelCode~ field.

				The DIAN requires that Colombian invoices specify the fiscal responsibilities of the
				supplier or customer using specific codes. If no value is provided, GOBL will
				automatically set ~R-99-PN~ as the default.

				| Code    | Description                   |
				| ------- | ----------------------------- |
				| O-13    | Gran contribuyente            |
				| O-15    | Autorretenedor                |
				| O-23    | Agente de retención IVA       |
				| O-47    | Régimen simple de tributación |
				| R-99-PN | No aplica - Otros             |

				For example:

				~~~js
				"customer": {
					"name": "EXAMPLE CUSTOMER S.A.S.",
					"tax_id": {
						"country": "CO",
						"code": "9014514812"
					},
					"ext": {
						"co-dian-fiscal-responsibility": "O-13"
					}
				}
				~~~

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
