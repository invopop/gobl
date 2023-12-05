package es

import (
	"github.com/invopop/gobl/i18n"
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
		Keys: []*tax.KeyDefinition{
			{
				Key: "goods",
				Name: i18n.String{
					i18n.ES: "Entrega de bienes",
					i18n.EN: "Delivery of goods",
				},
			},
			{
				Key: "services",
				Name: i18n.String{
					i18n.ES: "Prestacion de servicios",
					i18n.EN: "Provision of services",
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
		Codes: []*tax.CodeDefinition{
			{
				Code: "E1",
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to Article 20 of the Foral VAT Law",
					i18n.ES: "Exenta por el artículo 20 de la Norma Foral del IVA",
				},
			},
			{
				Code: "E2",
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to Article 21 of the Foral VAT Law",
					i18n.ES: "Exenta por el artículo 21 de la Norma Foral del IVA",
				},
			},
			{
				Code: "E3",
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to Article 22 of the Foral VAT Law",
					i18n.ES: "Exenta por el artículo 22 de la Norma Foral del IVA",
				},
			},
			{
				Code: "E4",
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to Articles 23 and 24 of the Foral VAT Law",
					i18n.ES: "Exenta por el artículos 23 y 24 de la Norma Foral del IVA",
				},
			},
			{
				Code: "E5",
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to Article 25 of the Foral VAT law",
					i18n.ES: "Exenta por el artículo 25 de la Norma Foral del IVA",
				},
			},
			{
				Code: "E6",
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to other reasons",
					i18n.ES: "Exenta por otra causa",
				},
			},
		},
	},
}
