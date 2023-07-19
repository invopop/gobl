// Package fr provides the tax region definition for France.
package fr

import (
	"context"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
)

// Identification keys used for additional codes not
// covered by the standard fields.
const (
	IdentityTypeSIREN cbc.Code = "SIREN" // SIREN is the main local tax code used in france, we use the normalized VAT version for the tax ID.
	IdentityTypeSIRET cbc.Code = "SIRET" // SIRET number combines the SIREN with a branch number.
	IdentityTypeRCS   cbc.Code = "RCS"   // Trade and Companies Register.
	IdentityTypeRM    cbc.Code = "RM"    // Directory of Traders.
	IdentityTypeNAF   cbc.Code = "NAF"   // Identifies the main branch of activity of the company or self-employed person.
)

func init() {
	tax.RegisterRegime(New())
}

// New provides the tax region definition
func New() *tax.Regime {
	return &tax.Regime{
		Country:  l10n.FR,
		Currency: "EUR",
		Name: i18n.String{
			i18n.EN: "France",
			i18n.FR: "La France",
		},
		Tags: invoiceTags,
		Scenarios: []*tax.ScenarioSet{
			invoiceScenarios,
		},
		Preceding: &tax.PrecedingDefinitions{
			// France supports both corrective methods
			Types: []cbc.Key{
				bill.InvoiceTypeCorrective, // Code 384
				bill.InvoiceTypeCreditNote, // Code 381
			},
		},
		Validator:  Validate,
		Calculator: Calculate,
		Categories: taxCategories,
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Calculate will attempt to clean the object passed to it.
func Calculate(_ context.Context, doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return normalizeTaxIdentity(obj)
	}
	return nil
}
