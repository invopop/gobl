// Package ee provides the Estonian tax regime.
package ee

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

// CountryCode is the tax country code for Estonia.
const CountryCode = "EE"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("ee", rules.GOBL.Add(CountryCode), taxIdentityRules())
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(func(id *tax.Identity) { tax.NormalizeIdentity(id) })),
	)
}

// New provides the tax regime definition for Estonia.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Estonia",
			i18n.ET: "Eesti",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Estonia's tax system is administered by the Estonian Tax and Customs
				Board (Maksu- ja Tolliamet, EMTA). As an EU member state, Estonia follows
				the EU VAT Directive with standard and reduced rates.

				VAT (Käibemaks) is applied to most goods and services. The standard rate
				is 24% since 1 July 2025, raised from 22% (in force since 1 January 2024),
				which in turn replaced the long-standing 20% rate in force since 2009.

				Businesses register for VAT via a KMKR number (käibemaksukohustuslase
				registreerimisnumber) in the format EE followed by 9 digits. Credit notes
				are supported for invoice corrections.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Estonian Tax and Customs Board - VAT rates"),
				URL:   "https://www.emta.ee/en/business-client/taxes-and-payment/value-added-tax/vat-rates-and-supply-exempt-tax/standard-vat-rate",
			},
		},
		TimeZone: "Europe/Tallinn",
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
