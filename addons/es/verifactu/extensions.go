package verifactu

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Extension keys for TicketBAI
const (
	ExtKeyExemption   cbc.Key = "es-verifactu-exemption"
	ExtKeyCorrection  cbc.Key = "es-verifactu-correction"
	ExtKeyTaxCategory cbc.Key = "es-verifactu-tax-category"
)

var extensions = []*cbc.KeyDefinition{
	{
		Key: ExtKeyExemption,
		Name: i18n.String{
			i18n.EN: "Verifactu Exemption code",
			i18n.ES: "Código de Exención de Verifactu",
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
		},
	},
	{
		Key: ExtKeyCorrection,
		Name: i18n.String{
			i18n.EN: "Verifactu Rectification Type Code",
			i18n.ES: "Verifactu Código de Factura Rectificativa",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Corrected or rectified invoices that need to be sent in the Verifactu format
				require a specific type code to be defined alongside the preceding invoice
				data.
			`),
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "R1",
				Name: i18n.String{
					i18n.EN: "Rectified invoice: error based on law and Article 80 One, Two and Six LIVA",
					i18n.ES: "Factura rectificativa: error fundado en derecho y Art. 80 Uno, Dos y Seis LIVA",
				},
			},
			{
				Value: "R2",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: artículo 80.3",
					i18n.EN: "Rectified invoice: error based on law and Article 80.3",
				},
			},
			{
				Value: "R3",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: artículo 80.4",
					i18n.EN: "Rectified invoice: error based on law and Article 80.4",
				},
			},
			{
				Value: "R4",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: Resto",
					i18n.EN: "Rectified invoice: Other",
				},
			},
			{
				Value: "R5",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: facturas simplificadas",
					i18n.EN: "Rectified invoice: simplified invoices",
				},
			},
		},
	},
	{
		Key: ExtKeyTaxCategory,
		Name: i18n.String{
			i18n.EN: "Verifactu Tax Category Code",
			i18n.ES: "Verifactu Impuesto de Aplicacion",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
			Tax category code used to identify the type of tax being applied to the invoice.
			The code must be one of the predefined values that correspond to the main Spanish
			tax regimes: IVA (Value Added Tax), IPSI (Tax on Production, Services and Imports),
			IGIC (Canary Islands General Indirect Tax), or Other.
			`),
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "01",
				Name: i18n.String{
					i18n.EN: "IVA",
					i18n.ES: "IVA",
				},
			},
			{
				Value: "02",
				Name: i18n.String{
					i18n.EN: "IPSI",
					i18n.ES: "IPSI",
				},
			},
			{
				Value: "03",
				Name: i18n.String{
					i18n.EN: "IGIC",
					i18n.ES: "IGIC",
				},
			},
			{
				Value: "04",
				Name: i18n.String{
					i18n.EN: "Other",
					i18n.ES: "Otro",
				},
			},
		},
	},
}
