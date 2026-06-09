// Package gb provides the United Kingdom tax regime.
package gb

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the tax country code for the United Kingdom.
const CountryCode = "GB"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register(
		"gb",
		rules.GOBL.Add(CountryCode),
		taxIdentityRules(),
	)
	norm.Register(
		// XI (Northern Ireland) and XU also resolve to this regime (see
		// AltCountryCodes), so normalize identities under any of them.
		norm.When(tax.IdentityIn(CountryCode, "XI", "XU"), norm.For(func(id *tax.Identity) { tax.NormalizeIdentity(id, altCountryCodes...) })),
	)
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

				VAT applies at standard, reduced, and zero rates. Zero-rated supplies include
				food, children's clothing, books, and newspapers. Some supplies are exempt
				from VAT, including financial services, education, and healthcare.
				Companies are identified by their VAT Registration Number (VRN) in the format
				GB followed by 9 digits, and optionally by their Company Registration Number
				(CRN) from Companies House.

				Northern Ireland (country code XI) has special arrangements for goods under
				the Windsor Framework, remaining aligned with EU VAT rules for goods while
				following UK rules for services. Credit notes are supported for invoice
				corrections.
			`),
		},
		TimeZone: "Europe/London",
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
