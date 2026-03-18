// Package dk provides a regime definition for Denmark.
package dk

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("dk", rules.GOBL.Add("DK"),
		billInvoiceRules(),
		orgIdentityRules(),
		taxIdentityRules(),
	)
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
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Denmark's tax system is administered by the Danish Tax Agency (Skattestyrelsen).
				As an EU member state, Denmark follows the EU VAT Directive.

				VAT (Moms, short for Merværdiafgift) applies at a single standard rate on
				most goods and services. Unlike many other EU countries, Denmark does not
				have reduced VAT rates, making it one of the simplest VAT systems in Europe.
				Certain supplies are zero-rated (e.g. exports, newspapers) or exempt (e.g.
				healthcare, education, financial services).

				Businesses are identified by their CVR number (Det Centrale Virksomhedsregister),
				an 8-digit number. The Danish VAT number uses the format DK followed by the
				8-digit CVR number. E-invoicing via the NemHandel/PEPPOL network is mandatory
				for all B2G transactions.
			`),
		},
		TimeZone:   "Europe/Copenhagen",
		Identities: identityTypeDefinitions,
		Categories: taxCategories,
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
			},
		},
		Normalizer: Normalize,
	}
}

// Normalize will perform any regime specific normalization.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
