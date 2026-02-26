// Package sa provides the tax regime definition for Saudi Arabia.
package sa

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the tax regime definition for SA.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "SA",
		Currency:  currency.SAR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Kingdom of Saudi Arabia",
			i18n.AR: "المملكة العربية السعودية",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Saudi Arabia's tax system is administered by [ZATCA](https://zatca.gov.sa) (Zakat, Tax
				and Customs Authority), which oversees the collection of VAT (Value Added Tax) introduced
				under the GCC VAT Framework Agreement.

				Tax identification uses 15-digit TIN (Tax Identification Number) codes issued by ZATCA,
				validated with a Luhn check digit. Businesses are additionally identified through several
				types: CRN (Commercial Registration Number) from the Ministry of Commerce, MOM (MOMRAH
				License) from the Ministry of Municipal and Rural Affairs, MLS (MHRSD License) from the
				Ministry of Human Resources, 700 (Unified Number), and SAG (MISA License) from the
				Ministry of Investment. Buyers may be identified by TIN, NAT (National ID), IQA (Iqama
				residency permit), PAS (Passport), or GCC (GCC ID). All seller identity types are also
				valid for buyers per ZATCA business rule BR-KSA-14.

				The standard VAT rate is 15%, effective since July 2020 under Royal Order No. A/638,
				increased from the original 5% rate introduced in January 2018. Zero-rated supplies cover
				exports, international transport, qualified means of transport, medicines and medical
				equipment, qualifying metals, and private education and healthcare for citizens. Exempt
				supplies include financial services and life insurance (Article 29) and real estate
				transactions (Article 30).

				Invoice validation enforces ZATCA business rules: supplier TaxID is required on all
				invoices (BR-KSA-39), supplier name is required (BR-06), customer name and identification
				are required on standard B2B invoices (BR-KSA-42, BR-KSA-81), while simplified B2C
				invoices skip customer requirements.

			`),
		},
		TimeZone:   "Asia/Riyadh",
		Validator:  Validate,
		Normalizer: Normalize,
		Categories: taxCategories,
		Identities: identityDefinitions,
		Scenarios: []*tax.ScenarioSet{
			bill.InvoiceScenarios(),
			invoiceScenarios,
		},
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

// Validate checks the document type and determines if it can be validated.
func Validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *org.Identity:
		return validateIdentity(obj)
	}
	return nil
}

// Normalize attempts to clean up the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	case *org.Identity:
		normalizeIdentity(obj)
	}
}
