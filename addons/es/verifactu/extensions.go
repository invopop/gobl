package verifactu

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Extension keys for Verifactu
const (
	ExtKeyDocType           cbc.Key = "es-verifactu-doc-type"
	ExtKeyIdentity          cbc.Key = "es-verifactu-identity"
	ExtKeyExemption         cbc.Key = "es-verifactu-exemption"
	ExtKeyTaxCategory       cbc.Key = "es-verifactu-tax-category"
	ExtKeyTaxClassification cbc.Key = "es-verifactu-tax-classification"
	ExtKeyTaxRegime         cbc.Key = "es-verifactu-tax-regime"
)

var extensions = []*cbc.KeyDefinition{
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
		Key: ExtKeyIdentity,
		Name: i18n.String{
			i18n.EN: "Verifactu Identity Type Code - L7",
			i18n.ES: "Código de Tipo de Identificación de Verifactu - L7",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Identity type code used to identify the type of identity being used.
			`),
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "02",
				Name: i18n.String{
					i18n.EN: "NIF-VAT",
					i18n.ES: "NIF-IVA",
				},
			},
			{
				Value: "03",
				Name: i18n.String{
					i18n.EN: "Passport",
					i18n.ES: "Pasaporte",
				},
			},
			{
				Value: "04",
				Name: i18n.String{
					i18n.EN: "Official identification document issued by the country or territory of residence",
					i18n.ES: "Documento oficial de identificación expedido por el país o territorio de residencia",
				},
			},
			{
				Value: "05",
				Name: i18n.String{
					i18n.EN: "Certificate of residence",
					i18n.ES: "Certificado de residencia",
				},
			},
			{
				Value: "06",
				Name: i18n.String{
					i18n.EN: "Other supporting document",
					i18n.ES: "Otro documento probatorio",
				},
			},
			{
				Value: "07",
				Name: i18n.String{
					i18n.EN: "Not registered",
					i18n.ES: "No censado",
				},
			},
		},
	},
	{
		Key: ExtKeyTaxRegime,
		Name: i18n.String{
			i18n.EN: "Verifactu Tax Regime Code - L8",
			i18n.ES: "Código de Régimen de Impuesto de Verifactu - L8",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Tax regime code used to identify the type of tax regime being applied to the invoice.
			`),
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "01",
				Name: i18n.String{
					i18n.EN: "General regime operation",
					i18n.ES: "Operación de régimen general",
				},
			},
			{
				Value: "02",
				Name: i18n.String{
					i18n.EN: "Export",
					i18n.ES: "Exportación",
				},
			},
			{
				Value: "03",
				Name: i18n.String{
					i18n.EN: "Special regime for used goods, art objects, antiques and collectibles",
					i18n.ES: "Operaciones a las que se aplique el régimen especial de bienes usados, objetos de arte, antigüedades y objetos de colección",
				},
			},
			{
				Value: "04",
				Name: i18n.String{
					i18n.EN: "Special regime for investment gold",
					i18n.ES: "Régimen especial del oro de inversión",
				},
			},
			{
				Value: "05",
				Name: i18n.String{
					i18n.EN: "Special regime for travel agencies",
					i18n.ES: "Régimen especial de las agencias de viajes",
				},
			},
			{
				Value: "06",
				Name: i18n.String{
					i18n.EN: "Special VAT regime for group entities (Advanced Level)",
					i18n.ES: "Régimen especial grupo de entidades en IVA (Nivel Avanzado)",
				},
			},
			{
				Value: "07",
				Name: i18n.String{
					i18n.EN: "Special cash accounting regime",
					i18n.ES: "Régimen especial del criterio de caja",
				},
			},
			{
				Value: "08",
				Name: i18n.String{
					i18n.EN: "Operations subject to IPSI/IGIC",
					i18n.ES: "Operaciones sujetas al IPSI/IGIC",
				},
			},
			{
				Value: "09",
				Name: i18n.String{
					i18n.EN: "Billing of travel agency services acting as intermediaries in name and on behalf of others",
					i18n.ES: "Facturación de las prestaciones de servicios de agencias de viaje que actúan como mediadoras en nombre y por cuenta ajena",
				},
			},
			{
				Value: "10",
				Name: i18n.String{
					i18n.EN: "Collection of professional fees or industrial property rights on behalf of third parties",
					i18n.ES: "Cobros por cuenta de terceros de honorarios profesionales o de derechos derivados de la propiedad industrial",
				},
			},
			{
				Value: "11",
				Name: i18n.String{
					i18n.EN: "Business premises rental operations",
					i18n.ES: "Operaciones de arrendamiento de local de negocio",
				},
			},
			{
				Value: "14",
				Name: i18n.String{
					i18n.EN: "Invoice with pending VAT accrual in work certifications for Public Administration",
					i18n.ES: "Factura con IVA pendiente de devengo en certificaciones de obra cuyo destinatario sea una Administración Pública",
				},
			},
			{
				Value: "15",
				Name: i18n.String{
					i18n.EN: "Invoice with pending VAT accrual in successive tract operations",
					i18n.ES: "Factura con IVA pendiente de devengo en operaciones de tracto sucesivo",
				},
			},
			{
				Value: "17",
				Name: i18n.String{
					i18n.EN: "Operation under OSS and IOSS regimes",
					i18n.ES: "Operación acogida a alguno de los regímenes previstos en el capítulo XI del título IX (OSS e IOSS)",
				},
			},
			{
				Value: "18",
				Name: i18n.String{
					i18n.EN: "Equivalence surcharge",
					i18n.ES: "Recargo de equivalencia",
				},
			},
			{
				Value: "19",
				Name: i18n.String{
					i18n.EN: "Operations included in the Special Regime for Agriculture, Livestock and Fisheries",
					i18n.ES: "Operaciones de actividades incluidas en el Régimen Especial de Agricultura, Ganadería y Pesca (REAGYP)",
				},
			},
			{
				Value: "20",
				Name: i18n.String{
					i18n.EN: "Simplified regime",
					i18n.ES: "Régimen simplificado",
				},
			},
		},
	},
	{
		Key: ExtKeyTaxClassification,
		Name: i18n.String{
			i18n.EN: "Verifactu Tax Classification Code - L9",
			i18n.ES: "Código de Clasificación de Impuesto de Verifactu - L9",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Tax classification code used to identify the type of tax being applied to the invoice.
			`),
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "S1",
				Name: i18n.String{
					i18n.EN: "Subject and Not Exempt - Without reverse charge",
					i18n.ES: "Operación Sujeta y No exenta - Sin inversión del sujeto pasivo",
				},
			},
			{
				Value: "S2",
				Name: i18n.String{
					i18n.EN: "Subject and Not Exempt - With reverse charge",
					i18n.ES: "Operación Sujeta y No exenta - Con Inversión del sujeto pasivo",
				},
			},
			{
				Value: "N1",
				Name: i18n.String{
					i18n.EN: "Not Subject - Articles 7, 14, others",
					i18n.ES: "Operación No Sujeta artículo 7, 14, otros",
				},
			},
			{
				Value: "N2",
				Name: i18n.String{
					i18n.EN: "Not Subject - Due to location rules",
					i18n.ES: "Operación No Sujeta por Reglas de localización",
				},
			},
		},
	},
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
}
