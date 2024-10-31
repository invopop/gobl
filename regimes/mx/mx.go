// Package mx provides the Mexican tax regime.
package mx

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/common"
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
		Country:  "MX",
		Currency: currency.MXN,
		Name: i18n.String{
			i18n.EN: "Mexico",
			i18n.ES: "México",
		},
		TimeZone:   "America/Mexico_City",
		Validator:  Validate,
		Normalizer: Normalize,
		Tags: []*tax.TagSet{
			common.InvoiceTags(),
		},
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
	case *org.Party:
		normalizeParty(obj)
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
