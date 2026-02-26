// Package nz provides the tax region definition for New Zealand.
package nz

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the tax region definition.
func New() *tax.RegimeDef {
	// Note: The current implementation does not include special validation for second-hand goods or livestock transactions, which have additional requirements under NZ law:
	//   - Supplier address required
	//   - Date of supply required
	//   - Quantity/volume information required
	return &tax.RegimeDef{
		Country:   "NZ",
		Currency:  currency.NZD,
		TaxScheme: tax.CategoryGST,
		Name: i18n.String{
			i18n.EN: "New Zealand",
		},
		Description: i18n.String{
			i18n.EN: `New Zealand's Goods and Services Tax (GST) uses a standard rate of 15% for most taxable supplies.
Zero-rated (0%) applies to specific supplies such as exports.
Exempt supplies are not charged GST and do not allow input tax credits.

IRD numbers are 8 or 9 digits and are validated for format only.

Invoice validation enforces threshold-based requirements:
- invoices > $200 require supplier GST number
- invoices > $1,000 require customer name and at least one identifier (address, phone, email, tax ID, or identities)

Sources:
- https://www.ird.govt.nz/gst/charging-gst
- https://www.ird.govt.nz/gst/tax-invoices-for-gst/how-tax-invoices-for-gst-work
- https://www.ird.govt.nz/en/gst/charging-gst/zero-rated-supplies
- https://www.ird.govt.nz/en/gst/charging-gst/exempt-supplies
- https://www.ird.govt.nz/myir-help/logging-in/ird-numbers`,
		},
		TimeZone: "Pacific/Auckland",
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
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
	}
	return nil
}

// Normalize attempts to clean up the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
