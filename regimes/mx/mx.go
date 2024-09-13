// Package mx provides the Mexican tax regime.
package mx

import (
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/regimes/mx/sat"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegime(New())

	// MX GOBL Schema Complements for CFDI

}

// New provides the tax region definition
func New() *tax.Regime {
	return &tax.Regime{
		Country:  "MX",
		Currency: currency.MXN,
		Name: i18n.String{
			i18n.EN: "Mexico",
			i18n.ES: "México",
		},
		TimeZone:    "America/Mexico_City",
		Validator:   Validate,
		Normalizer:  Normalize,
		Tags:        common.InvoiceTags(),
		Categories:  sat.TaxCategories(),
		Corrections: sat.CorrectionDefinitions(),
	}
}

// Validate validates a document against the tax regime.
func Validate(doc any) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return sat.ValidateTaxIdentity(obj)
	}
	return nil
}

// Normalize performs regime specific calculations.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
