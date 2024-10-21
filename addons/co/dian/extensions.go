package dian

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Extension keys used in Colombia.
const (
	ExtKeyMunicipality cbc.Key = "co-dian-municipality"
	ExtKeyCreditCode   cbc.Key = "co-dian-credit-code"
	ExtKeyDebitCode    cbc.Key = "co-dian-debit-code"
)

var extensions = []*cbc.KeyDefinition{
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
		Values: []*cbc.ValueDefinition{
			{
				Value: "1",
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
				Value: "2",
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
				Value: "3",
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
				Value: "4",
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
				Value: "5",
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
		Values: []*cbc.ValueDefinition{
			{
				Value: "1",
				Name: i18n.String{
					i18n.EN: "Interest",
					i18n.ES: "Intereses",
				},
			},
			{
				Value: "2",
				Name: i18n.String{
					i18n.EN: "Pending charges",
					i18n.ES: "Gastos por cobrar",
				},
			},
			{
				Value: "3",
				Name: i18n.String{
					i18n.EN: "Change in value",
					i18n.ES: "Cambio del valor",
				},
			},
			{
				Value: "4",
				Name: i18n.String{
					i18n.EN: "Other",
					i18n.ES: "Otros",
				},
			},
		},
	},
}
