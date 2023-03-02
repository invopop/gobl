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
var correctionList = []*tax.KeyDefinition{
	{
		Key: CorrectionKeyCode,
		Desc: i18n.String{
			i18n.EN: "Invoice code",
			i18n.ES: "Número de la factura",
		},
		Meta: cbc.Meta{KeyFacturaE: "01"},
	},
	{
		Key: CorrectionKeySeries,
		Desc: i18n.String{
			i18n.EN: "Invoice series",
			i18n.ES: "Serie de la factura",
		},
		Meta: cbc.Meta{KeyFacturaE: "02"},
	},
	{
		Key: CorrectionKeyIssueDate,
		Desc: i18n.String{
			i18n.EN: "Issue date",
			i18n.ES: "Fecha expedición",
		},
		Meta: cbc.Meta{KeyFacturaE: "03"},
	},
	{
		Key: CorrectionKeySupplierName,
		Desc: i18n.String{
			i18n.EN: "Name and surnames/Corporate name – Issuer (Sender)",
			i18n.ES: "Nombre y apellidos/Razón Social-Emisor",
		},
		Meta: cbc.Meta{KeyFacturaE: "04"},
	},
	{
		Key: CorrectionKeyCustomerName,
		Desc: i18n.String{
			i18n.EN: "Name and surnames/Corporate name - Receiver",
			i18n.ES: "Nombre y apellidos/Razón Social-Receptor",
		},
		Meta: cbc.Meta{KeyFacturaE: "05"},
	},
	{
		Key: CorrectionKeySupplierTaxID,
		Desc: i18n.String{
			i18n.EN: "Issuer's Tax Identification Number",
			i18n.ES: "Identificación fiscal Emisor/obligado",
		},
		Meta: cbc.Meta{KeyFacturaE: "06"},
	},
	{
		Key: CorrectionKeyCustomerTaxID,
		Desc: i18n.String{
			i18n.EN: "Receiver's Tax Identification Number",
			i18n.ES: "Identificación fiscal Receptor",
		},
		Meta: cbc.Meta{KeyFacturaE: "07"},
	},
	{
		Key: CorrectionKeySupplierAddress,
		Desc: i18n.String{
			i18n.EN: "Issuer's address",
			i18n.ES: "Domicilio Emisor/Obligado",
		},
		Meta: cbc.Meta{KeyFacturaE: "08"},
	},
	{
		Key: CorrectionKeyCustomerAddress,
		Desc: i18n.String{
			i18n.EN: "Receiver's address",
			i18n.ES: "Domicilio Receptor",
		},
		Meta: cbc.Meta{KeyFacturaE: "09"},
	},
	{
		Key: CorrectionKeyLine,
		Desc: i18n.String{
			i18n.EN: "Item line",
			i18n.ES: "Detalle Operación",
		},
		Meta: cbc.Meta{KeyFacturaE: "10"},
	},
	{
		Key: CorrectionKeyTaxRate,
		Desc: i18n.String{
			i18n.EN: "Applicable Tax Rate",
			i18n.ES: "Porcentaje impositivo a aplicar",
		},
		Meta: cbc.Meta{KeyFacturaE: "11"},
	},
	{
		Key: CorrectionKeyTaxAmount,
		Desc: i18n.String{
			i18n.EN: "Applicable Tax Amount",
			i18n.ES: "Cuota tributaria a aplicar",
		},
		Meta: cbc.Meta{KeyFacturaE: "12"},
	},
	{
		Key: CorrectionKeyPeriod,
		Desc: i18n.String{
			i18n.EN: "Applicable Date/Period",
			i18n.ES: "Fecha/Periodo a aplicar",
		},
		Meta: cbc.Meta{KeyFacturaE: "13"},
	},
	{
		Key: CorrectionKeyType,
		Desc: i18n.String{
			i18n.EN: "Invoice Class",
			i18n.ES: "Clase de factura",
		},
		Meta: cbc.Meta{KeyFacturaE: "14"},
	},
	{
		Key: CorrectionKeyLegalDetails,
		Desc: i18n.String{
			i18n.EN: "Legal literals",
			i18n.ES: "Literales legales",
		},
		Meta: cbc.Meta{KeyFacturaE: "15"},
	},
	{
		Key: CorrectionKeyTaxBase,
		Desc: i18n.String{
			i18n.EN: "Taxable Base",
			i18n.ES: "Base imponible",
		},
		Meta: cbc.Meta{KeyFacturaE: "16"},
	},
	{
		Key: CorrectionKeyTax,
		Desc: i18n.String{
			i18n.EN: "Calculation of tax outputs",
			i18n.ES: "Cálculo de cuotas repercutidas",
		},
		Meta: cbc.Meta{KeyFacturaE: "80"},
	},
	{
		Key: CorrectionKeyTaxRetained,
		Desc: i18n.String{
			i18n.EN: "Calculation of tax inputs",
			i18n.ES: "Cálculo de cuotas retenidas",
		},
		Meta: cbc.Meta{KeyFacturaE: "81"},
	},
	{
		Key: CorrectionKeyRefund,
		Desc: i18n.String{
			i18n.EN: "Taxable Base modified due to return of packages and packaging materials",
			i18n.ES: "Base imponible modificada por devolución de envases / embalajes",
		},
		Meta: cbc.Meta{KeyFacturaE: "82"},
	},
	{
		Key: CorrectionKeyDiscount,
		Desc: i18n.String{
			i18n.EN: "Taxable Base modified due to discounts and rebates",
			i18n.ES: "Base imponible modificada por descuentos y bonificaciones",
		},
		Meta: cbc.Meta{KeyFacturaE: "83"},
	},
	{
		Key: CorrectionKeyJudicial,
		Desc: i18n.String{
			i18n.EN: "Taxable Base modified due to firm court ruling or administrative decision",
			i18n.ES: "Base imponible modificada por resolución firme, judicial o administrativa",
		},
		Meta: cbc.Meta{KeyFacturaE: "84"},
	},
	{
		Key: CorrectionKeyInsolvency,
		Desc: i18n.String{
			i18n.EN: "Taxable Base modified due to unpaid outputs where there is a judgement opening insolvency proceedings",
			i18n.ES: "Base imponible modificada cuotas repercutidas no satisfechas. Auto de declaración de concurso",
		},
		Meta: cbc.Meta{KeyFacturaE: "85"},
	},
}

