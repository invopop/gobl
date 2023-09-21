package es

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// List of correction codes derived from the Spanish FacturaE format.
const (
	CorrectionKeyCode            cbc.Key = "code"       // Invoice Code
	CorrectionKeySeries          cbc.Key = "series"     // Invoice series number
	CorrectionKeyIssueDate       cbc.Key = "issue-date" // Issue Date
	CorrectionKeySupplier        cbc.Key = "supplier"   // General supplier details
	CorrectionKeyCustomer        cbc.Key = "customer"   // General customer details
	CorrectionKeySupplierName    cbc.Key = "supplier-name"
	CorrectionKeyCustomerName    cbc.Key = "customer-name"
	CorrectionKeySupplierTaxID   cbc.Key = "supplier-tax-id"
	CorrectionKeyCustomerTaxID   cbc.Key = "customer-tax-id"
	CorrectionKeySupplierAddress cbc.Key = "supplier-addr"
	CorrectionKeyCustomerAddress cbc.Key = "customer-addr"
	CorrectionKeyLine            cbc.Key = "line"
	CorrectionKeyPeriod          cbc.Key = "period"
	CorrectionKeyType            cbc.Key = "type"
	CorrectionKeyLegalDetails    cbc.Key = "legal-details"
	CorrectionKeyTaxRate         cbc.Key = "tax-rate"
	CorrectionKeyTaxAmount       cbc.Key = "tax-amount"
	CorrectionKeyTaxBase         cbc.Key = "tax-base"
	CorrectionKeyTax             cbc.Key = "tax"          // General issue with tax calculations
	CorrectionKeyTaxRetained     cbc.Key = "tax-retained" // Error in retained tax calculations
	CorrectionKeyRefund          cbc.Key = "refund"       // Goods or materials have been returned to supplier
	CorrectionKeyDiscount        cbc.Key = "discount"     // New discounts or rebates added
	CorrectionKeyJudicial        cbc.Key = "judicial"     // Court ruling or administrative decision
	CorrectionKeyInsolvency      cbc.Key = "insolvency"   // the customer is insolvent and cannot pay
)

// List of correction methods derived from the Spanish FacturaE format.
const (
	CorrectionMethodKeyComplete   cbc.Key = "complete"   // everything has changed
	CorrectionMethodKeyPartial    cbc.Key = "partial"    // only differences corrected
	CorrectionMethodKeyDiscount   cbc.Key = "discount"   // deducted from future invoices
	CorrectionMethodKeyAuthorized cbc.Key = "authorized" // Permitted by tax agency
)

