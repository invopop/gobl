package es

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/i18n"
)

// CorrectionReason defines expected correction reasons in Spanish invoices.
type CorrectionReason struct {
	Code string      `json:"code"`
	Desc i18n.String `json:"desc"`
}

// CorrectionMethod is used to define a correction method considered acceptable.
type CorrectionMethod struct {
	Code string      `json:"code"`
	Desc i18n.String `json:"desc"`
}

// CorrectionReasonMap maps GOBL Correction Codes to reason models acceptable for
// spanish invoices.
var CorrectionReasonMap = map[bill.CorrectionKey]*CorrectionReason{
	bill.CodeCorrectionKey: {
		Code: "01",
		Desc: i18n.String{
			i18n.EN: "Invoice Number",
			i18n.ES: "Número de la factura",
		},
	},
	bill.SeriesCorrectionKey: {
		Code: "02",
		Desc: i18n.String{
			i18n.EN: "Invoice serial number",
			i18n.ES: "Serie de la factura",
		},
	},
	bill.IssueDateCorrectionKey: {
		Code: "03",
		Desc: i18n.String{
			i18n.EN: "Issue date",
			i18n.ES: "Fecha expedición",
		},
	},
	bill.SupplierNameCorrectionKey: {
		Code: "04",
		Desc: i18n.String{
			i18n.EN: "Name and surnames/Corporate name – Issuer (Sender)",
			i18n.ES: "Nombre y apellidos/Razón Social-Emisor",
		},
	},
	bill.CustomerNameCorrectionKey: {
		Code: "05",
		Desc: i18n.String{
			i18n.EN: "Name and surnames/Corporate name - Receiver",
			i18n.ES: "Nombre y apellidos/Razón Social-Receptor",
		},
	},
	bill.SupplierTaxIDCorrectionKey: {
		Code: "06",
		Desc: i18n.String{
			i18n.EN: "Issuer's Tax Identification Number",
			i18n.ES: "Identificación fiscal Emisor/obligado",
		},
	},
	bill.CustomerTaxIDCorrectionKey: {
		Code: "07",
		Desc: i18n.String{
			i18n.EN: "Receiver's Tax Identification Number",
			i18n.ES: "Identificación fiscal Receptor",
		},
	},
	bill.SupplierAddressCorrectionKey: {
		Code: "08",
		Desc: i18n.String{
			i18n.EN: "Issuer's address",
			i18n.ES: "Domicilio Emisor/Obligado",
		},
	},
	bill.CustomerAddressCorrectionKey: {
		Code: "09",
		Desc: i18n.String{
			i18n.EN: "Receiver's address",
			i18n.ES: "Domicilio Receptor",
		},
	},
	bill.LineCorrectionKey: {
		Code: "10",
		Desc: i18n.String{
			i18n.EN: "Item line",
			i18n.ES: "Detalle Operación",
		},
	},
	bill.TaxRateCorrectionKey: {
		Code: "11",
		Desc: i18n.String{
			i18n.EN: "Applicable Tax Rate",
			i18n.ES: "Porcentaje impositivo a aplicar",
		},
	},
	bill.TaxAmountCorrectionKey: {
		Code: "12",
		Desc: i18n.String{
			i18n.EN: "Applicable Tax Amount",
			i18n.ES: "Cuota tributaria a aplicar",
		},
	},
	bill.PeriodCorrectionKey: {
		Code: "13",
		Desc: i18n.String{
			i18n.EN: "Applicable Date/Period",
			i18n.ES: "Fecha/Periodo a aplicar",
		},
	},
	bill.TypeCorrectionKey: {
		Code: "14",
		Desc: i18n.String{
			i18n.EN: "Invoice Class",
			i18n.ES: "Clase de factura",
		},
	},
	bill.LegalDetailsCorrectionKey: {
		Code: "15",
		Desc: i18n.String{
			i18n.EN: "Legal literals",
			i18n.ES: "Literales legales",
		},
	},
	bill.TaxBaseCorrectionKey: {
		Code: "16",
		Desc: i18n.String{
			i18n.EN: "Taxable Base",
			i18n.ES: "Base imponible",
		},
	},
	bill.TaxCorrectionKey: {
		Code: "80",
		Desc: i18n.String{
			i18n.EN: "Calculation of tax outputs",
			i18n.ES: "Cálculo de cuotas repercutidas",
		},
	},
	bill.TaxRetainedCorrectionKey: {
		Code: "81",
		Desc: i18n.String{
			i18n.EN: "Calculation of tax inputs",
			i18n.ES: "Cálculo de cuotas retenidas",
		},
	},
	bill.RefundCorrectionKey: {
		Code: "82",
		Desc: i18n.String{
			i18n.EN: "Taxable Base modified due to return of packages and packaging materials",
			i18n.ES: "Base imponible modificada por devolución de envases / embalajes",
		},
	},
	bill.DiscountCorrectionKey: {
		Code: "83",
		Desc: i18n.String{
			i18n.EN: "Taxable Base modified due to discounts and rebates",
			i18n.ES: "Base imponible modificada por descuentos y bonificaciones",
		},
	},
	bill.JudicialCorrectionKey: {
		Code: "84",
		Desc: i18n.String{
			i18n.EN: "Taxable Base modified due to firm court ruling or administrative decision",
			i18n.ES: "Base imponible modificada por resolución firme, judicial o administrativa",
		},
	},
	bill.InsolvencyCorrectionKey: {
		Code: "85",
		Desc: i18n.String{
			i18n.EN: "Taxable Base modified due to unpaid outputs where there is a judgement opening insolvency proceedings",
			i18n.ES: "Base imponible modificada cuotas repercutidas no satisfechas. Auto de declaración de concurso",
		},
	},
}

// CorrectionMethodMap defines the codes and texts expected by Spanish electronic invoices
// for the types of corrections being made to an invoice.
var CorrectionMethodMap = map[bill.CorrectionMethodKey]*CorrectionMethod{
	bill.CompleteCorrectionMethodKey: {
		Code: "01",
		Desc: i18n.String{
			i18n.EN: "Complete",
			i18n.ES: "Rectificaticón íntegra",
		},
	},
	bill.PartialCorrectionMethodKey: {
		Code: "02",
		Desc: i18n.String{
			i18n.EN: "Corrected items only",
			i18n.ES: "Rectificación por diferencias",
		},
	},
	bill.DiscountCorrectionMethodKey: {
		Code: "03",
		Desc: i18n.String{
			i18n.EN: "Bulk deal in a given period",
			i18n.ES: "Rectificación por descuento por volumen de operaciones durante un periodo",
		},
	},
	bill.AuthorizedCorrectionMethodKey: {
		Code: "04",
		Desc: i18n.String{
			i18n.EN: "Authorized by the Tax Agency",
			i18n.ES: "Autorizadas por la Agencia Tributaria",
		},
	},
}

func correctionReasonKeys() []interface{} {
	keys := make([]interface{}, len(CorrectionReasonMap))
	i := 0
	for k := range CorrectionReasonMap {
		keys[i] = k
		i++
	}
	return keys
}

func correctionMethodKeys() []interface{} {
	keys := make([]interface{}, len(CorrectionMethodMap))
	i := 0
	for k := range CorrectionMethodMap {
		keys[i] = k
		i++
	}
	return keys
}
