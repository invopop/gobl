package co

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

const (
	ExtKeyDIANCorrection cbc.Key = "co-dian-correction"
)

var extensionKeys = []*tax.KeyDefinition{
	{
		Key: ExtKeyDIANCorrection,
		Name: i18n.String{
			i18n.EN: "DIAN Correction Code",
			i18n.ES: "Código de corrección DIAN",
		},
		Codes: []*tax.CodeDefinition{
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
}
