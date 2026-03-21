// Package cz provides a regime definition for the Czech Republic.
package cz

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New instantiates a new Czech Republic regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   l10n.CZ.Tax(),
		Currency:  currency.CZK,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Czech Republic",
			i18n.CS: "Česká republika",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				The Czech Republic's tax system is administered by the Financial
				Administration of the Czech Republic (Finanční správa ČR). As an EU member
				state, the Czech Republic follows the EU VAT Directive.

				VAT (DPH — Daň z přidané hodnoty) applies at a standard rate of 21% and a
				single reduced rate of 12% (since January 2024, when the previous first
				reduced rate of 15% and second reduced rate of 10% were merged). Certain
				supplies are zero-rated (e.g. exports) or exempt (e.g. healthcare,
				education, financial services).

				Businesses are identified by their DIČ (Daňové identifikační číslo), which
				consists of the prefix CZ followed by 8 to 10 digits. For legal entities
				the DIČ is 8 digits with a modulo-11 checksum; for individuals it is
				derived from the birth number (Rodné číslo) and is 9 or 10 digits.
			`),
		},
		TimeZone:   "Europe/Prague",
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
	}
	return nil
}

// Normalize will perform any regime specific normalization.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
