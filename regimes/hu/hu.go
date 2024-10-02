// Package hu provides the Hungarian tax regime.
package hu

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegime(New())
}

// New instantiates a new Hungarian regime.
func New() *tax.Regime {
	return &tax.Regime{
		Country:  "HU",
		Currency: currency.HUF,
		Name: i18n.String{
			i18n.EN: "Hungary",
			i18n.HU: "Magyarország",
		},
		TimeZone:   "Europe/Budapest",
		Extensions: extensionKeys,
		Categories: taxCategories,
		Tags:       invoiceTags,
		Validator:  Validate,
		Calculator: Calculate,
		Scenarios:  scenarios,
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

// Calculate will perform any regime specific calculations.
func Calculate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return normalizeTaxIdentity(obj)
	}
	return nil
}
