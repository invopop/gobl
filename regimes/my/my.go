// Package it provides the Malaysian tax regime.
package my

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New instantiates a new Malaysian regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "MY",
		Currency:  currency.MYR,
		TaxScheme: tax.CategoryST, // tax/constants.go. Malaysia uses SST — a form of ST in GOBL terms.
		Name: i18n.String{
			i18n.EN: "Malaysia",
			i18n.MS: "Malaysia", // Write Malay in Malay
		},
		TimeZone: "Asia/Kuala_Lumpur",
		Tags: []*tax.TagSet{
			common.InvoiceTags(), // CHECK which tags I used use
		},
		Validator:  Validate,
		Normalizer: Normalize,
		Categories: taxCategories, // tax_categories.go
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

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	}
}
