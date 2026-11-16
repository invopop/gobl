// Package sk provides the Slovak tax regime.
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
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the tax country code for Slovakia.
const CountryCode = "SK"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("sk", rules.GOBL.Add(CountryCode),
		billInvoiceRules(),
		orgIdentityRules(),
		taxIdentityRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(func(tID *tax.Identity) { tax.NormalizeIdentity(tID) })),
	)
	norm.RegisterWithGuard(is.InContext(tax.RegimeIn(CountryCode)),
		norm.For(normalizeOrgIdentity),
	)
}

// New provides the tax region definition for Slovakia.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Slovakia",
			i18n.SK: "Slovensko",
		},

		// Short description
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Slovakia's tax system is administered by the Financial Administration of the
				Slovak Republic (Finančná správa Slovenskej republiky). As an EU member state,
				Slovakia follows the EU VAT Directive.

				VAT (daň z pridanej hodnoty, DPH) applies at a standard rate, with reduced rates
				for certain goods and services such as foodstuffs, medicines, medical devices,
				books, and accommodation. Exports and intra-EU supplies to VAT-registered buyers
				are exempt with the right to deduct.

				Every taxpayer has a DIČ (daňové identifikačné číslo), a ten digit number whose
				first digit gives the type of person. Businesses registered for VAT also have an
				IČ DPH, formed by the prefix SK followed by the same ten digits (e.g.
				SK2020273893). Both go in the tax ID code, which the supplier of an invoice must
				have. The IČO, the eight digit number the Statistical Office assigns to legal
				persons and entrepreneurs, is the identity type "ICO".
			`),
		},

		// Trustworthy sources about tax systems in Slovakia (Finančná správa)
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
			{
				Title: i18n.NewString("Peppol Code Lists - Participant identifier schemes 0245 (SK:DIC) and 0158 (SK:ICO)"),
				URL:   "https://docs.peppol.eu/edelivery/codelists/v9.7/Peppol%20Code%20Lists%20-%20Participant%20identifier%20schemes%20v9.7.html",
				At:    cal.NewDateTime(2026, 11, 16, 0, 0, 0),
			},
		},
		TimeZone:   "Europe/Bratislava",
		Identities: identityTypeDefinitions,
		Categories: taxCategories,
		Scenarios:  []*tax.ScenarioSet{bill.InvoiceScenarios()},
	}
}
