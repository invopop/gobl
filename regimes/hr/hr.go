// Package hr provides a regime definition for Croatia.
package hr

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the ISO 3166-1 alpha-2 code for Croatia.
const CountryCode = "HR"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("hr", rules.GOBL.Add(CountryCode),
		orgIdentityRules(),
		taxIdentityRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(func(tID *tax.Identity) { tax.NormalizeIdentity(tID) })),
	)
}

// New instantiates a new Croatian regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Croatia",
			i18n.HR: "Hrvatska",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Croatia's tax system is administered by the Tax Administration (Porezna
				uprava). Croatia joined the European Union on 1 July 2013 and adopted the
				euro on 1 January 2023, replacing the Croatian kuna (HRK).

				VAT (porez na dodanu vrijednost, PDV) is regulated by the VAT Act (Zakon o
				porezu na dodanu vrijednost). It applies at a general rate, with reduced
				rates for the goods and services listed in Article 38 of the Act.

				Every legal and natural person in Croatia is assigned a personal
				identification number (osobni identifikacijski broj, OIB): 10 digits
				followed by a check digit calculated with ISO 7064 MOD 11,10. The OIB is
				the tax identification number, and the Croatian VAT number is the OIB with
				the HR prefix (e.g. HR12345678903).
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Porezna uprava - Porez na dodanu vrijednost (VAT)"),
				URL:   "https://porezna-uprava.gov.hr/en/vat/7362",
				At:    cal.NewDateTime(2026, 10, 19, 0, 0, 0),
			},
			{
				Title: i18n.NewString("Narodne novine 60/2008 - Zakon o osobnom identifikacijskom broju"),
				URL:   "https://narodne-novine.nn.hr/clanci/sluzbeni/2008_05_60_2033.html",
				At:    cal.NewDateTime(2026, 10, 19, 0, 0, 0),
			},
			{
				Title: i18n.NewString("Narodne novine 1/2009 - Pravilnik o osobnom identifikacijskom broju"),
				URL:   "https://narodne-novine.nn.hr/clanci/sluzbeni/2009_01_1_6.html",
				At:    cal.NewDateTime(2026, 10, 19, 0, 0, 0),
			},
		},
		TimeZone:   "Europe/Zagreb",
		Tags:       []*tax.TagSet{invoiceTags},
		Identities: identityTypeDefinitions,
		Categories: taxCategories,
		Scenarios:  []*tax.ScenarioSet{invoiceScenarios},
	}
}
