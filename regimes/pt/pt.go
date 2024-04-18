// Package pt provides models for dealing with the Portuguese tax regime.
package pt

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegime(New())
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
	StampProviderATATCUD cbc.Key = "at-atcud"
	StampProviderATQR    cbc.Key = "at-qr"
)

// New instantiates a new Portugal regime for the given zone.
func New() *tax.Regime {
	return &tax.Regime{
		Country:  l10n.PT,
		Currency: currency.EUR,
		Name: i18n.String{
			i18n.EN: "Portugal",
			i18n.PT: "Portugal",
		},
		TimeZone:   "Europe/Lisbon",
		Extensions: extensionKeys,
		Tags:       invoiceTags,
		Scenarios:  scenarios,
		Validator:  Validate,
		Calculator: Calculate,
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
			},
		},
		Categories: taxCategories,
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Calculate will attempt to clean the object passed to it.
func Calculate(doc interface{}) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		if err := migrateTaxIDZoneToLines(obj); err != nil {
			return err
		}
		return migrateInvoiceRates(obj)
	case *tax.Identity:
		return normalizeTaxIdentity(obj)
	}
	return nil
}
