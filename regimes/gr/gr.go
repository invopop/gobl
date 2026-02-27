// Package gr provides the tax region definition for Greece.
package gr

import (
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

// Official IAPR codes to include in stamps.
const (
	StampIAPRQR       cbc.Key = "iapr-qr"
	StampIAPRMark     cbc.Key = "iapr-mark"
	StampIAPRHash     cbc.Key = "iapr-hash"
	StampIAPRUID      cbc.Key = "iapr-uid"
	StampIAPRProvider cbc.Key = "iapr-provider"
)

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country: "EL",
		AltCountryCodes: []l10n.Code{
			"GR", // regular ISO code
		},
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Greece",
			i18n.EL: "Ελλάδα",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Greece's tax system is administered by the Independent Authority for Public
				Revenue (IAPR / AADE). As an EU member state, Greece follows the EU VAT
				Directive with locally adapted rates.

				FPA (Fóros Prostithémenis Axías) applies at standard, reduced, and
				super-reduced rates. The islands of Leros, Lesbos, Kos, Samos, and Chios
				benefit from a reduction on all standard rates.

				Businesses are identified by their AFM (Arithmós Forologikoú Mitróou), a
				9-digit tax identification number. The Greek VAT number uses the format EL
				followed by the 9-digit AFM.

				Greece uses the myDATA platform for tax reporting, where invoices must be
				classified with specific invoice type codes, VAT category codes, income
				classifications, and exemption codes. Payment method codes must also be
				reported. PEPPOL BIS Billing 3.0 is used for B2G e-invoicing. Credit notes
				are supported for invoice corrections.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("myDATA API Documentation v1.0.7"),
				URL:   "https://www.aade.gr/sites/default/files/2023-10/myDATA%20API%20Documentation_v1.0.7_eng.pdf",
			},
			{
				Title: i18n.NewString("Greek PEPPOL BIS Billing 3.0"),
				URL:   "https://www.gsis.gr/sites/default/files/eInvoice/Instructions%20to%20B2G%20Suppliers%20and%20certified%20PEPPOL%20Providers%20for%20the%20Greek%20PEPPOL%20BIS-EN-%20v1.0.pdf",
			},
		},
		TimeZone:               "Europe/Athens",
		CalculatorRoundingRule: tax.RoundingRuleCurrency,
		Scenarios:              scenarios,
		Corrections:            corrections,
		Validator:              Validate,
		Normalizer:             Normalize,
		Categories:             taxCategories,
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc any) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	}
}
