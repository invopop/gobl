// Package de provides the tax region definition for Germany.
package de

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
	tax.RegisterRegime(New())
}

// New provides the tax region definition
func New() *tax.Regime {
	return &tax.Regime{
		Country:  l10n.DE,
		Currency: currency.EUR,
		Name: i18n.String{
			i18n.EN: "Germany",
			i18n.DE: "Deutschland",
		},
		TimeZone: "Europe/Berlin",
		Tags:     common.InvoiceTags(),
		Scenarios: []*tax.ScenarioSet{
			invoiceScenarios,
		},
		IdentityKeys: identityKeyDefinitions, // identities.go
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				// Germany only supports credit notes to correct an invoice
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
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
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Calculate will attempt to clean the object passed to it.
func Calculate(doc interface{}) error {
	switch obj := doc.(type) {
	case *org.Identity:
		return normalizeIdentity(obj)
	case *tax.Identity:
		return common.NormalizeTaxIdentity(obj)
	}
	return nil
}
