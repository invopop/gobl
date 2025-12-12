package arca

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

const (
	ExtKeyDocType      cbc.Key = "ar-arca-doc-type"
	ExtKeyConcept      cbc.Key = "ar-arca-concept"
	ExtKeyIdentityType cbc.Key = "ar-arca-identity-type"
	ExtKeyTributeType  cbc.Key = "ar-arca-tribute-type"
	ExtKeyVATRate      cbc.Key = "ar-arca-vat-rate"
	ExtKeyVATStatus    cbc.Key = "ar-arca-vat-status"
)

// Tribute type codes
const (
	TributeTypeNationalTaxes               cbc.Code = "1"
	TributeTypeProvincialTaxes             cbc.Code = "2"
	TributeTypeMunicipalTaxes              cbc.Code = "3"
	TributeTypeInternalTaxes               cbc.Code = "4"
	TributeTypeGrossIncomeTax              cbc.Code = "5"
	TributeTypeVATPrepayment               cbc.Code = "6"
	TributeTypeGrossIncomeTaxPrepayment    cbc.Code = "7"
	TributeTypeMunicipalTaxesPrepayment    cbc.Code = "8"
	TributeTypeOtherPrepayments            cbc.Code = "9"
	TributeTypeVATNotCategorizedPrepayment cbc.Code = "13"
	TributeTypeOther                       cbc.Code = "99"
)

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyDocType,
		Name: i18n.String{
			i18n.EN: "Argentina ARCA Document Type",
			i18n.ES: "Tipo de comprobante Argentina ARCA",
		},
		Values: []*cbc.Definition{
			{
				Code: "001",
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
				Code: "002",
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
				Code: "003",
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
				Code: "006",
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
				Code: "007",
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
				Code: "008",
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
				Code: "011",
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
				Code: "012",
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
				Code: "013",
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
				Code: "019",
				Name: i18n.String{
					i18n.EN: "Export Invoice",
					i18n.ES: "Factura de Exportación",
				},
				Desc: i18n.String{
					i18n.EN: "Export invoice issued by registered taxpayers or monotributistas that export goods or services to clients outside the country",
					i18n.ES: "Factura emitida por responsables inscriptos o monotributistas que exportan bienes o servicios a clientes en el exterior.",
				},
			},
			{
				Code: "020",
				Name: i18n.String{
					i18n.EN: "Debit Note for Foreign Operations",
					i18n.ES: "Nota de Débito por Operaciones con el Exterior",
				},
				Desc: i18n.String{
					i18n.EN: "Debit note issued by companies and businesses that export goods or services to clients outside the country",
					i18n.ES: "Nota de débito emitida por las empresas y comercios que exportan bienes o servicios a clientes en el exterior.",
				},
			},
			{
				Code: "021",
				Name: i18n.String{
					i18n.EN: "Credit Note for Foreign Operations",
					i18n.ES: "Nota de Crédito por Operaciones con el Exterior",
				},
				Desc: i18n.String{
					i18n.EN: "Credit note issued by companies and businesses that export goods or services to clients outside the country",
					i18n.ES: "Nota de crédito emitida por las empresas y comercios que exportan bienes o servicios a clientes en el exterior.",
				},
			},
			{
				Code: "195",
				Name: i18n.String{
					i18n.EN: "Invoice T",
					i18n.ES: "Factura T",
				},
				Desc: i18n.String{
					i18n.EN: "Invoice issued by a registered taxpayer to tourists from abroad for accommodation operations.",
					i18n.ES: "Factura emitida por un responsable inscripto a turistas del extranjero para operaciones de alojamiento.",
				},
			},
			{
				Code: "196",
				Name: i18n.String{
					i18n.EN: "Debit Note T",
					i18n.ES: "Nota de Débito T",
				},
				Desc: i18n.String{
					i18n.EN: "Debit note issued by a registered taxpayer to tourists from abroad for accommodation operations.",
					i18n.ES: "Nota de débito emitida por un responsable inscripto a turistas del extranjero para operaciones de alojamiento.",
				},
			},
			{
				Code: "197",
				Name: i18n.String{
					i18n.EN: "Credit Note T",
					i18n.ES: "Nota de Crédito T",
				},
				Desc: i18n.String{
					i18n.EN: "Credit note issued by a registered taxpayer to tourists from abroad for accommodation operations.",
					i18n.ES: "Nota de crédito emitida por un responsable inscripto a turistas del extranjero para operaciones de alojamiento.",
				},
			},
		},
	},
	{
		Key: ExtKeyConcept,
		Name: i18n.String{
			i18n.EN: "Argentina ARCA Concept",
			i18n.ES: "Concepto Argentina ARCA",
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
		Key: ExtKeyTributeType,
		Name: i18n.String{
			i18n.EN: "Argentina ARCA Tribute Type",
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
