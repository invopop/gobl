package arca

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Extension keys for Argentina ARCA v4
const (
	ExtKeyDocType      cbc.Key = "ar-arca-doc-type"
	ExtKeyConcept      cbc.Key = "ar-arca-concept"
	ExtKeyIdentityType cbc.Key = "ar-arca-identity-type"
	ExtKeyTaxType      cbc.Key = "ar-arca-tax-type"
	ExtKeyVATRate      cbc.Key = "ar-arca-vat-rate"
	ExtKeyVATStatus    cbc.Key = "ar-arca-vat-status"
)

// DocTypesA are document codes (Invoice A, Debit Note A, Credit Note A, and variants)
// Used for validating the document type against the VAT status.
var DocTypesA = []cbc.Code{
	"1",   // Invoice A
	"2",   // Debit Note A
	"3",   // Credit Note A
	"4",   // Receipt A
	"5",   // Cash Sales Note A
	"34",  // Invoices A per Annex I, Section A, subsection f), R.G. No. 1415
	"39",  // Other Vouchers A compliant with R.G. No. 1415
	"51",  // Invoice A with Legend 'Subject to Withholding'
	"52",  // Debit Note A with Legend 'Subject to Withholding'
	"53",  // Credit Note A with Legend 'Subject to Withholding'
	"54",  // Receipt A with Legend 'Subject to Withholding'
	"60",  // Sales Account and Product Settlement A
	"63",  // Settlement A
	"201", // MiPyMEs Electronic Credit Invoice (FCE) A
	"202", // MiPyMEs Electronic Debit Note (FCE) A
	"203", // MiPyMEs Electronic Credit Note (FCE) A
}

// DocTypesB are document codes (Invoice B, Debit Note B, Credit Note B, and FCE variants)
// Used for validating the document type against the VAT status.
var DocTypesB = []cbc.Code{
	"6",   // Invoice B
	"7",   // Debit Note B
	"8",   // Credit Note B
	"9",   // Receipt B
	"10",  // Cash Sales Note B
	"35",  // Vouchers B per Annex I, Section A, subsection f), R.G. No. 1415
	"40",  // Other Vouchers B compliant with R.G. No. 1415
	"61",  // Sales Account and Product Settlement B
	"64",  // Settlement B
	"206", // MiPyMEs Electronic Credit Invoice (FCE) B
	"207", // MiPyMEs Electronic Debit Note (FCE) B
	"208", // MiPyMEs Electronic Credit Note (FCE) B
}

// DocTypesC are document codes (Invoice C, Debit Note C, Credit Note C, and FCE variants)
// Used for validating the document type against the VAT status.
var DocTypesC = []cbc.Code{
	"11",  // Invoice C
	"12",  // Debit Note C
	"13",  // Credit Note C
	"15",  // Receipt C
	"211", // MiPyMEs Electronic Credit Invoice (FCE) C
	"212", // MiPyMEs Electronic Debit Note (FCE) C
	"213", // MiPyMEs Electronic Credit Note (FCE) C
}

// TypeUsedGoodsPurchaseInvoice is the code for the used goods purchase invoice
const TypeUsedGoodsPurchaseInvoice = "49"

// vatStatusesTypeA are VAT status codes that require type A documents
// Used for validating the document type against the different VAT statuses.
var vatStatusesTypeA = []cbc.Code{
	"1",  // Registered VAT Company
	"6",  // Monotributo Responsible
	"13", // Social Monotributista
	"16", // Promoted Independent Worker Monotributista
}

// DocTypesCreditNote are document codes for all credit notes (A, B, C, and FCE variants)
// Used for validating the arca document type extension agains GOBL bill.Invoice type.
var DocTypesCreditNote = []cbc.Code{
	"3",   // Credit Note A
	"8",   // Credit Note B
	"13",  // Credit Note C
	"53",  // Credit Note A with withholding legend
	"203", // MiPyMEs Electronic Credit Note (FCE) A
	"208", // MiPyMEs Electronic Credit Note (FCE) B
	"213", // MiPyMEs Electronic Credit Note (FCE) C
}

