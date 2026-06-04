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
				Norway applies VAT (merverdiavgift, MVA), administered by the
				Norwegian Tax Administration (Skatteetaten), at a standard rate
				with reduced and special rates for certain goods and services.

				Businesses are identified by their organisasjonsnummer, validated
				with a mod-11 check digit; the VAT number is the organisation
				number followed by "MVA".
			`),
		},
		TimeZone:   "Europe/Oslo",
		Identities: identityTypeDefinitions,
		Categories: taxCategories,
		// Only the kreditnota (credit note) exists in Norwegian law (§ 5-2-7).
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
