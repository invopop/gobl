// Package mx provides the Mexican tax regime.
package mx

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// Official SAT codes to include in stamps.
const (
	StampSATSignature   cbc.Key = "sat-sig"          // Signature - Sello Digital del SAT (optional)
	StampSATSerial      cbc.Key = "sat-serial"       // Cert Serial - Número de Certificado SAT
	StampSATTimestamp   cbc.Key = "sat-timestamp"    // Timestamp - Fecha y hora de certificación del SAT
	StampSATUUID        cbc.Key = "sat-uuid"         // Folio Fiscal
	StampSATURL         cbc.Key = "sat-url"          // URL QR Code
	StampSATProviderRFC cbc.Key = "sat-provider-rfc" // Provider RFC - RFC del Proveedor de Certificación
	StampSATChain       cbc.Key = "sat-chain"        // Cadena original del complemento de certificación digital del SAT
)

// Custom keys used typically in meta or codes information.
const (
	KeyFormaPago    cbc.Key = "sat-forma-pago"    // for mapping to c_FormaPago’s codes
	KeyTipoRelacion cbc.Key = "sat-tipo-relacion" // for mapping to c_TipoRelacion’s codes
	KeyImpuesto     cbc.Key = "sat-impuesto"      // for mapping to c_Impuesto’s codes
)

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "MX",
		Currency:  currency.MXN,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Mexico",
			i18n.ES: "México",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Mexico's tax system is administered by the SAT (Servicio de Administración
				Tributaria). Electronic invoicing through CFDI (Comprobante Fiscal Digital por
				Internet) version 4.0 is mandatory for all businesses.

				IVA (Impuesto al Valor Agregado) rates include a 16% general rate for most
				goods and services, 0% for food, medicine, and exports, and exempt categories
				for educational and medical services.

				Businesses are identified by their RFC (Registro Federal de Contribuyentes),
				a 12-character code for companies or 13-character code for individuals, which
				includes a date component and check digits. Every supplier and customer must
				be associated with a fiscal regime code (RegimenFiscal).

				CFDI invoices require specific fields including issue place (LugarExpedicion),
				CFDI use (UsoCFDI), payment method (MetodoPago distinguishing between fully
				paid PUE and pending PPD invoices), payment means (FormaPago), and
				product/service codes (ClaveProdServ) from the SAT catalog. For B2C sales,
				the simplified tag triggers use of the generic RFC code for final consumers.
				For foreign customers, their country and local tax code are mapped
				automatically. Invoices can include complements for fuel account balances and
				food vouchers among others.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("SAT - Anexo 20 (Invoice Format)"),
				URL:   "http://omawww.sat.gob.mx/tramitesyservicios/Paginas/anexo_20.htm",
			},
			{
				Title: i18n.NewString("SAT - CFDI 4.0 Filling Guide"),
				URL:   "http://omawww.sat.gob.mx/tramitesyservicios/Paginas/documentos/Anexo_20_Guia_de_llenado_CFDI.pdf",
			},
			{
				Title: i18n.NewString("SAT - Global CFDI 4.0 Filling Guide"),
				URL:   "http://omawww.sat.gob.mx/tramitesyservicios/Paginas/documentos/GuiallenadoCFDIglobal311221.pdf",
			},
		},
		TimeZone: "America/Mexico_City",
		Validator:   Validate,
		Normalizer:  Normalize,
		Categories:  taxCategories,
		Corrections: correctionDefinitions,
	}
}

// Normalize performs regime specific calculations.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeInvoice(obj)
	case *tax.Identity:
		NormalizeTaxIdentity(obj)
	}
}

// Validate validates a document against the tax regime.
func Validate(doc any) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return ValidateTaxIdentity(obj)
	}
	return nil
}
