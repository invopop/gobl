// Package pl provides the Polish tax regime.
package pl

import (
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("pl", rules.GOBL.Add("PL"),
		billInvoiceRules(),
		taxIdentityRules(),
	)
}

// CountryCode is the tax country code for Poland.
const CountryCode l10n.TaxCountryCode = "PL"

// New instantiates a new Polish regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.PLN,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Poland",
			i18n.PL: "Polska",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Poland's tax system is administered by the Krajowa Administracja Skarbowa
				(National Revenue Administration, KAS). As an EU member state, Poland follows
				the EU VAT Directive with locally adapted rates.

				VAT (Podatek od towarów i usług, PTU) applies at standard, reduced, and
				super-reduced rates. Zero-rated supplies include exports and intra-community
				supplies.

				Businesses are identified by their NIP (Numer Identyfikacji Podatkowej), a
				10-digit tax identification number. The Polish VAT number uses the format PL
				followed by the 10-digit NIP.

				Poland has implemented the KSeF (Krajowy System e-Faktur) national e-invoicing
				system, which is progressively becoming mandatory for B2B transactions.
				E-invoicing via PEPPOL is used for cross-border and B2G transactions.
			`),
		},
		TimeZone:   "Europe/Warsaw",
		Normalizer: Normalize,
		Categories: taxCategories, // tax_categories.go
	}
}

// Normalize will perform any regime specific normalizations.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
