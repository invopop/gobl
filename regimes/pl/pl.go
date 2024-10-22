// Package pl provides the Polish tax regime.
package pl

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

// KSeF official codes to include.
const (
	StampProviderKSeFID   cbc.Key = "ksef-id"
	StampProviderKSeFHash cbc.Key = "ksef-hash"
	StampProviderKSeFQR   cbc.Key = "ksef-qr"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// Custom keys used typically in meta or codes information.
const (
	KeyFAVATPaymentType cbc.Key = "favat-forma-platnosci" // for mapping to TFormaPlatnosci's codes
	KeyFAVATInvoiceType cbc.Key = "favat-rodzaj-faktury"  // for mapping to TRodzajFaktury's codes
)

// New instantiates a new Polish regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:  "PL",
		Currency: currency.PLN,
		Name: i18n.String{
			i18n.EN: "Poland",
			i18n.PL: "Polska",
		},
		TimeZone: "Europe/Warsaw",
		// ChargeKeys:       chargeKeyDefinitions,       // charges.go
		PaymentMeansKeys: paymentMeansKeyDefinitions, // pay.go
		Extensions:       extensionKeys,              // extensions.go
		Tags: []*tax.TagSet{
			common.InvoiceTags().Merge(invoiceTags),
		},
		Scenarios:    scenarios,              // scenarios.go
		IdentityKeys: identityKeyDefinitions, // identities.go
		Validator:    Validate,
		Normalizer:   Normalize,
		Categories:   taxCategories, // tax_categories.go
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
				ReasonRequired: true,
				Stamps: []cbc.Key{
					StampProviderKSeFID,
				},
				Extensions: []cbc.Key{
					ExtKeyKSeFEffectiveDate,
				},
			},
		},
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *bill.Invoice:
		return validateInvoice(obj)
	case *org.Identity:
		return validateTaxNumber(obj)
		// case *pay.Instructions:
		// 	return validatePayInstructions(obj)
		// case *pay.Advance:
		// 	return validatePayAdvance(obj)
	}
	return nil
}

// Normalize will perform any regime specific normalizations.
func Normalize(doc interface{}) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
