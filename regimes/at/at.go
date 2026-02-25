// Package at provides the Austrian tax regime.
package at

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

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "AT",
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Austria",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Austria's tax system is administered by the Federal Ministry of
				Finance (Bundesministerium für Finanzen). As an EU member state,
				Austria follows the EU VAT Directive with standard, reduced, and
				intermediate rates.

				VAT (Umsatzsteuer, USt) applies to most goods and services.
				Businesses are identified by their UID-Nummer (VAT identification
				number) in the format ATU followed by 8 digits. Austria supports
				credit notes for invoice corrections.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("E-Rechnung - Austrian E-Invoicing"),
				URL:   "https://www.erechnung.gv.at/erb",
			},
		},
		TimeZone:   "Europe/Vienna",
		Validator:  Validate,
		Normalizer: Normalize,
		Scenarios: []*tax.ScenarioSet{
			bill.InvoiceScenarios(),
		},
		Categories: taxCategories,
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
			},
		},
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc any) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
