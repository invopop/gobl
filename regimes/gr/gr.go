// Package gr provides the tax region definition for Greece.
package gr

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
	tax.RegisterRegimeDef(New())
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
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country: "EL",
		AltCountryCodes: []l10n.Code{
			"GR", // regular ISO code
		},
		Currency: currency.EUR,
		Name: i18n.String{
			i18n.EN: "Greece",
			i18n.EL: "Ελλάδα",
		},
		TimeZone:               "Europe/Athens",
		CalculatorRoundingRule: tax.CalculatorRoundThenSum,
		Tags: []*tax.TagSet{
			common.InvoiceTags().Merge(invoiceTags),
		},
		Scenarios:        scenarios,
		Corrections:      corrections,
		Validator:        Validate,
		Normalizer:       Normalize,
		Categories:       taxCategories,
		PaymentMeansKeys: paymentMeansKeys,
		Extensions:       extensionKeys,
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *org.Address:
		return validateAddress(obj)
	case *tax.Combo:
		return validateTaxCombo(obj)
	}
	return nil
}

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	}
}
