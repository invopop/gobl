// Package it provides the Italian tax regime.
package it

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New instantiates a new Italian regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:  "IT",
		Currency: currency.EUR,
		Name: i18n.String{
			i18n.EN: "Italy",
			i18n.IT: "Italia",
		},
		TimeZone:     "Europe/Rome",
		ChargeKeys:   chargeKeyDefinitions,   // charges.go
		IdentityKeys: identityKeyDefinitions, // identities.go
		Scenarios:    scenarios,              // scenarios.go
		Tags: []*tax.TagSet{
			common.InvoiceTags(),
		},
		Validator:  Validate,
		Normalizer: Normalize,
		Categories: categories, // categories.go
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
			},
		},
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *org.Identity:
		return validateIdentity(obj)
	}
	return nil
}

// Normalize will perform any regime specific calculations.
func Normalize(doc interface{}) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	case *org.Identity:
		normalizeIdentity(obj)
	case *org.Party:
		normalizeParty(obj)
	case *tax.Combo:
		normalizeTaxCombo(obj)
	}
}
