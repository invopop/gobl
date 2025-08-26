// Package gr provides the tax region definition for Greece.
package gr

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// Official IAPR codes to include in stamps.
const (
	StampIAPRQR       cbc.Key = "iapr-qr"
	StampIAPRMark     cbc.Key = "iapr-mark"
	StampIAPRHash     cbc.Key = "iapr-hash"
	StampIAPRUID      cbc.Key = "iapr-uid"
	StampIAPRProvider cbc.Key = "iapr-provider"
)

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country: "EL",
		AltCountryCodes: []l10n.Code{
			"GR", // regular ISO code
		},
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Greece",
			i18n.EL: "Ελλάδα",
		},
		TimeZone:               "Europe/Athens",
		CalculatorRoundingRule: tax.RoundingRuleCurrency,
		Scenarios:              scenarios,
		Corrections:            corrections,
		Validator:              Validate,
		Normalizer:             Normalize,
		Categories:             taxCategories,
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

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	}
}
