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
	KeyATTaxExemptionCode cbc.Key = "at-tax-exemption-code"
	KeyATInvoiceType      cbc.Key = "at-invoice-type"
)

// Zone code definitions for Portugal based on districts and
// autonomous regions based on ISO 3166-2:PT.
const (
	ZoneAveiro         l10n.Code = "01"
	ZoneBeja           l10n.Code = "02"
	ZoneBraga          l10n.Code = "03"
	ZoneBraganca       l10n.Code = "04"
	ZoneCasteloBranco  l10n.Code = "05"
	ZoneCoimbra        l10n.Code = "06"
	ZoneEvora          l10n.Code = "07"
	ZoneFaro           l10n.Code = "08"
	ZoneGuarda         l10n.Code = "09"
	ZoneLeiria         l10n.Code = "10"
	ZoneLisboa         l10n.Code = "11"
	ZonePortalegre     l10n.Code = "12"
	ZonePorto          l10n.Code = "13"
	ZoneSantarem       l10n.Code = "14"
	ZoneSetubal        l10n.Code = "15"
	ZoneVianaDoCastelo l10n.Code = "16"
	ZoneVilaReal       l10n.Code = "17"
	ZoneViseu          l10n.Code = "18"
	ZoneAzores         l10n.Code = "20" // Autonomous Region
	ZoneMadeira        l10n.Code = "30" // Autonomous Region
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
		Zones:      zones,
		Tags:       invoiceTags,
		Scenarios:  scenarios,
		Validator:  Validate,
		Calculator: Calculate,
		Preceding: &tax.PrecedingDefinitions{
			Types: []cbc.Key{
				bill.InvoiceTypeCreditNote,
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
	case *tax.Identity:
		return normalizeTaxIdentity(obj)
	}
	return nil
}
