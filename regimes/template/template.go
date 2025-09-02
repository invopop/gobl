//go:build ignore
// +build ignore

// Package template provides a template for creating new regimes.
package template

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Name: i18n.String{
			i18n.EN: "Template",
			// Add official local name here.
			// i18n.XX: "Template",
		},
		Country:   l10n.TaxCountryCode("XX"),
		Currency:  currency.XXX,
		TaxScheme: tax.CategoryVAT,
		TimeZone:  "Europe/London",
		Tags: []*tax.TagSet{
			common.InvoiceTags(),
		},
		Identities: identityTypeDefinitions(), // org_identities.go
		Categories: taxCategories(),           // tax_categories.go
		Scenarios:  scenarios(),               // scenarios.go
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
			},
		},
		Normalizer: Normalize,
		Validator:  Validate,
	}
}

// Normalize will perform any regime specific calculations.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	case *org.Identity:
		normalizeOrgIdentity(obj)
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc any) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *org.Identity:
		return validateOrgIdentity(obj)
	case *org.Party:
		return validateOrgParty(obj)
	}
	return nil
}
