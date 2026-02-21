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
	return &tax.RegimeDef{
		Country:   "NZ",
		Currency:  currency.NZD,
		TaxScheme: tax.CategoryGST,
		Name: i18n.String{
			i18n.EN: "New Zealand",
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