var correctionMethodList = []*tax.KeyDefinition{
	{
		Key: CorrectionMethodKeyComplete,
		Desc: i18n.String{
			i18n.EN: "Complete",
			i18n.ES: "Rectificaticón íntegra",
		},
		Meta: cbc.Meta{KeyFacturaE: "01"},
	},
	{
		Key: CorrectionMethodKeyPartial,
		Desc: i18n.String{
			i18n.EN: "Corrected items only",
			i18n.ES: "Rectificación por diferencias",
		},
		Meta: cbc.Meta{KeyFacturaE: "02"},
	},
	{
		Key: CorrectionMethodKeyDiscount,
		Desc: i18n.String{
			i18n.EN: "Bulk deal in a given period",
			i18n.ES: "Rectificación por descuento por volumen de operaciones durante un periodo",
		},
		Meta: cbc.Meta{KeyFacturaE: "03"},
	},
	{
		Key: CorrectionMethodKeyAuthorized,
		Desc: i18n.String{
			i18n.EN: "Authorized by the Tax Agency",
			i18n.ES: "Autorizadas por la Agencia Tributaria",
		},
		Meta: cbc.Meta{KeyFacturaE: "04"},
	},
}

func correctionKeys() []interface{} {
	keys := make([]interface{}, len(correctionList))
	i := 0
	for _, v := range correctionList {
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

var isValidCorrectionKey = validation.In(correctionKeys()...)

var isValidCorrectionMethodKey = validation.In(correctionMethodKeys()...)
