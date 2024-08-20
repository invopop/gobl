// Package gr provides the tax region definition for Greece.
package gr

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegime(New())
}

// Custom keys used typically in meta or codes information.
const (
	KeyIAPRPaymentMethod cbc.Key = "iapr-payment-method"
)

// Official IAPR codes to include in stamps.
const (
	StampIAPRQR       cbc.Key = "iapr-qr"
	StampIAPRMark     cbc.Key = "iapr-mark"
	StampIAPRHash     cbc.Key = "iapr-hash"
	StampIAPRUID      cbc.Key = "iapr-uid"
	StampIAPRProvider cbc.Key = "iapr-provider"
)

// New provides the tax region definition
func New() *tax.Regime {
	return &tax.Regime{
		Country: "EL",
		AltCountryCodes: []l10n.Code{
			"GR", // regular ISO code
		},
		Currency: currency.EUR,
		Name: i18n.String{
			i18n.EN: "Greece",
			i18n.EL: "Ελλάδα",
		},
		TimeZone:         "Europe/Athens",
		Tags:             invoiceTags,
		Scenarios:        scenarios,
		Corrections:      corrections,
		Validator:        Validate,
		Calculator:       Calculate,
		Categories:       taxCategories,
		PaymentMeansKeys: paymentMeansKeys,
		Extensions:       extensionKeys,
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *tax.Combo:
		return validateTaxCombo(obj)
	case *org.Address:
		return validateAddress(obj)
	}
	return nil
}

// Calculate will attempt to clean the object passed to it.
func Calculate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return normalizeTaxIdentity(obj)
	case *tax.Combo:
		return normalizeTaxCombo(obj)
	}
	return nil
}
