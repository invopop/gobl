// Package mt provides tax regime support for Malta.
package mt

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

// CountryCode is the tax country code for Malta.
const CountryCode = "MT"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("mt", rules.GOBL.Add(CountryCode), taxIdentityRules())
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(func(id *tax.Identity) { tax.NormalizeIdentity(id) })),
	)
}

// New instantiates a new Malta regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Malta",
			i18n.MT: "Malta",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Malta's tax system is administered by the Malta Tax and Customs
				Administration (MTCA). As an EU member state Malta follows the EU VAT
				Directive, transposed by the Value Added Tax Act (Chapter 406).

				VAT applies at a standard rate of 18% under article 19(1), plus reduced
				rates of 12%, 7% and 5% set by the Eighth Schedule and bounded by article
				19(2) at no more than 18% and no less than 5%. Malta gives these rates no
				names, referring to them only by value, so they map onto GOBL's rate keys
				by magnitude rather than by translation.

				Exemptions come in two kinds and only one is GOBL's exempt. Part One of the
				Fifth Schedule is "exempt with credit": no VAT is charged but input tax
				stays recoverable under article 22(4), matching the zero, export and
				intra-community keys. Part Two is "exempt without credit" and carries no
				right of deduction, which is what exempt means.

				Businesses are identified by a VAT number of MT plus 8 digits, the last two
				a check over the first six. Article 13(3) gives that prefix to article 10
				and article 12 registrations only; an article 11 small-enterprise number has
				none, is not a VAT identification number, and never appears on a tax invoice.

				Corrections are not restricted to a document type: the Twelfth Schedule
				treats any document unambiguously amending an earlier invoice as an invoice,
				and the Eleventh Schedule recognizes credit and debit notes alike.

				Malta has no e-invoicing mandate. Contracting authorities must be able to
				receive EN 16931 invoices over PEPPOL for procurement above the EU
				thresholds, but suppliers need not issue electronically and there is no B2B
				mandate or clearance system.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Value Added Tax Act (Chapter 406)"),
				URL:   "https://legislation.mt/eli/cap/406/eng",
			},
			{
				Title: i18n.NewString("Legal Notice 231 of 2023 - Amendment of Eighth Schedule"),
				URL:   "https://legislation.mt/eli/ln/2023/231/eng",
			},
			{
				Title: i18n.NewString("MTCA - VAT Rates & Exemptions FAQs"),
				URL:   "https://mtca.gov.mt/docs/default-source/documents/business-tax/vat/faqs/vat-rates-exemptions---faqs.pdf",
			},
			{
				Title: i18n.NewString("MTCA - Tax Invoices and Fiscal Receipts FAQs"),
				URL:   "https://mtca.gov.mt/docs/default-source/documents/business-tax/vat/faqs/tax-invoice-and-fiscal-receipts---faqs.pdf",
			},
			{
				Title: i18n.NewString("MTCA - Registrations & De-Registrations FAQs"),
				URL:   "https://mtca.gov.mt/docs/default-source/documents/business-tax/vat/faqs/registrations-de-registrations-faqs.pdf",
			},
			{
				Title: i18n.NewString("MTCA - Guidelines on the VAT rules applicable to small enterprises"),
				URL:   "https://mtca.gov.mt/docs/default-source/documents/news/2025/guidelines-on-article-11-(2025)---final-clean.pdf",
			},
			{
				Title: i18n.NewString("European Commission - eInvoicing in Malta"),
				URL:   "https://ec.europa.eu/digital-building-blocks/sites/spaces/DIGITAL/pages/467108894/eInvoicing+in+Malta",
			},
		},
		TimeZone:   "Europe/Malta",
		Categories: taxCategories,
		Scenarios:  []*tax.ScenarioSet{bill.InvoiceScenarios()},

		// Corrections, Identities and a bill_invoices.go are all deliberately absent,
		// not overlooked. Maltese law does not restrict corrections to a document type
		// (Twelfth Sch. item 1(2), Eleventh Sch. item 1(1)(h)); a tax invoice carries no
		// non-tax business identifier (Twelfth Sch. item 3); and the only candidate
		// invoice rule, requiring the customer's VAT number, would reject the fiscal
		// receipts that article 51 provides for.
	}
}
