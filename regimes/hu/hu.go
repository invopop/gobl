// Package hu provides the Hungarian tax regime.
package hu

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New instantiates a new Hungarian regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:  "HU",
		Currency: currency.HUF,
		Name: i18n.String{
			i18n.EN: "Hungary",
			i18n.HU: "Magyarország",
		},
		TimeZone:   "Europe/Budapest",
		Categories: taxCategories,
		Tags: []*tax.TagSet{
			common.InvoiceTags().Merge(invoiceTags),
		},
		Validator:    Validate,
		Normalizer:   Normalize,
		Scenarios:    scenarios,
		IdentityKeys: identityKeyDefinitions,
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *org.Identity:
		return validateIdentity(obj)
	}
	return nil
}

// Normalize will perform any regime specific normalizations.
func Normalize(doc interface{}) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
