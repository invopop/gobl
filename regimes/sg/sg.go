// Package sg provides the tax region definition for Singapore.
package sg

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "SG",
		Currency:  currency.SGD,
		TaxScheme: tax.CategoryGST,
		Name: i18n.String{
			i18n.EN: "Singapore",
		},
		TimeZone: "Asia/Singapore",
		Tags: []*tax.TagSet{
			invoiceTags(),
		},
		Description: i18n.String{
			i18n.EN: `Singapore offers a simple GST model with a standard rate along with a few exceptions. GST is handled by the Inland Revenue Authority of Singapore ([IRAS](https://www.iras.gov.sg/taxes/goods-services-tax-(gst)))

For GST to be chargeable on a supply of goods and services, the following four conditions must be satisfied:

1. The supply must be made in Singapore
2. The supply is a taxable supply
3. The supply is made by a taxable person
4. The supply is made in the course of furtherance of any business carried on by the taxable person, i.e, GST is not chargeable on personal transactions

GST is chargeable on all imported goods (whether for domestic consumption, sale, or re-export), regardless of whether the importer is GST-registered or not. The importer is required to take up the appropriate import permit and pay GST upon importation of the goods into Singapore. Import GST is not chargeable under the following circumstances:

1. Importation of investment precious metals
2. Importation of goods that are specifically given GST reliefs5 under the GST
Act
3. Importation of goods into Zero-GST/Licensed warehouses administered by
Singapore Customs 
4. Importation of goods by GST-registered businesses that are under Major
Exporter Scheme or other approved schemes.`,
		},
		Scenarios: []*tax.ScenarioSet{
			invoiceScenarios(),
		},
		Identities: identityDefinitions, // identities.go
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				// Singpore only supports credit notes to correct an invoice:
				// https://www.iras.gov.sg/taxes/goods-services-tax-(gst)/basics-of-gst/invoicing-price-display-and-record-keeping/invoicing-customers
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
			},
		},
		Validator:  Validate,
		Normalizer: Normalize,
		Categories: taxCategories(),
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc interface{}) error {
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

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	case *org.Identity:
		normalizeIdentity(obj)
	}
}
