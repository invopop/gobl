// Package pt provides models for dealing with the Portuguese tax regime.
package pt

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// Custom keys used typically in meta information
const (
	KeyATTaxCountryRegion cbc.Key = "at-tax-country-region"
	KeyATTaxCode          cbc.Key = "at-tax-code"
	KeyATTaxExemptionCode cbc.Key = "at-tax-exemption-code"
	KeyATInvoiceType      cbc.Key = "at-invoice-type"
)

// AT official codes to include in stamps.
const (
	StampProviderATATCUD     cbc.Key = "at-atcud"
	StampProviderATQR        cbc.Key = "at-qr"
	StampProviderATHash      cbc.Key = "at-hash"
	StampProviderATHashFull  cbc.Key = "at-hash-full"
	StampProviderATAppID     cbc.Key = "at-app-id"
	StampProviderATTimestamp cbc.Key = "at-ts"
)

// New instantiates a new Portugal regime for the given zone.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "PT",
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Portugal",
			i18n.PT: "Portugal",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Portugal's tax system is administered by the Autoridade Tributária e Aduaneira
				(AT). As an EU member state, Portugal follows the EU VAT Directive with
				locally adapted rates that vary by region.

				IVA (Imposto sobre o Valor Acrescentado) rates on the mainland include a 23%
				standard rate, a 13% intermediate rate for food and beverages in restaurants,
				and a 6% reduced rate for basic food items, books, and pharmaceuticals. The
				autonomous regions of Açores and Madeira apply reduced rates (16%/9%/4% and
				22%/12%/5% respectively).

				Businesses are identified by their NIF (Número de Identificação Fiscal), a
				9-digit number. The Portuguese VAT number uses the format PT followed by the
				NIF.

				Portugal requires all invoicing software to be certified by the AT and
				invoices must include a unique document identifier (ATCUD) and a hash chain
				linking sequential documents. The SAF-T (Standard Audit File for Tax Purposes)
				format is used for tax reporting. Both credit notes and debit notes are
				supported for invoice corrections.
			`),
		},
		TimeZone:   "Europe/Lisbon",
		Extensions: extensionKeys,
		Validator:  Validate,
		Normalizer: Normalize,
		Tags: []*tax.TagSet{
			invoiceTags,
		},
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
			},
		},
		Categories: taxCategories,
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *tax.Combo:
		return validateTaxCombo(obj)
	}
	return nil
}

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		migrateInvoiceRates(obj)
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	case *tax.Combo:
		normalizeTaxCombo(obj)
	}
}
