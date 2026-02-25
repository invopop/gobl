// Package de provides the tax region definition for Germany.
package de

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "DE",
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Germany",
			i18n.DE: "Deutschland",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Germany's tax system is administered by the Federal Central Tax Office
				(Bundeszentralamt für Steuern, BZSt). As an EU member state, Germany follows
				the EU VAT Directive with locally adapted rates.

				VAT (Umsatzsteuer, USt) applies at standard and reduced rates. The reduced
				rate covers food, books, newspapers, public transport, and cultural events.

				Businesses are identified by their Umsatzsteuer-Identifikationsnummer (USt-IdNr)
				in the format DE followed by 9 digits for cross-border transactions, and by
				their Steuernummer (tax number) in regional formats for domestic purposes.

				Germany supports credit notes for invoice corrections. E-invoicing is
				progressively becoming mandatory, with XRechnung as the standard for B2G
				transactions and ZUGFeRD/Factur-X widely used for B2B.
			`),
		},
		TimeZone: "Europe/Berlin",
		Scenarios: []*tax.ScenarioSet{
			invoiceScenarios,
		},
		Identities: identityDefinitions, // identities.go
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
		Normalizer: Normalize,
		Categories: taxCategories,
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *org.Identity:
		return validateTaxNumber(obj)
	}
	return nil
}

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	case *org.Identity:
		normalizeTaxNumber(obj)
	}
}
