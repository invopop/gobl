package es

import (
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

// Spanish regime extension codes for local electronic formats.
const (
	ExtKeyTBAIExemption = "es-tbai-exemption"
	ExtKeyTBAIProduct   = "es-tbai-product"
)

var extensionKeys = []*tax.KeyDefinition{
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
		Keys: []*tax.KeyDefinition{
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
		Codes: []*tax.CodeDefinition{
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
					i18n.EN: "Not subject: pursuant to Article 7 of the VAT Law. Other cases of non-subject.",
					i18n.ES: "No sujeto: por el artículo 7 de la Ley del IVA. Otros supuestos de no sujeción.",
				},
			},
			{
				Code: "RL",
				Name: i18n.String{
					i18n.EN: "Not subject: pursuant to localization rules.",
					i18n.ES: "No sujeto: por reglas de localización.",
				},
			},
		},
	},
}
