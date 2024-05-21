// Package fr provides the tax region definition for France.
package fr

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/regimes/common"
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
		Currency: currency.EUR,
		Name: i18n.String{
			i18n.EN: "France",
			i18n.FR: "La France",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				The French tax regime covers the basics.
			`),
		},
		TimeZone: "Europe/Paris",
		Tags:     common.InvoiceTags(),
		Scenarios: []*tax.ScenarioSet{
			invoiceScenarios,
		},
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				// France supports both corrective methods
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote, // Code 381
					bill.InvoiceTypeCorrective, // Code 384
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
	case *tax.Identity:
		return normalizeTaxIdentity(obj)
	}
	return nil
}
