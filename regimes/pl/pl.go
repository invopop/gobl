// Package pl provides the Polish tax regime.
package pl

import (
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New instantiates a new Polish regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "PL",
		Currency:  currency.PLN,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Poland",
			i18n.PL: "Polska",
		},
		TimeZone:   "Europe/Warsaw",
		Validator:  Validate,
		Normalizer: Normalize,
		Categories: taxCategories, // tax_categories.go
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc any) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Normalize will perform any regime specific normalizations.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
