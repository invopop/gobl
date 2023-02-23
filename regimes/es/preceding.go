package es

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// List of correction codes that are supported by the Spanish regime.
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

const (
	CorrectionMethodKeyComplete   cbc.Key = "complete"   // everything has changed
	CorrectionMethodKeyPartial    cbc.Key = "partial"    // only differences corrected
	CorrectionMethodKeyDiscount   cbc.Key = "discount"   // deducted from future invoices
	CorrectionMethodKeyAuthorized cbc.Key = "authorized" // Permitted by tax agency
)

// CorrectionMap maps GOBL Correction Codes to reason models acceptable for
// spanish invoices.
var correctionList = []*tax.KeyDefinition{
	{
		Key:  CorrectionKeyCode,
		Code: "01",
		Desc: i18n.String{
			i18n.EN: "Invoice code",
			i18n.ES: "Número de la factura",
		},
	},
	{
		Key:  CorrectionKeySeries,
		Code: "02",
		Desc: i18n.String{
			i18n.EN: "Invoice series",
			i18n.ES: "Serie de la factura",
		},
	},
	{
		Key:  CorrectionKeyIssueDate,
		Code: "03",
		Desc: i18n.String{
			i18n.EN: "Issue date",
			i18n.ES: "Fecha expedición",
		},
	},
	{
		Key:  CorrectionKeySupplierName,
		Code: "04",
		Desc: i18n.String{
			i18n.EN: "Name and surnames/Corporate name – Issuer (Sender)",
			i18n.ES: "Nombre y apellidos/Razón Social-Emisor",
		},
	},
	{
		Key:  CorrectionKeyCustomerName,
		Code: "05",
		Desc: i18n.String{
			i18n.EN: "Name and surnames/Corporate name - Receiver",
			i18n.ES: "Nombre y apellidos/Razón Social-Receptor",
		},
	},
	{
		Key:  CorrectionKeySupplierTaxID,
		Code: "06",
		Desc: i18n.String{
			i18n.EN: "Issuer's Tax Identification Number",
			i18n.ES: "Identificación fiscal Emisor/obligado",
		},
	},
	{
		Key:  CorrectionKeyCustomerTaxID,
		Code: "07",
		Desc: i18n.String{
			i18n.EN: "Receiver's Tax Identification Number",
			i18n.ES: "Identificación fiscal Receptor",
		},
	},
	{
		Key:  CorrectionKeySupplierAddress,
		Code: "08",
		Desc: i18n.String{
			i18n.EN: "Issuer's address",
			i18n.ES: "Domicilio Emisor/Obligado",
		},
	},
	{
		Key:  CorrectionKeyCustomerAddress,
		Code: "09",
		Desc: i18n.String{
			i18n.EN: "Receiver's address",
			i18n.ES: "Domicilio Receptor",
		},
	},
	{
		Key:  CorrectionKeyLine,
		Code: "10",
		Desc: i18n.String{
			i18n.EN: "Item line",
			i18n.ES: "Detalle Operación",
		},
	},
	{
		Key:  CorrectionKeyTaxRate,
		Code: "11",
		Desc: i18n.String{
			i18n.EN: "Applicable Tax Rate",
			i18n.ES: "Porcentaje impositivo a aplicar",
		},
	},
	{
		Key:  CorrectionKeyTaxAmount,
		Code: "12",
		Desc: i18n.String{
			i18n.EN: "Applicable Tax Amount",
			i18n.ES: "Cuota tributaria a aplicar",
		},
	},
	{
		Key:  CorrectionKeyPeriod,
		Code: "13",
		Desc: i18n.String{
			i18n.EN: "Applicable Date/Period",
			i18n.ES: "Fecha/Periodo a aplicar",
		},
	},
	{
		Key:  CorrectionKeyType,
		Code: "14",
		Desc: i18n.String{
			i18n.EN: "Invoice Class",
			i18n.ES: "Clase de factura",
		},
	},
	{
		Key:  CorrectionKeyLegalDetails,
		Code: "15",
		Desc: i18n.String{
			i18n.EN: "Legal literals",
			i18n.ES: "Literales legales",
		},
	},
	{
		Key:  CorrectionKeyTaxBase,
		Code: "16",
		Desc: i18n.String{
			i18n.EN: "Taxable Base",
			i18n.ES: "Base imponible",
		},
	},
	{
		Key:  CorrectionKeyTax,
		Code: "80",
		Desc: i18n.String{
			i18n.EN: "Calculation of tax outputs",
			i18n.ES: "Cálculo de cuotas repercutidas",
		},
	},
	{
		Key:  CorrectionKeyTaxRetained,
		Code: "81",
		Desc: i18n.String{
			i18n.EN: "Calculation of tax inputs",
			i18n.ES: "Cálculo de cuotas retenidas",
		},
	},
	{
		Key:  CorrectionKeyRefund,
		Code: "82",
		Desc: i18n.String{
			i18n.EN: "Taxable Base modified due to return of packages and packaging materials",
			i18n.ES: "Base imponible modificada por devolución de envases / embalajes",
		},
	},
	{
		Key:  CorrectionKeyDiscount,
		Code: "83",
		Desc: i18n.String{
			i18n.EN: "Taxable Base modified due to discounts and rebates",
			i18n.ES: "Base imponible modificada por descuentos y bonificaciones",
		},
	},
	{
		Key:  CorrectionKeyJudicial,
		Code: "84",
		Desc: i18n.String{
			i18n.EN: "Taxable Base modified due to firm court ruling or administrative decision",
			i18n.ES: "Base imponible modificada por resolución firme, judicial o administrativa",
		},
	},
	{
		Key:  CorrectionKeyInsolvency,
		Code: "85",
		Desc: i18n.String{
			i18n.EN: "Taxable Base modified due to unpaid outputs where there is a judgement opening insolvency proceedings",
			i18n.ES: "Base imponible modificada cuotas repercutidas no satisfechas. Auto de declaración de concurso",
		},
	},
}

var correctionMethodList = []*tax.KeyDefinition{
	{
		Key:  CorrectionMethodKeyComplete,
		Code: "01",
		Desc: i18n.String{
			i18n.EN: "Complete",
			i18n.ES: "Rectificaticón íntegra",
		},
	},
	{
		Key:  CorrectionMethodKeyPartial,
		Code: "02",
		Desc: i18n.String{
			i18n.EN: "Corrected items only",
			i18n.ES: "Rectificación por diferencias",
		},
	},
	{
		Key:  CorrectionMethodKeyDiscount,
		Code: "03",
		Desc: i18n.String{
			i18n.EN: "Bulk deal in a given period",
			i18n.ES: "Rectificación por descuento por volumen de operaciones durante un periodo",
		},
	},
	{
		Key:  CorrectionMethodKeyAuthorized,
		Code: "04",
		Desc: i18n.String{
			i18n.EN: "Authorized by the Tax Agency",
			i18n.ES: "Autorizadas por la Agencia Tributaria",
		},
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
