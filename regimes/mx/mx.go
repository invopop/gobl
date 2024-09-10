// Package mx provides the Mexican tax regime.
package mx

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegime(New())

	// MX GOBL Schema Complements
	schema.Register(schema.GOBL.Add("regimes/mx"),
		FuelAccountBalance{},
		FoodVouchers{},
	)
}

// Custom keys used typically in meta or codes information.
const (
	KeySATFormaPago    cbc.Key = "sat-forma-pago"    // for mapping to c_FormaPago’s codes
	KeySATTipoRelacion cbc.Key = "sat-tipo-relacion" // for mapping to c_TipoRelacion’s codes
	KeySATImpuesto     cbc.Key = "sat-impuesto"      // for mapping to c_Impuesto’s codes
)

// Official SAT and CFDI codes to include in stamps.
const (
	StampCFDISignature  cbc.Key = "cfdi-sig"         // Signature - Sello Digital del CFDI
	StampCFDISerial     cbc.Key = "cfdi-serial"      // Cert Serial - Número de Certificado del CFDI
	StampSATSignature   cbc.Key = "sat-sig"          // Signature - Sello Digital del SAT (optional)
	StampSATSerial      cbc.Key = "sat-serial"       // Cert Serial - Número de Certificado SAT
	StampSATTimestamp   cbc.Key = "sat-timestamp"    // Timestamp - Fecha y hora de certificación del SAT
	StampSATUUID        cbc.Key = "sat-uuid"         // Folio Fiscal
	StampSATURL         cbc.Key = "sat-url"          // URL QR Code
	StampSATProviderRFC cbc.Key = "sat-provider-rfc" // Provider RFC - RFC del Proveedor de Certificación
	StampSATChain       cbc.Key = "sat-chain"        // Cadena original del complemento de certificación digital del SAT
)

// New provides the tax region definition
func New() *tax.Regime {
	return &tax.Regime{
		Country:  "MX",
		Currency: currency.MXN,
		Name: i18n.String{
			i18n.EN: "Mexico",
			i18n.ES: "México",
		},
		TimeZone:         "America/Mexico_City",
		Validator:        Validate,
		Calculator:       Calculate,
		Tags:             common.InvoiceTags(),
		PaymentMeansKeys: paymentMeansKeyDefinitions, // pay.go
		Extensions:       extensionKeys,              // extensions.go
		Scenarios:        scenarios,                  // scenarios.go
		Categories:       taxCategories,              // categories.go
		Corrections:      correctionDefinitions,      // corrections.go
	}
}

// Validate validates a document against the tax regime.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *org.Item:
		return validateItem(obj)
	}
	return nil
}

// Calculate performs regime specific calculations.
func Calculate(doc interface{}) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return normalizeInvoice(obj)
	case *tax.Identity:
		return tax.NormalizeIdentity(obj)
	case *org.Party:
		return normalizeParty(obj)
	case *org.Item:
		return normalizeItem(obj)
	}
	return nil
}
