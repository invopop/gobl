// Package nz provides the New Zealand tax regime.
package nz

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the ISO 3166-1 alpha-2 code for New Zealand.
const CountryCode = "NZ"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("nz", rules.GOBL.Add(CountryCode),
		taxIdentityRules(),
		orgIdentityRules(),
	)
}

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.NZD,
		TaxScheme: tax.CategoryGST,
		Name: i18n.String{
			i18n.EN: "New Zealand",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Inland Revenue New Zealand",
				},
				URL: "https://www.ird.govt.nz",
			},
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				New Zealand's tax system is administered by Inland Revenue
				(Te Tari Taake). Goods and Services Tax (GST) applies at a standard
				rate of 15% to most goods and services. A reduced rate applies to
				long-term accommodation of 28 or more consecutive days, equivalent
				to 60% of the full standard rate.

				Zero-rated supplies include exported goods, international transport
				services, and certain land transfers between GST-registered parties.
				Exempt supplies include residential dwelling rent and most financial
				services.

				Businesses with annual turnover of NZD 60,000 or more must register
				for GST. The GST registration number is the same as the IRD number
				(Inland Revenue Department number), an 8 or 9 digit identifier
				validated using a two-pass modulo-11 checksum algorithm.

				Businesses are also identified by a New Zealand Business Number
				(NZBN), a 13-digit GS1 Global Location Number used for Peppol
				e-invoicing. E-invoicing via the Peppol network using the PINT A-NZ
				specification is mandatory for government agencies and large suppliers
				to government.

				Supply Correction Information (credit and debit notes) may be issued
				to correct a previously issued invoice.
			`),
		},
		TimeZone:   "Pacific/Auckland",
		Normalizer: Normalize,
		Identities: identityDefs,
		Scenarios: []*tax.ScenarioSet{
			invoiceScenarios,
		},
		Categories: taxCategories(),
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
			},
		},
	}
}

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	case *org.Identity:
		normalizeOrgIdentity(obj)
	}
}
