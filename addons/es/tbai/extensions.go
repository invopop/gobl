package tbai

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Extension keys for TicketBAI
const (
	ExtKeyRegion     cbc.Key = "es-tbai-region"
	ExtKeyExemption  cbc.Key = "es-tbai-exemption"
	ExtKeyProduct    cbc.Key = "es-tbai-product"
	ExtKeyCorrection cbc.Key = "es-tbai-correction"
)

var extensions = []*cbc.KeyDefinition{
	{
		Key: ExtKeyRegion,
		Name: i18n.String{
			i18n.EN: "TicketBAI Region Code",
			i18n.ES: "Código de Región TicketBAI",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Region codes are used by TicketBAI to differentiate between the different
				subdivisions of the Basque Country. This is used to determine the correct
				API endpoint to use when submitting documents.
			`),
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "VI",
				Name: i18n.String{
					i18n.EN: "Araba",
					i18n.ES: "Álava",
				},
			},
			{
				Value: "BI",
				Name: i18n.String{
					i18n.EN: "Bizkaia",
					i18n.ES: "Vizcaya",
				},
			},
			{
				Value: "SS",
				Name: i18n.String{
					i18n.EN: "Gipuzkoa",
					i18n.ES: "Guipúzcoa",
				},
			},
		},
	},
	{
		Key: ExtKeyProduct,
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
		Values: []*cbc.ValueDefinition{
			{
				Value: "goods",
				Name: i18n.String{
					i18n.EN: "Delivery of goods",
					i18n.ES: "Entrega de bienes",
				},
			},
			{
				Value: "services",
				Name: i18n.String{
					i18n.EN: "Provision of services",
					i18n.ES: "Prestacion de servicios",
				},
			},
			{
				Value: "resale",
				Name: i18n.String{
					i18n.EN: "Resale of goods without modification by vendor in the simplified regime",
					i18n.ES: "Reventa de bienes sin modificación por vendedor en regimen simplificado",
				},
			},
		},
	},
	{
		Key: ExtKeyExemption,
		Name: i18n.String{
			i18n.EN: "TicketBAI Exemption code",
			i18n.ES: "Código de Exención de TicketBAI",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Codes used by TicketBAI for both "exempt", "not-subject", and reverse
				charge transactions. In the TicketBAI format these are separated,
				but in order to simplify GOBL and be more closely aligned with
				other countries we've combined them into one.
			`),
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "E1",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 20 of the Foral VAT Law",
					i18n.ES: "Exenta: por el artículo 20 de la Norma Foral del IVA",
				},
			},
			{
				Value: "E2",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 21 of the Foral VAT Law",
					i18n.ES: "Exenta: por el artículo 21 de la Norma Foral del IVA",
				},
			},
			{
				Value: "E3",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 22 of the Foral VAT Law",
					i18n.ES: "Exenta: por el artículo 22 de la Norma Foral del IVA",
				},
			},
			{
				Value: "E4",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Articles 23 and 24 of the Foral VAT Law",
					i18n.ES: "Exenta: por el artículos 23 y 24 de la Norma Foral del IVA",
				},
			},
			{
				Value: "E5",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 25 of the Foral VAT law",
					i18n.ES: "Exenta: por el artículo 25 de la Norma Foral del IVA",
				},
			},
			{
				Value: "E6",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to other reasons",
					i18n.ES: "Exenta: por otra causa",
				},
			},
			{
				Value: "OT",
				Name: i18n.String{
					i18n.EN: "Not subject: pursuant to Article 7 of the VAT Law - other cases of non-subject",
					i18n.ES: "No sujeto: por el artículo 7 de la Ley del IVA - otros supuestos de no sujeción",
				},
			},
			{
				Value: "RL",
				Name: i18n.String{
					i18n.EN: "Not subject: pursuant to localization rules",
					i18n.ES: "No sujeto: por reglas de localización",
				},
			},
			{
				Value: "VT",
				Name: i18n.String{
					i18n.EN: "Not subject: sales made on behalf of third parties (amount not computable for VAT or IRPF purposes)",
					i18n.ES: "No sujeto: ventas realizadas por cuenta de terceros (importe no computable a efectos de IVA ni de IRPF)",
				},
			},
			{
				Value: "IE",
				Name: i18n.String{
					i18n.EN: "Not subject in the TAI due to localization rules, but foreign tax, IPS/IGIC or VAT from another EU member state is passed on",
					i18n.ES: "No sujeto en el TAI por reglas de localización, pero repercute impuesto extranjero, IPS/IGIC o IVA de otro estado miembro UE",
				},
			},
			/*
				// S1 is the default value for regular invoices, so we don't need to include it here
				// alongside the exemption codes.
				{
					Value: "S1",
					Name: i18n.String{
						i18n.EN: "Subject and not exempt: without reverse charge",
						i18n.ES: "Sujeto y no exenta: sin inversión del sujeto pasivo",
					},
				},
			*/
			{
				Value: "S2",
				Name: i18n.String{
					i18n.EN: "Subject and not exempt: with reverse charge",
					i18n.ES: "Sujeto y no exenta: con inversión del sujeto pasivo",
				},
			},
		},
	},
	{
		Key: ExtKeyCorrection,
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
		Values: []*cbc.ValueDefinition{
			{
				Value: "R1",
				Name: i18n.String{
					i18n.EN: "Rectified invoice: error based on law and Article 80 One, Two and Six of the Provincial Tax Law of VAT",
					i18n.ES: "Factura rectificativa: error fundado en derecho y Art. 80 Uno, Dos y Seis de la Norma Foral del IVA",
					i18n.EU: "Faktura zuzentzailea: zuzenbidean oinarritutako akatsa eta BEZaren Foru Arauaren 80.artikuluko Bat, Bi eta Sei",
				},
			},
			{
				Value: "R2",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: artículo 80 Tres de la Norma Foral del IVA",
					i18n.EN: "Rectified invoice: error based on law and Article 80 Three of the Provincial Tax Law of VAT",
					i18n.EU: "Faktura zuzentzailea: BEZari buruzko Foru Arauko 80. artikuluko Hiru",
				},
			},
			{
				Value: "R3",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: artículo 80 Cuatro de la Norma Foral del IVA",
					i18n.EN: "Rectified invoice: error based on law and Article 80 Four of the Provincial Tax Law of VAT",
					i18n.EU: "Faktura zuzentzailea: BEZari buruzko Foru Arauko 80. artikuluko Lau",
				},
			},
			{
				Value: "R4",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: Resto",
					i18n.EN: "Rectified invoice: Other",
					i18n.EU: "Faktura zuzentzailea: Gainerakoak",
				},
			},
			{
				Value: "R5",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: facturas simplificadas",
					i18n.EN: "Rectified invoice: simplified invoices",
					i18n.EU: "Faktura zuzentzaile: faktura erraztuetan",
				},
			},
		},
	},
}
