// Package pl provides the Polish tax regime.
package pl

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegime(New())
}

// Custom keys used typically in meta or codes information.
const (
	KeyFAVATPaymentType cbc.Key = "favat-forma-platnosci" // for mapping to TFormaPlatnosci's codes
	KeyFAVATInvoiceType cbc.Key = "favat-rodzaj-faktury"  // for mapping to TRodzajFaktury's codes
)

// New instantiates a new Polish regime.
func New() *tax.Regime {
	return &tax.Regime{
		Country:  l10n.PL,
		Currency: currency.PLN,
		Name: i18n.String{
			i18n.EN: "Poland",
			i18n.PL: "Polska",
		},
		TimeZone: "Europe/Warsaw",
		// ChargeKeys:       chargeKeyDefinitions,       // charges.go
		PaymentMeansKeys: paymentMeansKeyDefinitions, // pay.go
		Extensions:       extensionKeys,              // extensions.go
		Tags:             invoiceTags,
		Scenarios:        scenarios, // scenarios.go
		Validator:        Validate,
		// Calculator:       Calculate,
		Categories: taxCategories, // tax_categories.go
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *bill.Invoice:
		return validateInvoice(obj)
		// case *pay.Instructions:
		// 	return validatePayInstructions(obj)
		// case *pay.Advance:
		// 	return validatePayAdvance(obj)
	}
	return nil
}

// Calculate will perform any regime specific calculations.
// func Calculate(doc interface{}) error {
// 	switch obj := doc.(type) {
// 	case *tax.Identity:
// 		return normalizeTaxIdentity(obj)
// 	}
// 	return nil
// }
