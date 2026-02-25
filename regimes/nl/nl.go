// Package nl provides the Dutch region definition
package nl

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the Dutch region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "NL",
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "The Netherlands",
			i18n.NL: "Nederland",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				The Netherlands' tax system is administered by the Belastingdienst (Tax and
				Customs Administration). As an EU member state, the Netherlands follows the
				EU VAT Directive with locally adapted rates.

				BTW (Belasting over de Toegevoegde Waarde) applies at standard and reduced
				rates. The reduced rate covers food, water, pharmaceuticals, books, passenger
				transport, hotel accommodation, and cultural and sporting events.

				Businesses are identified by their BTW-nummer (VAT number) in the format NL
				followed by 9 digits, the letter B, and 2 check digits (e.g.
				NL123456789B01). The KVK (Kamer van Koophandel) number is the commercial
				register number.

				The Netherlands supports credit notes for invoice corrections. E-invoicing
				via PEPPOL is commonly used, and is mandatory for B2G transactions with the
				central government.
			`),
		},
		TimeZone:   "Europe/Amsterdam",
		Validator:  Validate,
		Normalizer: Normalize,
		Scenarios: []*tax.ScenarioSet{
			bill.InvoiceScenarios(),
		},
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
			},
		},
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

// Normalize performs region specific calculations on the document.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
