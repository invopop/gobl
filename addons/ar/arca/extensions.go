package arca

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
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

// Tax type codes
const (
	ChargeTaxCodeNationalTaxes               cbc.Code = "1"
	ChargeTaxCodeProvincialTaxes             cbc.Code = "2"
	ChargeTaxCodeMunicipalTaxes              cbc.Code = "3"
	ChargeTaxCodeInternalTaxes               cbc.Code = "4"
	ChargeTaxCodeGrossIncomeTax              cbc.Code = "5"
	ChargeTaxCodeVATPrepayment               cbc.Code = "6"
	ChargeTaxCodeGrossIncomeTaxPrepayment    cbc.Code = "7"
	ChargeTaxCodeMunicipalTaxesPrepayment    cbc.Code = "8"
	ChargeTaxCodeOtherPrepayments            cbc.Code = "9"
	ChargeTaxCodeVATNotCategorizedPrepayment cbc.Code = "13"
	ChargeTaxCodeOther                       cbc.Code = "99"
)

const (
	VATStatusRegisteredCompany                       cbc.Code = "1"
	VATStatusExemptSubject                           cbc.Code = "4"
	VATStatusFinalConsumer                           cbc.Code = "5"
	VATStatusMonotributoResponsible                  cbc.Code = "6"
	VATStatusUncategorizedSubject                    cbc.Code = "7"
	VATStatusForeignSupplier                         cbc.Code = "8"
	VATStatusForeignCustomer                         cbc.Code = "9"
	VATStatusVATExemptLaw19640                       cbc.Code = "10"
	VATStatusSocialMonotributista                    cbc.Code = "13"
	VATStatusVATNotApplicable                        cbc.Code = "15"
	VATStatusPromotedIndependentWorkerMonotributista cbc.Code = "16"
)

const (
	ConceptGoods               cbc.Code = "1"
	ConceptServices            cbc.Code = "2"
	ConceptProductsAndServices cbc.Code = "3"
)

// DocTypesA are document codes (Invoice A, Debit Note A, Credit Note A, and variants)
// Used for validating the document type against the VAT status.
var DocTypesA = []cbc.Code{"1", "2", "3", "4", "5", "34", "39", "51", "52", "53", "54", "60", "63", "201", "202", "203"}

// DocTypesB are document codes (Invoice B, Debit Note B, Credit Note B, and FCE variants)
// Used for validating the document type against the VAT status.
var DocTypesB = []cbc.Code{"6", "7", "8", "9", "10", "35", "40", "61", "64", "206", "207", "208"}

// DocTypesC are document codes (Invoice C, Debit Note C, Credit Note C, and FCE variants)
// Used for validating the document type against the VAT status.
var DocTypesC = []cbc.Code{"11", "12", "13", "15", "211", "212", "213"}

// TypeUsedGoodsPurchaseInvoice is the code for the used goods purchase invoice
const TypeUsedGoodsPurchaseInvoice = "49"

// vatStatusesTypeA are VAT status codes that require type A documents
// Used for validating the document type against the different VAT statuses.
var vatStatusesTypeA = []cbc.Code{
	VATStatusRegisteredCompany,                       // 1
	VATStatusMonotributoResponsible,                  // 6
	VATStatusSocialMonotributista,                    // 13
	VATStatusPromotedIndependentWorkerMonotributista, // 16
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
			i18n.EN: "Argentina ARCA Concept (Product, service, or both)",
			i18n.ES: "Concepto Argentina ARCA (Producto, servicio, o ambos)",
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Products",
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
					i18n.EN: "Products and services",
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
			i18n.EN: "Argentina ARCA Identity Type",
			i18n.ES: "Tipo de Identidad Argentina ARCA",
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
			i18n.EN: "Argentina ARCA Customer VAT Status",
			i18n.ES: "Condición IVA del Receptor Argentina ARCA",
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
