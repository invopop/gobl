// Package se provides a regime definition for Sweden.
package se

import (
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New instantiates a new Swedish regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "SE",
		Currency:  currency.SEK,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Sweden",
			i18n.SE: "Sverige",
		},
		TimeZone:   "Europe/Stockholm",
		// Identities: identityKeyDefinitions, // identities.go
		// Categories: taxCategories,          // tax_categories.go
		Scenarios: []*tax.ScenarioSet{
			bill.InvoiceScenarios(),
		},
		// Validator:  Validate,
		// Normalizer: Normalize,
	}
}
