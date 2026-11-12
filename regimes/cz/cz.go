// Package cz provides tax regime support for the Czech Republic.
package cz

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

// CountryCode is the ISO 3166-1 alpha-2 code for the Czech Republic.
const CountryCode = "CZ"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("cz", rules.GOBL.Add(CountryCode),
		taxIdentityRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(func(tID *tax.Identity) { tax.NormalizeIdentity(tID) })),
	)
}

// New instantiates a new Czech regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.CZK,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Czech Republic",
			i18n.CS: "Česká republika",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				The Czech tax system is administered by the Financial Administration
				(Finanční správa České republiky). As an EU member state, the Czech
				Republic follows the EU VAT Directive.

				VAT (daň z přidané hodnoty, DPH) applies at a standard rate and a reduced
				rate. Exports and intra-EU supplies to VAT-registered buyers are exempt
				with the right to deduct.

				Taxpayers are identified by their DIČ (daňové identifikační číslo): the
				prefix CZ followed by the 8-digit IČO of a legal person (e.g. CZ12345679),
				or by the 9 or 10-digit birth number or tax office identifier of a natural
				person.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Finanční správa - Daň z přidané hodnoty"),
				URL:   "https://financnisprava.gov.cz/cs/dane/dane/dan-z-pridane-hodnoty",
				At:    cal.NewDateTime(2026, 11, 6, 0, 0, 0),
			},
			{
				Title: i18n.NewString("e-Sbírka - Zákon č. 235/2004 Sb., o dani z přidané hodnoty"),
				URL:   "https://e-sbirka.gov.cz/sb/2004/235",
				At:    cal.NewDateTime(2026, 11, 6, 0, 0, 0),
			},
			{
				Title: i18n.NewString("Zákony pro lidi - Zákon č. 235/2004 Sb., o dani z přidané hodnoty"),
				URL:   "https://www.zakonyprolidi.cz/cs/2004-235",
				At:    cal.NewDateTime(2026, 11, 6, 0, 0, 0),
			},
			{
				Title: i18n.NewString("Zákon č. 280/2009 Sb., daňový řád, § 130 (Daňové identifikační číslo)"),
				URL:   "https://www.zakonyprolidi.cz/cs/2009-280#p130",
				At:    cal.NewDateTime(2026, 11, 6, 0, 0, 0),
			},
		},
		TimeZone:   "Europe/Prague",
		Categories: taxCategories,
		Scenarios:  []*tax.ScenarioSet{bill.InvoiceScenarios()},
	}
}
