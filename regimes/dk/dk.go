// Package dk provides a regime definition for Denmark.
package dk

import (
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New instantiates a new Danish regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   l10n.DK.Tax(),
		Currency:  currency.DKK,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Denmark",
			i18n.DA: "Danmark",
		},
		TimeZone:   "Europe/Copenhagen",
		Identities: identityKeyDefinitions,
		Categories: taxCategories,
		Validator:  Validate,
		Normalizer: Normalize,
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc any) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *org.Identity:
		return validateIdentity(obj)
	}
	return nil
}

// Normalize will perform any regime specific normalization.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
