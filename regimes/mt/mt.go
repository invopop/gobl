// Package mt provides the tax regime definition for Malta.
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
	rules.Register("mt", rules.GOBL.Add(CountryCode),
		taxIdentityRules(),
		billInvoiceRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(func(id *tax.Identity) { tax.NormalizeIdentity(id) })),
	)
}

// New instantiates a new Maltese tax regime.
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
				Malta's VAT is governed by the Value Added Tax Act (Cap. 406) and administered
				by the Malta Tax and Customs Administration (MTCA). As an EU member state Malta
				follows the EU VAT Directive (2006/112/EC).

				The standard rate is 18%, with reduced rates of 12%, 7% and 5%. Supplies that are
				"exempt with credit" (Fifth Schedule Part One - exports, intra-community supplies,
				food, pharmaceuticals) are zero-rated with input-VAT recovery; supplies "exempt
				without credit" (Fifth Schedule Part Two - insurance, finance, letting of
				immovable property, health, education) carry no input-VAT recovery.

				Businesses are identified by a VAT identification number: the prefix "MT" followed
				by eight digits (Article 10 and Article 12 registrations). Full tax invoices must
				show both the supplier's and the customer's VAT number; simplified invoices are
				allowed for amounts up to EUR 100. Corrections use credit notes and debit notes.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Value Added Tax Act (Cap. 406)"),
				URL:   "https://legislation.mt/eli/cap/406/eng",
			},
			{
				Title: i18n.NewString("MTCA - VAT Rates"),
				URL:   "https://mtca.gov.mt/business-tax/vat1/vat-compliance/vat-rates/vat-rates",
			},
			{
				Title: i18n.NewString("MTCA - VAT Rates and Exemptions FAQs"),
				URL:   "https://mtca.gov.mt/docs/default-source/documents/business-tax/vat/faqs/vat-rates-and-exemptions---faqs.pdf",
			},
		},
		TimeZone:   "Europe/Malta",
		Categories: taxCategories,
		Scenarios: []*tax.ScenarioSet{
			bill.InvoiceScenarios(),
		},
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				// Cap. 406 Tenth Schedule item 1(h) recognises credit notes (a decrease in
				// consideration) and debit notes (an increase) as the documents that correct
				// an issued VAT invoice.
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
			},
		},
	}
}
