// Package mx provides the Mexican tax regime.
package mx

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegime(New())
}

// Custom keys used typically in meta or codes information.
const (
	KeySATFormaPago         cbc.Key = "sat-forma-pago"          // for mapping to c_FormaPago’s codes
	KeySATTipoDeComprobante cbc.Key = "sat-tipo-de-comprobante" // for mapping to c_TipoDeComprobante’s codes
	KeySATTipoRelacion      cbc.Key = "sat-tipo-relacion"       // for mapping to c_TipoRelacion’s codes
)

// SAT official codes to include in stamps.
const (
	StampProviderSATUUID cbc.Key = "sat-uuid" // a.k.a. Folio Fiscal
)

// New provides the tax region definition
func New() *tax.Regime {
	return &tax.Regime{
		Country:  l10n.MX,
		Currency: currency.MXN,
		Name: i18n.String{
			i18n.EN: "Mexico",
			i18n.ES: "México",
		},
		TimeZone:         "America/Mexico_City",
		Validator:        Validate,
		Calculator:       Calculate,
		PaymentMeansKeys: paymentMeansKeyDefinitions, // pay.go
		Extensions:       extensionKeys,              // extensions.go
		Scenarios:        scenarios,                  // scenarios.go
		Categories:       taxCategories,              // categories.go
		Preceding:        precedingDefinitions,       // preceding.go
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
	case *tax.Identity:
		return common.NormalizeTaxIdentity(obj)
	case *org.Item:
		return normalizeItem(obj)
	}
	return nil
}
