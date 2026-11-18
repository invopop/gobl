// Package ae provides the tax region definition for United Arab Emirates.
package ae

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

// CountryCode is the tax country code for United Arab Emirates.
const CountryCode = "AE"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("ae", rules.GOBL.Add(CountryCode),
		taxIdentityRules(),
		orgIdentityRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode),
			norm.For(func(id *tax.Identity) {
				tax.NormalizeIdentity(id)
			})),
	)
}

// New provides the tax region definition for AE.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.AED,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "United Arab Emirates",
			i18n.AR: "الإمارات العربية المتحدة",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				The UAE tax system is administered by the Federal Tax Authority (FTA).
				VAT was introduced on January 1, 2018, with standard, zero, and
				exempt rate categories.

				Businesses must register for VAT if taxable supplies and imports
				exceed AED 375,000 in a 12-month period, with voluntary registration
				available above AED 187,500. Registered businesses receive a Tax
				Registration Number (TRN) which must be included on all tax invoices.

				Simplified VAT invoices may be used when the recipient is not VAT
				registered, or when the transaction value does not exceed AED 10,000
				for VAT-registered recipients. Credit notes are supported for
				correcting invoices.

				Supplies of eligible goods under the profit margin scheme use the
				"margin-scheme" tag, with the "exempt" VAT key on their lines. VAT
				is due on the margin and reported in the VAT return, so the invoice
				shows no VAT amount, only a statement that the scheme applies.
				Using the "exempt" key is deliberate: GOBL classifies margin scheme
				lines the same way for Italy, Greece and Portugal.

				Electronic invoicing is mandatory for any person conducting
				business in the UAE, unless excluded by Ministerial Decision No. 243
				of 2025, following the phased timeline of Ministerial Decision
				No. 244 of 2025. Electronic invoices follow PINT AE, the UAE's
				Peppol specification, and identify the parties by their 10-digit
				Tax Identification Number (TIN).
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Federal Tax Authority - VAT"),
				URL:   "https://tax.gov.ae/en/taxes/Vat/vat.topics/registration.for.vat.aspx",
			},
			mofMandatoryFieldsSource,
		},
		TimeZone: "Asia/Dubai",
		Tags: []*tax.TagSet{
			invoiceTags,
		},
		Scenarios: []*tax.ScenarioSet{
			invoiceScenarios,
		},
		Identities: identityDefinitions,
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
			},
		},
		Categories: taxCategories,
	}
}
