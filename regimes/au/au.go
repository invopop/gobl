// Package au provides the Australian tax regime.
package au

import (
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the tax region definition for Australia.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "AU",
		Currency:  currency.AUD,
		TaxScheme: tax.CategoryGST,
		Name: i18n.String{
			i18n.EN: "Australia",
		},
		TimeZone:   "Australia/Sydney",
		Validator:  Validate,
		Normalizer: Normalize,
		Categories: taxCategories,
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

// Normalize will attempt to clean the object passed to it.
func Normalize(doc interface{}) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
