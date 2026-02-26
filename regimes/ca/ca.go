// Package ca provides models for dealing with Canada.
package ca

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// Tax categories specific for Canada.
const (
	TaxCategoryHST cbc.Code = "HST"
	TaxCategoryPST cbc.Code = "PST"
)

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "CA",
		Currency:  currency.CAD,
		TaxScheme: tax.CategoryGST,
		Name: i18n.String{
			i18n.EN: "Canada",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Canada's tax system is administered by the Canada Revenue Agency
				(CRA). The country uses a multi-layered sales tax system consisting
				of the federal Goods and Services Tax (GST) and various provincial
				taxes.

				The Harmonized Sales Tax (HST) combines GST and provincial sales
				tax in participating provinces. Non-participating provinces levy a
				separate Provincial Sales Tax (PST) at varying rates. Zero-rated
				supplies include basic groceries, agricultural products, and
				exports. Exempt supplies include certain financial services,
				educational services, and healthcare services.

				Businesses with annual taxable revenues exceeding CAD 30,000 must
				register for GST/HST. Tax identification is through the Business
				Number (BN) assigned by the CRA. Canada supports both credit notes
				and debit notes for invoice corrections.
			`),
		},
		TimeZone:   "America/Toronto", // Toronto
		Validator:  Validate,
		Normalizer: Normalize,
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
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	}
	return nil
}

// Normalize will attempt to clean the object passed to it.
func Normalize(doc interface{}) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
