// Package fr provides the tax region definition for France.
package fr

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:  "FR",
		Currency: currency.EUR,
		Name: i18n.String{
			i18n.EN: "France",
			i18n.FR: "La France",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				The French tax regime covers the basics.
			`),
		},
		TimeZone: "Europe/Paris",
		Tags: []*tax.TagSet{
			common.InvoiceTags(),
		},
		Scenarios: []*tax.ScenarioSet{
			invoiceScenarios,
		},
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				// France supports both corrective methods
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote, // Code 381
					bill.InvoiceTypeCorrective, // Code 384
				},
			},
		},
		Validator:     Validate,
		Normalizer:    Normalize,
		Categories:    taxCategories,
		IdentityTypes: identityTypeDefinitions, // identities.go
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

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	case *org.Identity:
		normalizeIdentity(obj)
	}
}
