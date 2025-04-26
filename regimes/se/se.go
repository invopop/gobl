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
		// Identities: identityKeyDefinitions,
		Categories: taxCategories,
		Scenarios: []*tax.ScenarioSet{
			bill.InvoiceScenarios(),
		},
		Validator:  Validate,
		Normalizer: Normalize,
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc any) error {
	switch obj := doc.(type) {
	// case *bill.Invoice:
	// 	return validateInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Normalize will perform any regime specific calculations.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
		// case *org.Identity:
		// 	normalizeOrgIdentity(obj)
	}
}
