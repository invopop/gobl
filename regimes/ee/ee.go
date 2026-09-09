// Package ee provides tax regime support for Estonia.
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

// CountryCode is the ISO 3166-1 alpha-2 code for Estonia.
const CountryCode = "EE"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("ee", rules.GOBL.Add(CountryCode),
		taxIdentityRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(func(tID *tax.Identity) { tax.NormalizeIdentity(tID) })),
	)
}

// New instantiates a new Estonia regime
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
				Estonia's tax system is administered by the Estonian Tax and Customs Board (Maksu- ja Tolliamet or EMTA).
				As an EU member state, Estonia follows the EU VAT Directive.
				A general, intermediate, reduced, and zero rates apply, with certain exemptions for specific goods and services.
				Estonian businesses are required to register for the VAT system if their revenue exceeds the 40,000 EUR threshold.
				Businesses are identified by a unique VAT number (käibemaksukohustuslase number) with prefix EE.
				E-invoicing is mandated for accounting entities regiistered as e-invoice recipients in the commerical register.
				Given that all public sector entities are registered as e-invoice recipients, this effectively mandates e-invoicing for all B2G transactions.
				The European standard for e-invoicing (EN 16931) is used. There is no mandate for B2C e-invoicing.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Estonian Tax and Customs Board - Value Added Tax"),
				URL:   "https://www.emta.ee/en/business-client/taxes-and-payment/value-added-tax",
			},
			{
				Title: i18n.NewString("Riigi Teataja - Käibemaksuseadus"),
				URL:   "https://www.riigiteataja.ee/en/akt/511072022006",
			},
			{
				Title: i18n.NewString("Estonian Tax and Customs Board - VAT rates and supply exempt from tax"),
				URL:   "https://www.emta.ee/en/business-client/taxes-and-payment/value-added-tax/vat-rates-and-supply-exempt-tax/standard-vat-rate",
			},
			{
				Title: i18n.NewString("European Commission - eInvoicing in Estonia"),
				URL:   "https://ec.europa.eu/digital-building-blocks/sites/spaces/DIGITAL/pages/467108883/eInvoicing+in+Estonia",
			},
		},
		TimeZone:   "Europe/Tallinn",
		Categories: taxCategories,
		Scenarios:  []*tax.ScenarioSet{bill.InvoiceScenarios()},
	}
}
