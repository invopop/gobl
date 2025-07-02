// Package ar provides the AR tax regime.
package ar

import (
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New returns the AR regime definition.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "AR",
		Currency:  currency.ARS,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Argentina",
			i18n.ES: "Argentina",
		},
		TimeZone: "America/Argentina/Buenos_Aires",
		Tags: []*tax.TagSet{
			common.InvoiceTags(),
		},
		Validator:  Validate,
		Normalizer: Normalize,
		Categories: categories, // categories.go
	}
}

// Validate performs validation on tax identity.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return ValidateTaxIdentity(obj)
	}
	return nil
}

// Normalize normalizes tax identity before validation.
func Normalize(doc interface{}) {
	switch obj := doc.(type) {
	case *tax.Identity:
		NormalizeTaxIdentity(obj)
	}
}
