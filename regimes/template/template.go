//go:build ignore

// Package template provides a template for creating new regimes.
package template

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New instantiates a new regime definition.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   l10n.TaxCountryCode("XX"),
		Currency:  currency.EUR, 		// Replace with the country's currency
		TaxScheme: tax.CategoryVAT, // or tax.CategoryGST, tax.CategoryST
		Name: i18n.String{
			i18n.EN: "English Name",
			// i18n.XX: "Local name",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Provide a concise overview of the country's tax system covering: the tax
				authority, main tax scheme and rates, business identification numbers,
				e-invoicing requirements, and supported correction methods (credit notes,
				debit notes, corrective invoices). Should be used instead of a README.md.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Tax Authority - Main Reference"),
				URL:   "https://example.com",
			},
		},
		TimeZone:   "Europe/London", // Replace with the country's base timezone
		Categories: taxCategories, // tax_categories.go

		// Scenarios auto-inject notes and codes into documents based on
		// tax tags and document type. Most regimes only need the shared
		// baseline below. Add a country-specific scenarios.go when the
		// regime requires its own mappings (see ES, DE, FR for examples).
		Scenarios: []*tax.ScenarioSet{ bill.InvoiceScenarios() },

		// Corrections defines which correction types are allowed. Omit entirely if
		// the regime regulation doesn't explicitly constrain this: an absent
		// Corrections field means any correction type is accepted.
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
					// bill.InvoiceTypeDebitNote,
					// bill.InvoiceTypeCorrective,
				},
			},
		},
		Validator:  Validate,
		Normalizer: Normalize,

		// Less common fields (uncomment as needed):
		//
		// Identities defines non-tax identity types (e.g. business registration
		// numbers, personal IDs) that can appear on invoices as alternatives to
		// the tax ID. Only needed when the country allows non-tax identifiers —
		// e.g. Sweden accepts an org number instead of a VAT number. Most regimes
		// can omit this. Implement in org_identities.go; see regimes/se/ for
		// reference.
		//
		// Identities: identityTypeDefinitions,
		//
		//
		// AltCountryCodes: []l10n.Code{},          // e.g. GB uses "XI", "XU"
		// CalculatorRoundingRule: tax.RoundingRuleCurrency, // default is "precise"
		// Tags: []*tax.TagSet{},                   // only if regime defines custom tags
		// Extensions: []*cbc.Definition{},         // only if regime defines extension keys
		// PaymentMeansKeys: []*cbc.Definition{},   // regime-specific payment method keys
		// InboxKeys: []*cbc.Definition{},          // e.g. for e-invoice routing
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	// Less common cases — uncomment as needed:
	// case *org.Identity:
	//     return validateIdentity(obj)  // only if Identities field is used
	// case *org.Item:
	//     return validateItem(obj)     // e.g. India validates HSN codes on items
	// case *tax.Combo:
	//     return validateCombo(obj)    // e.g. Portugal validates tax combos
	}
	return nil
}

// Normalize will perform any regime specific normalization.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	// Less common cases — uncomment as needed:
	// case *org.Identity:
	//     normalizeIdentity(obj)       // e.g. Sweden normalizes org numbers
	// case *bill.Invoice:
	//     normalizeInvoice(obj)        // e.g. Mexico sets default extensions
	// case *org.Party:
	//     normalizeParty(obj)          // e.g. Italy extracts fiscal code from tax ID
	// case *tax.Combo:
	//     normalizeCombo(obj)          // e.g. Portugal migrates legacy rate keys
	}
}