// DocTypesDebitNote are document codes for all debit notes (A, B, C, and FCE variants)
// Used for validating the arca document type extension agains GOBL bill.Invoice type.
var DocTypesDebitNote = []cbc.Code{
	"2",   // Debit Note A
	"7",   // Debit Note B
	"12",  // Debit Note C
	"52",  // Debit Note A with withholding legend
	"202", // MiPyMEs Electronic Debit Note (FCE) A
	"207", // MiPyMEs Electronic Debit Note (FCE) B
	"212", // MiPyMEs Electronic Debit Note (FCE) C
}

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyDocType,
		Name: i18n.String{
			i18n.EN: "Argentina ARCA Document Type",
			i18n.ES: "Tipo de comprobante Argentina ARCA",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Code used to identify the type of document being issued.

				This can always be set directly. If not provided, GOBL will automatically determine it during
				normalization based on the following rules:

				Type C (Simplified Tax Scheme / Monotributo supplier): If the invoice has the
				~monotax~ tag set, the document type will be ~11~ (Invoice C), ~13~ (Credit Note C),
				or ~12~ (Debit Note C).

				Type A (B2B with VAT-registered customers): If the customer's ~ar-arca-vat-status~ extension
				is one of ~1~ (Registered Company), ~6~ (Monotributo Responsible), ~13~ (Social Monotributista),
				or ~16~ (Promoted Independent Worker), the document type will be ~1~ (Invoice A), ~3~ (Credit
				Note A), or ~2~ (Debit Note A).

				Type B (B2C or other scenarios): For all other cases (final consumers, foreign customers,
				exempt entities, or when no customer is provided), the document type will be ~6~ (Invoice B),
				~8~ (Credit Note B), or ~7~ (Debit Note B).

				For other document types the value must be set manually.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Invoice A",
					i18n.ES: "Factura A",
				},
				Desc: i18n.String{
					i18n.EN: "Invoice issued by a VAT registered company to another VAT registered company or a monotributista",
					i18n.ES: "Factura emitida por un responsable inscripto a otro responsable inscripto o a un monotributista",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Debit Note A",
					i18n.ES: "Nota de Débito A",
				},
				Desc: i18n.String{
					i18n.EN: "Debit note issued by a VAT registered company to another VAT registered company or a monotributista",
					i18n.ES: "Nota de débito emitida por un responsable inscripto a otro responsable inscripto o a un monotributista",
				},
			},
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "Credit Note A",
					i18n.ES: "Nota de Crédito A",
				},
				Desc: i18n.String{
					i18n.EN: "Credit note issued by a VAT registered company to another VAT registered company or a monotributista",
					i18n.ES: "Nota de crédito emitida por un responsable inscripto a otro responsable inscripto o a un monotributista",
				},
			},
			{
				Code: "4",
				Name: i18n.String{
					i18n.EN: "Receipt A",
					i18n.ES: "Recibos A",
				},
			},
			{
				Code: "5",
				Name: i18n.String{
					i18n.EN: "Cash Sales Note A",
					i18n.ES: "Notas de Venta al contado A",
				},
			},
			{
				Code: "6",
				Name: i18n.String{
					i18n.EN: "Invoice B",
					i18n.ES: "Factura B",
				},
				Desc: i18n.String{
					i18n.EN: "Invoice issued by a registered taxpayer in the value added tax by operations performed with final consumers, exempt subjects in such tax, subjects not reached, subjects not categorized and with tourists from abroad.",
					i18n.ES: "Factura emitida por responsables inscriptos en el impuesto al valor agregado por operaciones realizadas con consumidores finales, sujetos exentos en dicho impuesto, sujetos no alcanzados, sujetos no categorizados y con turistas del extranjero.",
				},
			},
			{
				Code: "7",
				Name: i18n.String{
					i18n.EN: "Debit Note B",
					i18n.ES: "Nota de Débito B",
				},
				Desc: i18n.String{
					i18n.EN: "Debit note issued by a registered taxpayer in the value added tax by operations performed with final consumers, exempt subjects in such tax, subjects not reached, subjects not categorized and with tourists from abroad.",
					i18n.ES: "Nota de débito emitida por responsables inscriptos en el impuesto al valor agregado por operaciones realizadas con consumidores finales, sujetos exentos en dicho impuesto, sujetos no alcanzados, sujetos no categorizados y con turistas del extranjero.",
				},
			},
			{
				Code: "8",
				Name: i18n.String{
					i18n.EN: "Credit Note B",
					i18n.ES: "Nota de Crédito B",
				},
				Desc: i18n.String{
					i18n.EN: "Credit note issued by a registered taxpayer in the value added tax by operations performed with final consumers, exempt subjects in such tax, subjects not reached, subjects not categorized and with tourists from abroad.",
					i18n.ES: "Nota de crédito emitida por responsables inscriptos en el impuesto al valor agregado por operaciones realizadas con consumidores finales, sujetos exentos en dicho impuesto, sujetos no alcanzados, sujetos no categorizados y con turistas del extranjero.",
				},
			},
			{
				Code: "9",
				Name: i18n.String{
					i18n.EN: "Receipt B",
					i18n.ES: "Recibos B",
				},
			},
			{
				Code: "10",
				Name: i18n.String{
					i18n.EN: "Cash Sales Note B",
					i18n.ES: "Notas de Venta al contado B",
				},
			},
			{
				Code: "11",
				Name: i18n.String{
					i18n.EN: "Invoice C",
					i18n.ES: "Factura C",
				},
				Desc: i18n.String{
					i18n.EN: "Invoice issued by a monotributista.",
					i18n.ES: "Factura emitida por un monotributista.",
				},
			},
			{
				Code: "12",
				Name: i18n.String{
					i18n.EN: "Debit Note C",
					i18n.ES: "Nota de Débito C",
				},
				Desc: i18n.String{
					i18n.EN: "Debit note issued by a monotributista",
					i18n.ES: "Nota de débito emitida por un monotributista",
				},
			},
			{
				Code: "13",
				Name: i18n.String{
					i18n.EN: "Credit Note C",
					i18n.ES: "Nota de Crédito C",
				},
				Desc: i18n.String{
					i18n.EN: "Credit note issued by a monotributista",
					i18n.ES: "Nota de crédito emitida por un monotributista",
				},
			},
			{
				Code: "15",
				Name: i18n.String{
					i18n.EN: "Receipt C",
					i18n.ES: "Recibo C",
				},
			},
			{
				Code: "34",
				Name: i18n.String{
					i18n.EN: "Invoices A per Annex I, Section A, subsection f), R.G. No. 1415",
					i18n.ES: "Cbtes. A del Anexo I, Apartado A, inc. f), R.G. Nro. 1415",
				},
			},
			{
				Code: "35",
				Name: i18n.String{
					i18n.EN: "Vouchers B per Annex I, Section A, subsection f), R.G. No. 1415",
					i18n.ES: "Cbtes. B del Anexo I, Apartado A, inc. f), R.G. Nro. 1415",
				},
			},
			{
				Code: "39",
				Name: i18n.String{
					i18n.EN: "Other Vouchers A compliant with R.G. No. 1415",
					i18n.ES: "Otros comprobantes A que cumplan con R.G. Nro. 1415",
				},
			},
			{
				Code: "40",
				Name: i18n.String{
					i18n.EN: "Other Vouchers B compliant with R.G. No. 1415",
					i18n.ES: "Otros comprobantes B que cumplan con R.G. Nro. 1415",
				},
			},
			{
				Code: "49",
				Name: i18n.String{
					i18n.EN: "Used Goods Purchase Invoice to Final Consumer",
					i18n.ES: "Comprobante de Compra de Bienes Usados a Consumidor Final",
				},
			},
			{
				Code: "51",
				Name: i18n.String{
					i18n.EN: "Invoice A with Legend 'Subject to Withholding'",
					i18n.ES: "Factura A con Leyenda 'Operación Sujeta a Retención'",
				},
			},
			{
				Code: "52",
				Name: i18n.String{
					i18n.EN: "Debit Note A with Legend 'Subject to Withholding'",
					i18n.ES: "Nota de Débito A con Leyenda 'Operación Sujeta a Retención'",
				},
			},
			{
				Code: "53",
				Name: i18n.String{
					i18n.EN: "Credit Note A with Legend 'Subject to Withholding'",
					i18n.ES: "Nota de Crédito A con Leyenda 'Operación Sujeta a Retención'",
				},
			},
			{
				Code: "54",
				Name: i18n.String{
					i18n.EN: "Receipt A with Legend 'Subject to Withholding'",
					i18n.ES: "Recibo A con Leyenda 'Operación Sujeta a Retención'",
				},
			},
			{
				Code: "60",
				Name: i18n.String{
					i18n.EN: "Sales Account and Product Settlement A",
					i18n.ES: "Cta de Vta y Liquido prod. A",
				},
			},
			{
				Code: "61",
				Name: i18n.String{
					i18n.EN: "Sales Account and Product Settlement B",
					i18n.ES: "Cta de Vta y Liquido prod. B",
				},
			},
			{
				Code: "63",
				Name: i18n.String{
					i18n.EN: "Settlement A",
					i18n.ES: "Liquidacion A",
				},
			},
			{
				Code: "64",
				Name: i18n.String{
					i18n.EN: "Settlement B",
					i18n.ES: "Liquidacion B",
				},
			},
			{
				Code: "201",
				Name: i18n.String{
					i18n.EN: "MiPyMEs Electronic Credit Invoice (FCE) A",
					i18n.ES: "Factura de Crédito electrónica MiPyMEs (FCE) A",
				},
			},
			{
				Code: "202",
				Name: i18n.String{
					i18n.EN: "MiPyMEs Electronic Debit Note (FCE) A",
					i18n.ES: "Nota de Débito electrónica MiPyMEs (FCE) A",
				},
			},
			{
				Code: "203",
				Name: i18n.String{
					i18n.EN: "MiPyMEs Electronic Credit Note (FCE) A",
					i18n.ES: "Nota de Crédito electrónica MiPyMEs (FCE) A",
				},
			},
			{
				Code: "206",
				Name: i18n.String{
					i18n.EN: "MiPyMEs Electronic Credit Invoice (FCE) B",
					i18n.ES: "Factura de Crédito electrónica MiPyMEs (FCE) B",
				},
			},
			{
				Code: "207",
				Name: i18n.String{
					i18n.EN: "MiPyMEs Electronic Debit Note (FCE) B",
					i18n.ES: "Nota de Débito electrónica MiPyMEs (FCE) B",
				},
			},
			{
				Code: "208",
				Name: i18n.String{
					i18n.EN: "MiPyMEs Electronic Credit Note (FCE) B",
					i18n.ES: "Nota de Crédito electrónica MiPyMEs (FCE) B",
				},
			},
			{
				Code: "211",
				Name: i18n.String{
					i18n.EN: "MiPyMEs Electronic Credit Invoice (FCE) C",
					i18n.ES: "Factura de Crédito electrónica MiPyMEs (FCE) C",
				},
			},
			{
				Code: "212",
				Name: i18n.String{
					i18n.EN: "MiPyMEs Electronic Debit Note (FCE) C",
					i18n.ES: "Nota de Débito electrónica MiPyMEs (FCE) C",
				},
			},
			{
				Code: "213",
				Name: i18n.String{
					i18n.EN: "MiPyMEs Electronic Credit Note (FCE) C",
					i18n.ES: "Nota de Crédito electrónica MiPyMEs (FCE) C",
				},
			},
		},
	},
	{
		Key: ExtKeyConcept,
		Name: i18n.String{
			i18n.EN: "Argentina ARCA Concept (Goods, services, or both)",
			i18n.ES: "Concepto Argentina ARCA (Productos, servicios, o ambos)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Code used to identify the invoice concept, indicating whether the invoice covers
				goods, services, or both.

				This extension is automatically determined by GOBL based on the ~key~ field of each
				line item:

				- ~1~ (Products): All line items have ~key~ set to ~goods~.
				- ~2~ (Services): All line items have ~key~ empty, unset, or set to any value other than ~goods~.
				- ~3~ (Products and Services): The invoice contains a mix of both goods and services.

				When the concept is ~2~ (Services) or ~3~ (Products and Services), the invoice must
				include an ordering period and payment terms with due dates.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Goods",
					i18n.ES: "Productos",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Services",
					i18n.ES: "Servicios",
				},
			},
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "Goods and services",
					i18n.ES: "Productos y servicios",
				},
			},
		},
	},
	{
		Key: ExtKeyTaxType,
		Name: i18n.String{
			i18n.EN: "Argentina ARCA Tax Type",
			i18n.ES: "Tipo de Tributo Argentina ARCA",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Code used to identify the type of tax applied as a charge. This is used for taxes other than
				VAT that need to be included in the invoice (national taxes, provincial taxes, withholdings, etc.).

				These taxes must be added as invoice charges (not in line items) with the charge ~key~ set
				to ~tax~, a ~percent~ value (required, amounts are not accepted), and this extension to
				specify the tax type.

				Example - Tax applied to invoice total:

				~~~json
				"charges": [
					{
						"key": "tax",
						"percent": "3%",
						"ext": {
							"ar-arca-tax-type": "1"
						}
					}
				]
				~~~

				Example - Tax applied to a specific base:

				~~~json
				"charges": [
					{
						"key": "tax",
						"percent": "3%",
						"base": "1000.00",
						"ext": {
							"ar-arca-tax-type": "1"
						}
					}
				]
				~~~

				Validation rules:

				- When a charge has ~key: tax~, this extension is required.
				- The ~percent~ field is required when this extension is present.
				- When using code ~99~ (Other), the charge ~reason~ field is required to describe the tax.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "National Taxes",
					i18n.ES: "Impuestos Nacionales",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Provincial Taxes",
					i18n.ES: "Impuestos Provinciales",
				},
			},
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "Municipal Taxes",
					i18n.ES: "Impuestos Municipales",
				},
			},
			{
				Code: "4",
				Name: i18n.String{
					i18n.EN: "Internal Taxes",
					i18n.ES: "Impuestos Internos",
				},
			},
			{
				Code: "5",
				Name: i18n.String{
					i18n.EN: "Gross Income Tax",
					i18n.ES: "IIBB",
				},
			},
			{
				Code: "6",
				Name: i18n.String{
					i18n.EN: "VAT Prepayment",
					i18n.ES: "Percepción de IVA",
				},
			},
			{
				Code: "7",
				Name: i18n.String{
					i18n.EN: "Gross Income Tax Prepayment",
					i18n.ES: "Percepción de IIBB",
				},
			},
			{
				Code: "8",
				Name: i18n.String{
					i18n.EN: "Municipal Taxes Prepayment",
					i18n.ES: "Percepciones por Impuestos Municipales",
				},
			},
			{
				Code: "9",
				Name: i18n.String{
					i18n.EN: "Other Prepayments",
					i18n.ES: "Otras Percepciones",
				},
			},
			{
				Code: "13",
				Name: i18n.String{
					i18n.EN: "VAT Not Categorized Prepayment",
					i18n.ES: "Percepción de IVA a no Categorizado",
				},
			},
			{
				Code: "99",
				Name: i18n.String{
					i18n.EN: "Other",
					i18n.ES: "Otras",
				},
			},
		},
	},
	{
		Key: ExtKeyVATRate,
		Name: i18n.String{
			i18n.EN: "Argentina ARCA VAT Rate",
			i18n.ES: "Tasa de IVA Argentina ARCA",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Code used to identify the VAT rate applied to line items.

				This extension is automatically determined by GOBL based on the tax rate key defined in
				the Argentine tax regime:

				- ~zero~ rate (0%) → code ~3~
				- ~reduced~ rate (10.5%) → code ~4~
				- ~standard~ rate (21%) → code ~5~
				- ~increased~ rate (27%) → code ~6~

				For the special reduced rates of 5% and 2.5%, the extension must be set manually on the
				line tax as these are not covered by the standard rate keys:

				- 5% → code ~8~
				- 2.5% → code ~9~
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "0%",
					i18n.ES: "0%",
				},
			},
			{
				Code: "4",
				Name: i18n.String{
					i18n.EN: "10.5%",
					i18n.ES: "10.5%",
				},
			},
			{
				Code: "5",
				Name: i18n.String{
					i18n.EN: "21%",
					i18n.ES: "21%",
				},
			},
			{
				Code: "6",
				Name: i18n.String{
					i18n.EN: "27%",
					i18n.ES: "27%",
				},
			},
			{
				Code: "8",
				Name: i18n.String{
					i18n.EN: "5%",
					i18n.ES: "5%",
				},
			},
			{
				Code: "9",
				Name: i18n.String{
					i18n.EN: "2.5%",
					i18n.ES: "2.5%",
				},
			},
		},
	},
	{
		Key: ExtKeyIdentityType,
		Name: i18n.String{
			i18n.EN: "Argentina ARCA Customer Identity Type",
			i18n.ES: "Código de Documento identificatorio del comprador",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Code used to identify the type of identity document in the customer's ~identities~ field.

				This extension is typically only needed for B2C operations where the customer needs to
				provide an identification number but does not have a tax ID. Common examples include
				DNI (code ~96~), passport, or provincial identity cards.

				When the customer has a CUIT, CUIL, or a valid foreign tax identification number,
				these should be included in the ~tax_id~ field instead, and this extension is not required.

				Example usage for a final consumer with DNI:

				~~~json
				"customer": {
					"name": "John Doe",
					"identities": [
						{
							"code": "12345678",
							"ext": {
								"ar-arca-identity-type": "96"
							}
						}
					],
					"ext": {
						"ar-arca-vat-status": "5"
					}
				}
				~~~
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "00",
				Name: i18n.String{
					i18n.EN: "CI Capital Federal",
					i18n.ES: "CI Capital Federal",
				},
			},
			{
				Code: "01",
				Name: i18n.String{
					i18n.EN: "CI Buenos Aires",
					i18n.ES: "CI Buenos Aires",
				},
			},
			{
				Code: "02",
				Name: i18n.String{
					i18n.EN: "CI Catamarca",
					i18n.ES: "CI Catamarca",
				},
			},
			{
				Code: "03",
				Name: i18n.String{
					i18n.EN: "CI Córdoba",
					i18n.ES: "CI Córdoba",
				},
			},
			{
				Code: "04",
				Name: i18n.String{
					i18n.EN: "CI Corrientes",
					i18n.ES: "CI Corrientes",
				},
			},
			{
				Code: "05",
				Name: i18n.String{
					i18n.EN: "CI Entre Ríos",
					i18n.ES: "CI Entre Ríos",
				},
			},
			{
				Code: "06",
				Name: i18n.String{
					i18n.EN: "CI Jujuy",
					i18n.ES: "CI Jujuy",
				},
			},
			{
				Code: "07",
				Name: i18n.String{
					i18n.EN: "CI Mendoza",
					i18n.ES: "CI Mendoza",
				},
			},
			{
				Code: "08",
				Name: i18n.String{
					i18n.EN: "CI La Rioja",
					i18n.ES: "CI La Rioja",
				},
			},
			{
				Code: "09",
				Name: i18n.String{
					i18n.EN: "CI Salta",
					i18n.ES: "CI Salta",
				},
			},
			{
				Code: "10",
				Name: i18n.String{
					i18n.EN: "CI San Juan",
					i18n.ES: "CI San Juan",
				},
			},
			{
				Code: "11",
				Name: i18n.String{
					i18n.EN: "CI San Luis",
					i18n.ES: "CI San Luis",
				},
			},
			{
				Code: "12",
				Name: i18n.String{
					i18n.EN: "CI Santa Fe",
					i18n.ES: "CI Santa Fe",
				},
			},
			{
				Code: "13",
				Name: i18n.String{
					i18n.EN: "CI Santiago del Estero",
					i18n.ES: "CI Santiago del Estero",
				},
			},
			{
				Code: "14",
				Name: i18n.String{
					i18n.EN: "CI Tucumán",
					i18n.ES: "CI Tucumán",
				},
			},
			{
				Code: "16",
				Name: i18n.String{
					i18n.EN: "CI Chaco",
					i18n.ES: "CI Chaco",
				},
			},
			{
				Code: "17",
				Name: i18n.String{
					i18n.EN: "CI Chubut",
					i18n.ES: "CI Chubut",
				},
			},
			{
				Code: "18",
				Name: i18n.String{
					i18n.EN: "CI Formosa",
					i18n.ES: "CI Formosa",
				},
			},
			{
				Code: "19",
				Name: i18n.String{
					i18n.EN: "CI Misiones",
					i18n.ES: "CI Misiones",
				},
			},
			{
				Code: "20",
				Name: i18n.String{
					i18n.EN: "CI Neuquén",
					i18n.ES: "CI Neuquén",
				},
			},
			{
				Code: "21",
				Name: i18n.String{
					i18n.EN: "CI La Pampa",
					i18n.ES: "CI La Pampa",
				},
			},
			{
				Code: "22",
				Name: i18n.String{
					i18n.EN: "CI Río Negro",
					i18n.ES: "CI Río Negro",
				},
			},
			{
				Code: "23",
				Name: i18n.String{
					i18n.EN: "CI Santa Cruz",
					i18n.ES: "CI Santa Cruz",
				},
			},
			{
				Code: "24",
				Name: i18n.String{
					i18n.EN: "CI Tierra del Fuego",
					i18n.ES: "CI Tierra del Fuego",
				},
			},
			{
				Code: "80",
				Name: i18n.String{
					i18n.EN: "CUIT (Unique Tax Identification Number)",
					i18n.ES: "CUIT (Clave Única de Identificación Tributaria)",
				},
			},
			{
				Code: "86",
				Name: i18n.String{
					i18n.EN: "CUIL (Unique Labor Identification Number)",
					i18n.ES: "CUIL (Clave Única de Identificación Laboral)",
				},
			},
			{
				Code: "87",
				Name: i18n.String{
					i18n.EN: "CDI",
					i18n.ES: "CDI",
				},
			},
			{
				Code: "89",
				Name: i18n.String{
					i18n.EN: "LE",
					i18n.ES: "LE",
				},
			},
			{
				Code: "90",
				Name: i18n.String{
					i18n.EN: "LC",
					i18n.ES: "LC",
				},
			},
			{
				Code: "91",
				Name: i18n.String{
					i18n.EN: "Foreign CI",
					i18n.ES: "CI extranjera",
				},
			},
			{
				Code: "92",
				Name: i18n.String{
					i18n.EN: "In Process",
					i18n.ES: "en trámite",
				},
			},
			{
				Code: "93",
				Name: i18n.String{
					i18n.EN: "Birth Certificate",
					i18n.ES: "Acta Nacimiento",
				},
			},
			{
				Code: "94",
				Name: i18n.String{
					i18n.EN: "Passport",
					i18n.ES: "Pasaporte",
				},
			},
			{
				Code: "95",
				Name: i18n.String{
					i18n.EN: "CI Bs. As. RNP",
					i18n.ES: "CI Bs. As. RNP",
				},
			},
			{
				Code: "96",
				Name: i18n.String{
					i18n.EN: "DNI",
					i18n.ES: "DNI",
				},
			},
			{
				Code: "99",
				Name: i18n.String{
					i18n.EN: "Final Consumer",
					i18n.ES: "Consumidor Final",
				},
			},
		},
	},
	{
		Key: ExtKeyVATStatus,
		Name: i18n.String{
			i18n.EN: "Customer VAT Status",
			i18n.ES: "Condición frente al IVA del receptor",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Code used to identify the VAT status of the customer. This extension must be included
				in the customer's ~ext~ field when the customer is required.

				GOBL will automatically normalize this value based on the customer's tax identification:

				Customer without tax ID (final consumers): Defaults to ~5~ (Final Consumer). Also
				accepts ~4~ (Exempt Subject), ~7~ (Uncategorized), ~10~ (VAT Exempt Law 19640),
				or ~15~ (VAT Not Applicable).

				Customer with foreign tax ID: Defaults to ~9~ (Foreign Customer). Also accepts
				~8~ (Foreign Supplier).

				Customer with Argentine tax ID (CUIT/CUIL): Defaults to ~1~ (Registered VAT Company).
				Also accepts ~6~ (Monotributo), ~13~ (Social Monotributista), ~16~ (Promoted Independent
				Worker), ~4~ (Exempt Subject), ~7~ (Uncategorized), ~10~ (VAT Exempt Law 19640),
				or ~15~ (VAT Not Applicable).

				The VAT status is validated against the document type:

				- Type A documents (~1~, ~2~, ~3~, etc.): Require VAT status ~1~, ~6~, ~13~, or ~16~.
				- Type B documents (~6~, ~7~, ~8~, etc.): Cannot have VAT status ~1~, ~6~, ~13~, or ~16~.
				- Document type ~49~ (Used Goods Purchase Invoice): Requires VAT status ~5~ (Final Consumer).
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Registered VAT Company",
					i18n.ES: "IVA Responsable Inscripto",
				},
			},
			{
				Code: "4",
				Name: i18n.String{
					i18n.EN: "VAT Exempt Subject",
					i18n.ES: "IVA Sujeto Exento",
				},
			},
			{
				Code: "5",
				Name: i18n.String{
					i18n.EN: "Final Consumer",
					i18n.ES: "Consumidor Final",
				},
			},
			{
				Code: "6",
				Name: i18n.String{
					i18n.EN: "Monotributo Responsible",
					i18n.ES: "Responsable Monotributo",
				},
			},
			{
				Code: "7",
				Name: i18n.String{
					i18n.EN: "Uncategorized Subject",
					i18n.ES: "Sujeto No Categorizado",
				},
			},
			{
				Code: "8",
				Name: i18n.String{
					i18n.EN: "Foreign Supplier",
					i18n.ES: "Proveedor del Exterior",
				},
			},
			{
				Code: "9",
				Name: i18n.String{
					i18n.EN: "Foreign Customer",
					i18n.ES: "Cliente del Exterior",
				},
			},
			{
				Code: "10",
				Name: i18n.String{
					i18n.EN: "VAT Exempt - Law N° 19.640",
					i18n.ES: "IVA Liberado – Ley N° 19.640",
				},
			},
			{
				Code: "13",
				Name: i18n.String{
					i18n.EN: "Social Monotributista",
					i18n.ES: "Monotributista Social",
				},
			},
			{
				Code: "15",
				Name: i18n.String{
					i18n.EN: "VAT Not Applicable",
					i18n.ES: "IVA No Alcanzado",
				},
			},
			{
				Code: "16",
				Name: i18n.String{
					i18n.EN: "Promoted Independent Worker Monotributista",
					i18n.ES: "Monotributo Trabajador Independiente Promovido",
				},
			},
		},
	},
}
