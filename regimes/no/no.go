// Package no provides the tax regime definition for Norway.
package no

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

// New instantiates a new Norwegian tax regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   l10n.NO.Tax(),
		Currency:  currency.NOK,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Norway",
			i18n.NB: "Norge",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				The Norwegian tax regime covers VAT (merverdiavgift) with four
				rates: general (25%), reduced (15%), super-reduced (12%), and
				special (11.11%). Identity validation supports
				organisasjonsnummer with mod-11 check digits.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Skatteetaten - VAT rates"),
				URL:   "https://www.skatteetaten.no/en/rates/value-added-tax/",
			},
			{
				Title: i18n.NewString("Brønnøysundregistrene - Organisation number"),
				URL:   "https://www.brreg.no/en/about-us-2/our-registers/about-the-central-coordinating-register-for-legal-entities-ccr/about-the-organisation-number/",
			},
			{
				Title: i18n.NewString("Lovdata - Merverdiavgiftsloven"),
				URL:   "https://lovdata.no/dokument/NL/lov/2009-06-19-58",
			},
		},
		TimeZone:   "Europe/Oslo",
		Identities: identityTypeDefinitions,
		Categories: taxCategories,
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
			},
		},
		Scenarios: []*tax.ScenarioSet{
			invoiceScenarios,
		},
		Validator:  Validate,
		Normalizer: Normalize,
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateBillInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *org.Identity:
		return validateOrgIdentity(obj)
	}
	return nil
}

// Normalize performs any regime-specific normalizations.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	case *org.Identity:
		normalizeOrgIdentity(obj)
	}
}
