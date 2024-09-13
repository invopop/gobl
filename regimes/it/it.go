// Package it provides the Italian tax regime.
package it

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegime(New())
}

// Keys used for meta data from external sources.
const (
	KeyFatturaPATipoDocumento     cbc.Key = "fatturapa-tipo-documento"
	KeyFatturaPATipoRitenuta      cbc.Key = "fatturapa-tipo-ritenuta"
	KeyFatturaPAModalitaPagamento cbc.Key = "fatturapa-modalita-pagamento"
)

// New instantiates a new Italian regime.
func New() *tax.Regime {
	return &tax.Regime{
		Country:  "IT",
		Currency: currency.EUR,
		Name: i18n.String{
			i18n.EN: "Italy",
			i18n.IT: "Italia",
		},
		TimeZone:         "Europe/Rome",
		ChargeKeys:       chargeKeyDefinitions,       // charges.go
		IdentityKeys:     identityKeyDefinitions,     // identities.go
		InboxKeys:        inboxKeyDefinitions,        // inboxes.go
		PaymentMeansKeys: paymentMeansKeyDefinitions, // pay.go
		Extensions:       extensionKeys,              // extensions.go
		Tags:             invoiceTags,
		Scenarios:        scenarios, // scenarios.go
		Validator:        Validate,
		Normalizer:       Normalize,
		Categories:       categories, // categories.go
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
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
	case *pay.Instructions:
		return validatePayInstructions(obj)
	case *org.Identity:
		return validateIdentity(obj)
	case *pay.Advance:
		return validatePayAdvance(obj)
	case *tax.Combo:
		return validateTaxCombo(obj)
	}
	return nil
}

// Normalize will perform any regime specific calculations.
func Normalize(doc interface{}) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeInvoice(obj)
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	case *org.Identity:
		normalizeIdentity(obj)
	}
}
