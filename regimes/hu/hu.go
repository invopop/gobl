// Package hu provides a regime definition for Hungary.
package hu

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the ISO 3166-1 alpha-2 code for Hungary.
const CountryCode = "HU"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("hu", rules.GOBL.Add(CountryCode),
		billInvoiceRules(),
		taxIdentityRules(),
	)
}

// New instantiates a new Hungarian regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.HUF,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Hungary",
			i18n.HU: "Magyarország",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Hungary's tax system is administered by the National Tax and Customs
				Administration (Nemzeti Adó- és Vámhivatal, NAV). As an EU member state,
				Hungary follows the EU VAT Directive.

				VAT (ÁFA — Általános forgalmi adó) applies at a standard rate of 27%
				(the highest in the EU), an intermediate reduced rate of 18%, and a lower
				reduced rate of 5%. Certain supplies are zero-rated (exports, intra-community)
				or exempt (healthcare, education, financial services).

				Businesses are identified by their adószám (tax number), an 11-digit
				identifier in the format TTTTTTTT-V-RR, where the first 8 digits include
				a modulo-10 check digit, the 9th digit indicates VAT status, and the
				last 2 digits identify the regional tax office. The EU VAT number uses
				the HU prefix followed by the 8-digit base number.

				Hungary operates a mandatory real-time invoice reporting (RTIR) system
				through NAV Online Számla, requiring all invoices to be reported
				electronically since April 2021.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("NAV - Nemzeti Adó- és Vámhivatal"),
				URL:   "https://nav.gov.hu",
			},
			{
				Title: i18n.NewString("NAV Online Számla - Technical Documentation"),
				URL:   "https://github.com/nav-gov-hu/Online-Invoice",
			},
			{
				Title: i18n.NewString("OECD - Hungary Tax Identification Number"),
				URL:   "https://www.oecd.org/content/dam/oecd/en/topics/policy-issue-focus/aeoi/hungary-tin.pdf",
			},
		},
		TimeZone:   "Europe/Budapest",
		Categories: taxCategories,
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
			},
		},
		Normalizer: Normalize,
	}
}

// Normalize will perform any regime specific normalization.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	}
}