// correctionList contains an array of Key Definitions describing each of the acceptable
// correction keys, descriptions, and their "code" as determined by the FacturaE specifications.
var correctionChangesList = []*tax.KeyDefinition{
	{
		Key: CorrectionKeyCode,
		Name: i18n.String{
			i18n.EN: "Invoice code",
			i18n.ES: "Número de la factura",
		},
		Map: cbc.CodeMap{KeyFacturaE: "01"},
	},
	{
		Key: CorrectionKeySeries,
		Name: i18n.String{
			i18n.EN: "Invoice series",
			i18n.ES: "Serie de la factura",
		},
		Map: cbc.CodeMap{KeyFacturaE: "02"},
	},
	{
		Key: CorrectionKeyIssueDate,
		Name: i18n.String{
			i18n.EN: "Issue date",
			i18n.ES: "Fecha expedición",
		},
		Map: cbc.CodeMap{KeyFacturaE: "03"},
	},
	{
		Key: CorrectionKeySupplierName,
		Name: i18n.String{
			i18n.EN: "Name and surnames/Corporate name - Issuer (Sender)",
			i18n.ES: "Nombre y apellidos/Razón Social-Emisor",
		},
		Map: cbc.CodeMap{KeyFacturaE: "04"},
	},
	{
		Key: CorrectionKeyCustomerName,
		Name: i18n.String{
			i18n.EN: "Name and surnames/Corporate name - Receiver",
			i18n.ES: "Nombre y apellidos/Razón Social-Receptor",
		},
		Map: cbc.CodeMap{KeyFacturaE: "05"},
	},
	{
		Key: CorrectionKeySupplierTaxID,
		Name: i18n.String{
			i18n.EN: "Issuer's Tax Identification Number",
			i18n.ES: "Identificación fiscal Emisor/obligado",
		},
		Map: cbc.CodeMap{KeyFacturaE: "06"},
	},
	{
		Key: CorrectionKeyCustomerTaxID,
		Name: i18n.String{
			i18n.EN: "Receiver's Tax Identification Number",
			i18n.ES: "Identificación fiscal Receptor",
		},
		Map: cbc.CodeMap{KeyFacturaE: "07"},
	},
	{
		Key: CorrectionKeySupplierAddress,
		Name: i18n.String{
			i18n.EN: "Issuer's address",
			i18n.ES: "Domicilio Emisor/Obligado",
		},
		Map: cbc.CodeMap{KeyFacturaE: "08"},
	},
	{
		Key: CorrectionKeyCustomerAddress,
		Name: i18n.String{
			i18n.EN: "Receiver's address",
			i18n.ES: "Domicilio Receptor",
		},
		Map: cbc.CodeMap{KeyFacturaE: "09"},
	},
	{
		Key: CorrectionKeyLine,
		Name: i18n.String{
			i18n.EN: "Item line",
			i18n.ES: "Detalle Operación",
		},
		Map: cbc.CodeMap{KeyFacturaE: "10"},
	},
	{
		Key: CorrectionKeyTaxRate,
		Name: i18n.String{
			i18n.EN: "Applicable Tax Rate",
			i18n.ES: "Porcentaje impositivo a aplicar",
		},
		Map: cbc.CodeMap{KeyFacturaE: "11"},
	},
	{
		Key: CorrectionKeyTaxAmount,
		Name: i18n.String{
			i18n.EN: "Applicable Tax Amount",
			i18n.ES: "Cuota tributaria a aplicar",
		},
		Map: cbc.CodeMap{KeyFacturaE: "12"},
	},
	{
		Key: CorrectionKeyPeriod,
		Name: i18n.String{
			i18n.EN: "Applicable Date/Period",
			i18n.ES: "Fecha/Periodo a aplicar",
		},
		Map: cbc.CodeMap{KeyFacturaE: "13"},
	},
	{
		Key: CorrectionKeyType,
		Name: i18n.String{
			i18n.EN: "Invoice Class",
			i18n.ES: "Clase de factura",
		},
		Map: cbc.CodeMap{KeyFacturaE: "14"},
	},
	{
		Key: CorrectionKeyLegalDetails,
		Name: i18n.String{
			i18n.EN: "Legal literals",
			i18n.ES: "Literales legales",
		},
		Map: cbc.CodeMap{KeyFacturaE: "15"},
	},
	{
		Key: CorrectionKeyTaxBase,
		Name: i18n.String{
			i18n.EN: "Taxable Base",
			i18n.ES: "Base imponible",
		},
		Map: cbc.CodeMap{KeyFacturaE: "16"},
	},
	{
		Key: CorrectionKeyTax,
		Name: i18n.String{
			i18n.EN: "Calculation of tax outputs",
			i18n.ES: "Cálculo de cuotas repercutidas",
		},
		Map: cbc.CodeMap{KeyFacturaE: "80"},
	},
	{
		Key: CorrectionKeyTaxRetained,
		Name: i18n.String{
			i18n.EN: "Calculation of tax inputs",
			i18n.ES: "Cálculo de cuotas retenidas",
		},
		Map: cbc.CodeMap{KeyFacturaE: "81"},
	},
	{
		Key: CorrectionKeyRefund,
		Name: i18n.String{
			i18n.EN: "Taxable Base modified due to return of packages and packaging materials",
			i18n.ES: "Base imponible modificada por devolución de envases / embalajes",
		},
		Map: cbc.CodeMap{KeyFacturaE: "82"},
	},
	{
		Key: CorrectionKeyDiscount,
		Name: i18n.String{
			i18n.EN: "Taxable Base modified due to discounts and rebates",
			i18n.ES: "Base imponible modificada por descuentos y bonificaciones",
		},
		Map: cbc.CodeMap{KeyFacturaE: "83"},
	},
	{
		Key: CorrectionKeyJudicial,
		Name: i18n.String{
			i18n.EN: "Taxable Base modified due to firm court ruling or administrative decision",
			i18n.ES: "Base imponible modificada por resolución firme, judicial o administrativa",
		},
		Map: cbc.CodeMap{KeyFacturaE: "84"},
	},
	{
		Key: CorrectionKeyInsolvency,
		Name: i18n.String{
			i18n.EN: "Taxable Base modified due to unpaid outputs where there is a judgement opening insolvency proceedings",
			i18n.ES: "Base imponible modificada cuotas repercutidas no satisfechas. Auto de declaración de concurso",
		},
		Map: cbc.CodeMap{KeyFacturaE: "85"},
	},
}

var correctionMethodList = []*tax.KeyDefinition{
	{
		Key: CorrectionMethodKeyComplete,
		Name: i18n.String{
			i18n.EN: "Complete",
			i18n.ES: "Rectificaticón íntegra",
		},
		Map: cbc.CodeMap{KeyFacturaE: "01"},
	},
	{
		Key: CorrectionMethodKeyPartial,
		Name: i18n.String{
			i18n.EN: "Corrected items only",
			i18n.ES: "Rectificación por diferencias",
		},
		Map: cbc.CodeMap{KeyFacturaE: "02"},
	},
	{
		Key: CorrectionMethodKeyDiscount,
		Name: i18n.String{
			i18n.EN: "Bulk deal in a given period",
			i18n.ES: "Rectificación por descuento por volumen de operaciones durante un periodo",
		},
		Map: cbc.CodeMap{KeyFacturaE: "03"},
	},
	{
		Key: CorrectionMethodKeyAuthorized,
		Name: i18n.String{
			i18n.EN: "Authorized by the Tax Agency",
			i18n.ES: "Autorizadas por la Agencia Tributaria",
		},
		Map: cbc.CodeMap{KeyFacturaE: "04"},
	},
}

func correctionChangeKeys() []interface{} {
	keys := make([]interface{}, len(correctionChangesList))
	i := 0
	for _, v := range correctionChangesList {
		keys[i] = v.Key
		i++
	}
	return keys
}

func correctionMethodKeys() []interface{} {
	keys := make([]interface{}, len(correctionMethodList))
	i := 0
	for _, v := range correctionMethodList {
		keys[i] = v.Key
		i++
	}
	return keys
}

var isValidCorrectionChangeKey = validation.In(correctionChangeKeys()...)

var isValidCorrectionMethodKey = validation.In(correctionMethodKeys()...)
