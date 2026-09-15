// Package sk provides tax regime support for Slovakia.
package sk

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the ISO 3166-1 alpha-2 code for Slovakia.
const CountryCode = "SK"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("sk", rules.GOBL.Add(CountryCode),
		taxIdentityRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(func(tID *tax.Identity) { tax.NormalizeIdentity(tID) })),
	)
}

// New instantiates a new Slovakia regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Slovakia",
			i18n.SK: "Slovensko",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Slovakia's tax system is administered by the Financial Administration of the
				Slovak Republic (Finančná správa Slovenskej republiky). As an EU member state,
				Slovakia follows the EU VAT Directive.

				VAT (daň z pridanej hodnoty, DPH) applies at a standard rate, with reduced rates
				for certain goods and services such as foodstuffs, medicines, medical devices,
				books, and accommodation. Exports and intra-EU supplies to VAT-registered buyers
				are zero-rated.

				Businesses registered for VAT are identified by their IČ DPH, formed by the
				prefix SK followed by ten digits (e.g. SK2020273893).
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Finančná správa - Daň z pridanej hodnoty (VAT)"),
				URL:   "https://www.financnasprava.sk/sk/podnikatelia/dane/dan-z-pridanej-hodnoty",
				At:    cal.NewDateTime(2026, 8, 6, 0, 0, 0),
			},
			{
				Title: i18n.NewString("Slov-Lex - Zákon č. 222/2004 Z. z. o dani z pridanej hodnoty"),
				URL:   "https://www.slov-lex.sk/ezbierky/pravne-predpisy/SK/ZZ/2004/222/",
				At:    cal.NewDateTime(2026, 8, 6, 0, 0, 0),
			},
			{
				Title: i18n.NewString("Financial Administration of the Slovak Republic (English)"),
				URL:   "https://www.financnasprava.sk/en/",
				At:    cal.NewDateTime(2026, 8, 6, 0, 0, 0),
			},
		},
		TimeZone:   "Europe/Bratislava",
		Categories: taxCategories,
		Scenarios:  []*tax.ScenarioSet{bill.InvoiceScenarios()},
	}
}
