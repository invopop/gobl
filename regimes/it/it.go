package it

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
)

// IdentityTypeFiscalCode represents the Italian Fiscal Code (Codice Fiscale).
// See https://en.wikipedia.org/wiki/Italian_fiscal_code. Every natural person
// has a fiscal code, and it is used to identify them for tax purposes. Not to
// be confused with the Italian VAT number (Partita IVA).
var IdentityTypeFiscalCode cbc.Code = "CF" // Codice Fiscale

// New instantiates a new Italian regime.
func New() *tax.Regime {
	return &tax.Regime{
		Country:  l10n.IT,
		Currency: "EUR",
		Name: i18n.String{
			i18n.EN: "Italy",
			i18n.IT: "Italia",
		},
		Validator:  Validate,
		Calculator: Calculate,
		Zones:      zones,         // zones.go
		Categories: taxCategories, // tax_categories.go
		Schemes:    schemes,       // schemes.go
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc interface{}) error {

	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Calculate will perform any regime specific calculations.
func Calculate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return normalizeTaxIdentity(obj)
	}
	return nil
}
