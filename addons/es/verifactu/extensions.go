package verifactu

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Extension keys for Verifactu
const (
	ExtKeyDocType     cbc.Key = "es-verifactu-doc-type"
	ExtKeyExemption   cbc.Key = "es-verifactu-exemption"
	ExtKeyTaxCategory cbc.Key = "es-verifactu-tax-category"
)

var extensions = []*cbc.KeyDefinition{
	{
		Key: ExtKeyExemption,
		Name: i18n.String{
			i18n.EN: "Verifactu Exemption code - L10",
			i18n.ES: "Código de Exención de Verifactu - L10",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Codes used by Verifactu for both "exempt", "not-subject", and reverse
				charge transactions. In the Verifactu format these are separated,
				but in order to simplify GOBL and be more closely aligned with
				other countries we've combined them into one.
			`),
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "E1",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 20. Exemptions in internal operations.",
					i18n.ES: "Exenta: por el artículo 20. Exenciones en operaciones interiores.",
				},
			},
			{
				Value: "E2",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 21. Exemptions in exports of goods.",
					i18n.ES: "Exenta: por el artículo 21. Exenciones en las exportaciones de bienes.",
				},
			},
			{
				Value: "E3",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 22. Exemptions in operations asimilated to exports.",
					i18n.ES: "Exenta: por el artículo 22. Exenciones en las operaciones asimiladas a las exportaciones.",
				},
			},
			{
				Value: "E4",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Articles 23 and 24. Exemptions related to temporary deposit, customs and fiscal regimes, and other situations.",
					i18n.ES: "Exenta: por el artículos 23 y 24. Exenciones relativas a las situaciones de depósito temporal, regímenes aduaneros y fiscales, y otras situaciones.",
				},
			},
			{
				Value: "E5",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 25. Exemptions in the delivery of goods destined to another Member State.",
					i18n.ES: "Exenta: por el artículo 25. Exenciones en las entregas de bienes destinados a otro Estado miembro.",
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
		Key: ExtKeyDocType,
		Name: i18n.String{
			i18n.EN: "Verifactu Invoice Type Code - L2",
			i18n.ES: "Código de Tipo de Factura de Verifactu - L2",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Invoice type code used to identify the type of invoice being sent.
			`),
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "F1",
				Name: i18n.String{
					i18n.EN: "Invoice (Article 6, 7.2 and 7.3 of RD 1619/2012)",
					i18n.ES: "Factura (Art. 6, 7.2 y 7.3 del RD 1619/2012)",
				},
			},
			{
				Value: "F2",
				Name: i18n.String{
					i18n.EN: "Simplified invoice (Article 6.1.d) of RD 1619/2012)",
					i18n.ES: "Factura Simplificada (Art. 6.1.d) del RD 1619/2012)",
				},
			},
			{
				Value: "F3",
				Name: i18n.String{
					i18n.EN: "Invoice in substitution of simplified invoices.",
					i18n.ES: "Factura emitida en sustitución de facturas simplificadas facturadas y declaradas.",
				},
			},
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
			i18n.EN: "Verifactu Tax Category Code - L1",
			i18n.ES: "Código de Tipo de Impuesto de Verifactu - L1",
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
				Value: "05",
				Name: i18n.String{
					i18n.EN: "Other",
					i18n.ES: "Otro",
				},
			},
		},
	},
}
