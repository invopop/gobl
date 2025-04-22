// Package no provides the tax region definition for Norway.
package no

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the tax region definition for NO.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "NO",
		Currency:  currency.NOK,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
            i18n.EN: "Norway",
            i18n.NO: "Norge",
		},
        TimeZone: "Europe/Oslo",
		Tags: []*tax.TagSet{
			common.InvoiceTags().Merge(invoiceTags),
		},
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
			},
		},
		Scenarios:  scenarios,
		Extensions: extensionKeys,
		Validator:  Validate,
		Normalizer: Normalize,
		Categories: taxCategories,
	}
}

// Validate function assesses the document type to determine if validation is required.
// Note that, under the NO tax regime, validation of the supplier's tax ID is not necessary if it does not meet the specified threshold (refer to the README section for more details).
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
