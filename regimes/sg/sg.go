// Package sg provides the tax region definition for Singapore.
package sg

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

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "SG",
		Currency:  currency.SGD,
		TaxScheme: tax.CategoryGST,
		Name: i18n.String{
			i18n.EN: "Singapore",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("IRAS - GST General Guide for Businesses"),
				URL:   "https://www.iras.gov.sg/media/docs/default-source/e-tax/etaxguide_gst_gst-general-guide-for-businesses(1).pdf",
			},
		},
		TimeZone: "Asia/Singapore",
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Singapore's tax system includes a Goods and Services Tax (GST) administered
				by the Inland Revenue Authority of Singapore (IRAS). Zero-rated supplies apply
				to international services and exports. Exempt supplies include financial
				services, sale and lease of residential properties, digital payment tokens,
				and investment precious metals.

				Businesses are identified by their Unique Entity Number (UEN). GST-registered
				suppliers must display their GST Registration Number on all tax invoices,
				which in most cases is the same as the UEN.

				Three invoicing methods are supported: tax invoices (standard, requiring full
				supplier and customer details), simplified tax invoices (for transactions
				up to 1000 SGD inclusive of GST), and receipts (for non-GST-registered
				customers). Credit notes are supported for correcting invoices; debit notes
				in Singapore are used for requesting payment on non-GST transactions, not
				for invoice corrections.
			`),
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
