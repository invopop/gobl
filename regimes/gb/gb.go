// Package gb provides the United Kingdom tax regime.
package gb

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// Identification code types unique to the United Kingdom.
const (
	IdentityTypeCRN cbc.Code = "CRN" // Company Registration Number
)

var (
	altCountryCodes = []l10n.Code{
		l10n.XI, // Northern Ireland
		l10n.XU, // UK except Northern Ireland (Brexit)
	}
)

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:         "GB",
		AltCountryCodes: altCountryCodes,
		Currency:        currency.GBP,
		TaxScheme:       tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "United Kingdom",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				The United Kingdom's tax system is administered by His Majesty's Revenue and
				Customs (HMRC). Following Brexit, the UK operates its own VAT system
				independently of the EU VAT Directive.

				VAT rates include a 20% standard rate for most goods and services, a 5%
				reduced rate for domestic fuel, children's car seats, and certain other goods,
				and a 0% zero rate for food, children's clothing, books, and newspapers.
				Some supplies are exempt from VAT, including financial services, education,
				and healthcare.

				Businesses with taxable turnover exceeding GBP 90,000 must register for VAT.
				Companies are identified by their VAT Registration Number (VRN) in the format
				GB followed by 9 digits, and optionally by their Company Registration Number
				(CRN) from Companies House.

				Northern Ireland (country code XI) has special arrangements for goods under
				the Windsor Framework, remaining aligned with EU VAT rules for goods while
				following UK rules for services. Credit notes are supported for invoice
				corrections.
			`),
		},
		TimeZone:   "Europe/London",
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

// Validate checks the document type and determines if it can be validated. Note that in
// the GB tax regime we don't need to validate the presence of the supplier's tax ID.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Normalize will attempt to clean the object passed to it.
func Normalize(doc interface{}) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj, altCountryCodes...)
	}
}
