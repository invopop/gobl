// Package lv provides the tax region definition for Latvia.
package lv

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

// CountryCode is the tax country code for Latvia.
const CountryCode = "LV"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("lv", rules.GOBL.Add(CountryCode), taxIdentityRules())
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(func(id *tax.Identity) { tax.NormalizeIdentity(id) })),
	)
}

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Latvia",
			i18n.LV: "Latvija",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Latvia's tax system is administered by the State Revenue Service
				(Valsts ieņēmumu dienests, VID). As an EU member state, Latvia
				follows the EU VAT Directive with standard, reduced, and
				super-reduced rates.

				VAT (Pievienotās vertības nodoklis, PVN) applies to most goods and services. Businesses are
				identified by their VAT registration number, which follows the format of the country code
				"LV" followed by 11 digits.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("State Revenue Service - VAT overview"),
				URL:   "https://www.vid.gov.lv/en/value-added-tax",
			},
		},
		TimeZone: "Europe/Riga",
		Scenarios: []*tax.ScenarioSet{
			bill.InvoiceScenarios(),
		},
		Categories: taxCategories,
	}
}
