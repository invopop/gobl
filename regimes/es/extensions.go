package es

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Spanish regime extension codes for local electronic formats.
const (
	ExtKeyTBAIExemption      = "es-tbai-exemption"
	ExtKeyTBAIProduct        = "es-tbai-product"
	ExtKeyTBAICorrection     = "es-tbai-correction"
	ExtKeyFacturaECorrection = "es-facturae-correction"
)

var extensionKeys = []*cbc.KeyDefinition{
	{
		Key: ExtKeyFacturaECorrection,
		Name: i18n.String{
			i18n.EN: "FacturaE Change",
			i18n.ES: "Cambio de FacturaE",
		},
		Desc: i18n.String{
			i18n.EN: "FacturaE requires a specific and single code that explains why the previous invoice is being corrected.",
			i18n.ES: "FacturaE requiere un código específico y único que explique por qué se está corrigiendo la factura anterior.",
		},
		// Codes take from FacturaE XSD
		Codes: []*cbc.CodeDefinition{
			{
				Code: "01",
				Name: i18n.String{
					i18n.EN: "Invoice code",
					i18n.ES: "Número de la factura",
				},
			},
			{
				Code: "02",
				Name: i18n.String{
					i18n.EN: "Invoice series",
					i18n.ES: "Serie de la factura",
				},
			},
			{
				Code: "03",
				Name: i18n.String{
					i18n.EN: "Issue date",
					i18n.ES: "Fecha expedición",
				},
			},
			{
				Code: "04",
				Name: i18n.String{
					i18n.EN: "Name and surnames/Corporate name - Issuer (Sender)",
					i18n.ES: "Nombre y apellidos/Razón Social-Emisor",
				},
			},
			{
				Code: "05",
				Name: i18n.String{
					i18n.EN: "Name and surnames/Corporate name - Receiver",
					i18n.ES: "Nombre y apellidos/Razón Social-Receptor",
				},
			},
			{
				Code: "06",
				Name: i18n.String{
					i18n.EN: "Issuer's Tax Identification Number",
					i18n.ES: "Identificación fiscal Emisor/obligado",
				},
			},
			{
				Code: "07",
				Name: i18n.String{
					i18n.EN: "Receiver's Tax Identification Number",
					i18n.ES: "Identificación fiscal Receptor",
				},
			},
			{
				Code: "08",
				Name: i18n.String{
					i18n.EN: "Supplier's address",
					i18n.ES: "Domicilio Emisor/Obligado",
				},
			},
			{
				Code: "09",
				Name: i18n.String{
					i18n.EN: "Customer's address",
					i18n.ES: "Domicilio Receptor",
				},
			},
			{
				Code: "10",
				Name: i18n.String{
					i18n.EN: "Item line",
					i18n.ES: "Detalle Operación",
				},
			},
			{
				Code: "11",
				Name: i18n.String{
					i18n.EN: "Applicable Tax Rate",
					i18n.ES: "Porcentaje impositivo a aplicar",
				},
			},
			{
				Code: "12",
				Name: i18n.String{
					i18n.EN: "Applicable Tax Amount",
					i18n.ES: "Cuota tributaria a aplicar",
				},
			},
			{
				Code: "13",
				Name: i18n.String{
					i18n.EN: "Applicable Date/Period",
					i18n.ES: "Fecha/Periodo a aplicar",
				},
			},
			{
				Code: "14",
				Name: i18n.String{
					i18n.EN: "Invoice Class",
					i18n.ES: "Clase de factura",
				},
			},
			{
				Code: "15",
				Name: i18n.String{
					i18n.EN: "Legal literals",
					i18n.ES: "Literales legales",
				},
			},
			{
				Code: "16",
				Name: i18n.String{
					i18n.EN: "Taxable Base",
					i18n.ES: "Base imponible",
				},
			},
			{
				Code: "80",
				Name: i18n.String{
					i18n.EN: "Calculation of tax outputs",
					i18n.ES: "Cálculo de cuotas repercutidas",
				},
			},
			{
				Code: "81",
				Name: i18n.String{
					i18n.EN: "Calculation of tax inputs",
					i18n.ES: "Cálculo de cuotas retenidas",
				},
			},
			{
				Code: "82",
				Name: i18n.String{
					i18n.EN: "Taxable Base modified due to return of packages and packaging materials",
					i18n.ES: "Base imponible modificada por devolución de envases / embalajes",
				},
			},
			{
				Code: "83",
				Name: i18n.String{
					i18n.EN: "Taxable Base modified due to discounts and rebates",
					i18n.ES: "Base imponible modificada por descuentos y bonificaciones",
				},
			},
			{
				Code: "84",
				Name: i18n.String{
					i18n.EN: "Taxable Base modified due to firm court ruling or administrative decision",
					i18n.ES: "Base imponible modificada por resolución firme, judicial o administrativa",
				},
			},
			{
				Code: "85",
				Name: i18n.String{
					i18n.EN: "Taxable Base modified due to unpaid outputs where there is a judgement opening insolvency proceedings",
					i18n.ES: "Base imponible modificada cuotas repercutidas no satisfechas. Auto de declaración de concurso",
				},
			},
		},
	},
	{
		Key: ExtKeyTBAIProduct,
		Name: i18n.String{
			i18n.EN: "TicketBAI Product Key",
			i18n.ES: "Clave de Producto TicketBAI",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Product keys are used by TicketBAI to differentiate between -exported- goods
				and services. It may be useful to classify all products regardless of wether
				they are exported or not.

				There is an additional exception case for goods that are resold without modification
				when the supplier is in the simplified tax regime. For must purposes this special
				case can be ignored.
			`),
		},
		Keys: []*cbc.KeyDefinition{
			{
				Key: "goods",
				Name: i18n.String{
					i18n.EN: "Delivery of goods",
					i18n.ES: "Entrega de bienes",
				},
			},
			{
				Key: "services",
				Name: i18n.String{
					i18n.EN: "Provision of services",
					i18n.ES: "Prestacion de servicios",
				},
			},
			{
				Key: "resale",
				Name: i18n.String{
					i18n.EN: "Resale of goods without modification by vendor in the simplified regime",
					i18n.ES: "Reventa de bienes sin modificación por vendedor en regimen simplificado",
				},
			},
		},
	},
	{
		Key: ExtKeyTBAIExemption,
		Name: i18n.String{
			i18n.EN: "TicketBAI Exemption code",
			i18n.ES: "Código de Exención de TicketBAI",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Codes used by TicketBAI for both "exempt" and "not-subject"
				transactions. In the TicketBAI format these are separated,
				but in order to simplify GOBL and be more closely aligned with
				other countries we've combined them into one.
			`),
		},
		Codes: []*cbc.CodeDefinition{
			{
				Code: "E1",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 20 of the Foral VAT Law",
					i18n.ES: "Exenta: por el artículo 20 de la Norma Foral del IVA",
				},
			},
			{
				Code: "E2",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 21 of the Foral VAT Law",
					i18n.ES: "Exenta: por el artículo 21 de la Norma Foral del IVA",
				},
			},
			{
				Code: "E3",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 22 of the Foral VAT Law",
					i18n.ES: "Exenta: por el artículo 22 de la Norma Foral del IVA",
				},
			},
			{
				Code: "E4",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Articles 23 and 24 of the Foral VAT Law",
					i18n.ES: "Exenta: por el artículos 23 y 24 de la Norma Foral del IVA",
				},
			},
			{
				Code: "E5",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 25 of the Foral VAT law",
					i18n.ES: "Exenta: por el artículo 25 de la Norma Foral del IVA",
				},
			},
			{
				Code: "E6",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to other reasons",
					i18n.ES: "Exenta: por otra causa",
				},
			},
			{
				Code: "OT",
				Name: i18n.String{
					i18n.EN: "Not subject: pursuant to Article 7 of the VAT Law - other cases of non-subject",
					i18n.ES: "No sujeto: por el artículo 7 de la Ley del IVA - otros supuestos de no sujeción",
				},
			},
			{
				Code: "RL",
				Name: i18n.String{
					i18n.EN: "Not subject: pursuant to localization rules",
					i18n.ES: "No sujeto: por reglas de localización",
				},
			},
			{
				Code: "VT",
				Name: i18n.String{
					i18n.EN: "Not subject: sales made on behalf of third parties (amount not computable for VAT or IRPF purposes)",
					i18n.ES: "No sujeto: ventas realizadas por cuenta de terceros (importe no computable a efectos de IVA ni de IRPF)",
				},
			},
			{
				Code: "IE",
				Name: i18n.String{
					i18n.EN: "Not subject in the TAI due to localization rules, but foreign tax, IPS/IGIC or VAT from another EU member state is passed on",
					i18n.ES: "No sujeto en el TAI por reglas de localización, pero repercute impuesto extranjero, IPS/IGIC o IVA de otro estado miembro UE",
				},
			},
			{
				Code: "S1",
				Name: i18n.String{
					i18n.EN: "Subject and not exempt: without reverse charge",
					i18n.ES: "Sujeto y no exenta: sin inversión del sujeto pasivo",
				},
			},
			{
				Code: "S2",
				Name: i18n.String{
					i18n.EN: "Subject and not exempt: with reverse charge",
					i18n.ES: "Sujeto y no exenta: con inversión del sujeto pasivo",
				},
			},
		},
	},
	{
		Key: ExtKeyTBAICorrection,
		Name: i18n.String{
			i18n.EN: "TicketBAI Rectification Type Code",
			i18n.ES: "TicketBAI Código de Factura Rectificativa",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Corrected or rectified invoices that need to be sent in the TicketBAI format
				require a specific type code to be defined alongside the preceding invoice
				data.
			`),
		},
		// Codes taken from TicketBAI XSD
		Codes: []*cbc.CodeDefinition{
			{
				Code: "R1",
				Name: i18n.String{
					i18n.EN: "Rectified invoice: error based on law and Article 80 One, Two and Six of the Provincial Tax Law of VAT",
					i18n.ES: "Factura rectificativa: error fundado en derecho y Art. 80 Uno, Dos y Seis de la Norma Foral del IVA",
					i18n.EU: "Faktura zuzentzailea: zuzenbidean oinarritutako akatsa eta BEZaren Foru Arauaren 80.artikuluko Bat, Bi eta Sei",
				},
			},
			{
				Code: "R2",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: artículo 80 Tres de la Norma Foral del IVA",
					i18n.EN: "Rectified invoice: error based on law and Article 80 Three of the Provincial Tax Law of VAT",
					i18n.EU: "Faktura zuzentzailea: BEZari buruzko Foru Arauko 80. artikuluko Hiru",
				},
			},
			{
				Code: "R3",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: artículo 80 Cuatro de la Norma Foral del IVA",
					i18n.EN: "Rectified invoice: error based on law and Article 80 Four of the Provincial Tax Law of VAT",
					i18n.EU: "Faktura zuzentzailea: BEZari buruzko Foru Arauko 80. artikuluko Lau",
				},
			},
			{
				Code: "R4",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: Resto",
					i18n.EN: "Rectified invoice: Other",
					i18n.EU: "Faktura zuzentzailea: Gainerakoak",
				},
			},
			{
				Code: "R5",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: facturas simplificadas",
					i18n.EN: "Rectified invoice: simplified invoices",
					i18n.EU: "Faktura zuzentzaile: faktura erraztuetan",
				},
			},
		},
	},
}
