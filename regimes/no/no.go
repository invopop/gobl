// Package no provides a regime definition for Norway.
package no

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"

	"github.com/invopop/gobl/cbc"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New instantiates a new Norwegian regime.
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
				Norway's tax system is administered by the Norwegian Tax Administration
				(Skatteetaten). While not an EU member state, Norway follows similar VAT
				rules through the EEA Agreement.

				VAT (Merverdiavgift, MVA) applies at a standard rate and two reduced rates.
				The standard rate covers most goods and services, one reduced rate applies
				to foodstuffs and water supply, and another reduced rate covers passenger
				transport, accommodation, broadcasting, and cultural events. Certain
				supplies are zero-rated (e.g. exports, newspapers) or exempt (e.g.
				healthcare, education, financial services).

				Businesses are identified by their organization number
				(Organisasjonsnummer), a 9-digit number. The Norwegian VAT number uses the
				format NO followed by the 9-digit organization number and the suffix MVA.
				E-invoicing via the PEPPOL network is mandatory for all B2G transactions.
			`),
		},
		TimeZone:   "Europe/Oslo",
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
		Validator:  Validate,
		Normalizer: Normalize,
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc any) error {
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

// Normalize will perform any regime specific normalization.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	case *org.Identity:
		normalizeOrgIdentity(obj)
	}
}
