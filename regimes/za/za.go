// Package za provides the tax regime for South Africa.
package za

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

// CountryCode is the ISO 3166-1 alpha-2 code for South Africa.
const CountryCode = "ZA"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("za", rules.GOBL.Add(CountryCode),
		taxIdentityRules(),
		orgIdentityRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(func(id *tax.Identity) { tax.NormalizeIdentity(id) })),
	)
}

// New instantiates a new South Africa regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.ZAR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "South Africa",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				South Africa's tax system is administered by the South African
				Revenue Service (SARS) under the Value-Added Tax Act No. 89 of 1991.

				VAT is levied at a standard rate of 15%, in effect since 1 April 2018
				(previously 14% since 1993). A 2025 budget proposal to raise the rate
				further to 15.5% and then 16% was withdrawn before taking effect, and
				the rate has remained at 15%.

				Zero-rated supplies (0%, with input VAT still recoverable) include
				basic foodstuffs, fuel, exports, international transport, and
				services consumed outside South Africa. Exempt supplies (no VAT
				charged, no input VAT recovery) include financial services,
				residential rental, non-international passenger transport, and
				education. The VAT Act treats these very differently, and this
				regime preserves the distinction rather than collapsing both into a
				single "0%" concept.

				Vendors must register for VAT once taxable supplies exceed the
				compulsory threshold (ZAR 2.3 million over 12 months, since 1 April
				2026; ZAR 1 million previously) and may register voluntarily above a
				lower threshold (ZAR 120,000, since 1 April 2026; ZAR 50,000
				previously). Businesses below these thresholds may trade without a
				VAT number, so one is not always present.

				Businesses are identified for VAT purposes by a 10-digit VAT
				registration number that always starts with the digit 4 (e.g.
				4480152117). Unlike most countries covered by GOBL, SARS does not
				publish a check digit algorithm for this number, so only its format
				can be validated here. Authoritative verification requires a live
				lookup against SARS's VAT Vendor Search service, which is outside
				the scope of this library.

				Registered companies must additionally display their CIPC company
				registration number (format YYYY/NNNNNN/XX) on invoices and other
				business correspondence, per section 32(4) of the Companies Act 71
				of 2008.

				Both credit and debit notes are supported for invoice corrections,
				per section 21 of the VAT Act.

				As of mid-2026, e-invoicing is not yet mandatory in South Africa.
				The Tax Administration Laws Amendment Act (published April 2026)
				establishes a legal framework for a phased mandate expected to
				complete by 2028, but the technical format has not yet been
				finalized, so no e-invoicing addon is implemented here.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("SARS - Value-Added Tax"),
				URL:   "https://www.sars.gov.za/types-of-tax/value-added-tax/",
			},
			{
				Title: i18n.NewString("Value-Added Tax Act No. 89 of 1991"),
				URL:   "https://www.gov.za/sites/default/files/gcis_document/201505/act-89-1991s.pdf",
			},
			{
				Title: i18n.NewString("SARS - VAT rate increase from 1 April 2018"),
				URL:   "https://www.sars.gov.za/wp-content/uploads/Docs/VAT/NON-BDE-Rate-change-letter-final-.pdf",
			},
			{
				Title: i18n.NewString("Zero-rated and exempt supplies (Parliamentary Monitoring Group)"),
				URL:   "https://static.pmg.org.za/docs/Zero-rated%20and%20exempt%20supplies.pdf",
			},
			{
				Title: i18n.NewString("SARS - Tax Invoices"),
				URL:   "https://www.sars.gov.za/businesses-and-employers/government/tax-invoices/",
			},
			{
				Title: i18n.NewString("Companies Act 71 of 2008, Section 32"),
				URL:   "https://marxgore.co.za/wp-content/uploads/2020/01/Section-32-Use-of-company-name-and-registration-number.pdf",
			},
		},
		TimeZone: "Africa/Johannesburg",
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
			},
		},
		Categories: taxCategories,
		Identities: identityTypeDefinitions,
		Scenarios:  []*tax.ScenarioSet{bill.InvoiceScenarios()},
	}
}
