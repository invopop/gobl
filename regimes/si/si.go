// Package si provides the tax region definition for Slovenia.
package si

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the tax country code for Slovenia.
const CountryCode = "SI"

// init registers the Slovenian regime, its validation rules and its
// normalization with GOBL so that they become available as soon as this
// package is imported.
func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("si", rules.GOBL.Add(CountryCode),
		billInvoiceRules(),
		orgIdentityRules(),
		taxIdentityRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(func(id *tax.Identity) { tax.NormalizeIdentity(id) })),
	)
}

// New provides the tax region definition for Slovenia.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Slovenia",
			i18n.SL: "Slovenija",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Slovenia's tax system is administered by the Financial Administration
				of the Republic of Slovenia (FURS). As an EU member state, Slovenia
				follows the EU VAT Directive with locally adapted rates.

				VAT (davek na dodano vrednost, DDV) applies at a standard rate with a
				reduced rate for specific goods and services such as food, medicines,
				passenger transport and accommodation, and a special reduced rate for
				books, newspapers and periodicals in force since 1 January 2020. The
				standard rate rose from 20% to 22% on 1 July 2013.

				Taxpayers are identified by an eight-digit tax number (davčna
				številka) whose last digit is a modulo-11 check digit. The VAT
				identification number is the same number prefixed with the SI
				country code (e.g. SI12345678).

				Companies are additionally identified in the Slovenian Business
				Register (Poslovni register Slovenije) by a ten-digit
				registration number (matična številka), also protected by a
				modulo-11 check digit.

				Slovenia supports credit notes for invoice corrections.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("FURS - Entry into the tax register and tax number"),
				URL:   "https://www.fu.gov.si/en/taxes_and_other_duties/work_with_us/entry_into_the_tax_register_and_tax_number",
			},
			{
				Title: i18n.NewString("SPOT - Invoicing"),
				URL:   "https://spot.gov.si/en/info/accountancy/invoicing",
			},
		},
		TimeZone:   "Europe/Ljubljana",
		Identities: identityDefinitions,
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
