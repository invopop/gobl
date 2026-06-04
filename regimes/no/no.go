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
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the tax country code for Norway.
const CountryCode l10n.TaxCountryCode = "NO"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("no", rules.GOBL.Add("NO"),
		taxIdentityRules(),
		orgIdentityRules(),
		billInvoiceRules(),
	)
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
		TimeZone:   "Europe/Oslo",
		Identities: identityTypeDefinitions,
		Categories: taxCategories,
		// Norwegian bookkeeping law (bokføringsforskriften § 5-2-7) only
		// recognises the kreditnota (credit note) for correcting an issued
		// sales document; there is no debit-note concept.
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
			},
		},
		Scenarios: []*tax.ScenarioSet{
			invoiceScenarios,
		},
		Normalizer: Normalize,
	}
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
