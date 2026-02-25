// Package be defines the tax regime data for Belgium
package be

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

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "BE",
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Belgium",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Belgium's tax system is administered by the Federal Public Service Finance
				(Service Public Fédéral Finances / Federale Overheidsdienst Financiën).
				As an EU member state, Belgium follows the EU VAT Directive.

				VAT (Taxe sur la Valeur Ajoutée, TVA / Belasting over de Toegevoegde Waarde,
				BTW) rates include a 21% standard rate for most goods and services, a 12%
				intermediate rate for certain goods including social housing, restaurant
				services, and some food products, and a 6% reduced rate for basic necessities
				such as food, water, pharmaceuticals, books, and passenger transport.

				Businesses are identified by their VAT number (Numéro de TVA / BTW-nummer)
				in the format BE followed by 10 digits. Belgium supports credit notes for
				invoice corrections.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("BOSA - Electronic Invoicing"),
				URL:   "https://bosa.belgium.be/fr/themes/administration-numerique/facturation-electronique",
			},
		},
		TimeZone:   "Europe/Brussels",
		Validator:  Validate,
		Normalizer: Normalize,
		Scenarios: []*tax.ScenarioSet{
			bill.InvoiceScenarios(),
		},
		Categories: taxCategories,
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
			},
		},
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
		tax.NormalizeIdentity(obj)
	}
}
