// Package pt provides models for dealing with the Portuguese tax regime.
package pt

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/regimes/common"
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
	StampProviderATATCUD    cbc.Key = "at-atcud"
	StampProviderATQR       cbc.Key = "at-qr"
	StampProviderATHash     cbc.Key = "at-hash"
	StampProviderATHashFull cbc.Key = "at-hash-full"
	StampProviderATAppID    cbc.Key = "at-app-id"
)

// New instantiates a new Portugal regime for the given zone.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:  "PT",
		Currency: currency.EUR,
		Name: i18n.String{
			i18n.EN: "Portugal",
			i18n.PT: "Portugal",
		},
		TimeZone:   "Europe/Lisbon",
		Extensions: extensionKeys,
		Validator:  Validate,
		Normalizer: Normalize,
		Tags: []*tax.TagSet{
			common.InvoiceTags().Merge(invoiceTags),
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
		migrateTaxIDZoneToLines(obj)
		migrateInvoiceRates(obj)
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	case *tax.Combo:
		normalizeTaxCombo(obj)
	}
}
