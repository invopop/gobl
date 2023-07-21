// Package it provides the Italian tax regime.
package it

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegime(New())
}

// Keys used for meta data from external sources.
const (
	KeyFatturaPATipoDocumento     cbc.Key = "fatturapa-tipo-documento"
	KeyFatturaPARegimeFiscale     cbc.Key = "fatturapa-regime-fiscale"
	KeyFatturaPANatura            cbc.Key = "fatturapa-natura"
	KeyFatturaPATipoRitenuta      cbc.Key = "fatturapa-tipo-ritenuta"
	KeyFatturaPAModalitaPagamento cbc.Key = "fatturapa-modalita-pagamento"
	KeyFatturaPACausalePagamento  cbc.Key = "fatturapa-causale-pagamento"
)

// New instantiates a new Italian regime.
func New() *tax.Regime {
	return &tax.Regime{
		Country:  l10n.IT,
		Currency: currency.EUR,
		Name: i18n.String{
			i18n.EN: "Italy",
			i18n.IT: "Italia",
		},
		ChargeKeys:       chargeKeyDefinitions,       // charges.go
		PaymentMeansKeys: paymentMeansKeyDefinitions, // pay.go
		IdentityTypeKeys: taxIdentityTypeDefinitions, // tax_identity.go
		Tags:             invoiceTags,
		Scenarios:        scenarios, // scenarios.go
		Validator:        Validate,
		Normalizer:       Normalize,
		Zones:            zones,      // zones.go
		Categories:       categories, // categories.go
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
	case *pay.Advance:
		return validatePayAdvance(obj)
	}
	return nil
}

// Normalize will perform any regime specific calculations.
func Normalize(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return normalizeTaxIdentity(obj)
	}
	return nil
}
