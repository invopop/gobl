// Package lu provides the Luxembourg tax regime.
package lu

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the tax country code for Luxembourg.
const CountryCode = "LU"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("lu", rules.GOBL.Add(CountryCode), taxIdentityRules())
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(func(id *tax.Identity) { tax.NormalizeIdentity(id) })),
	)
}

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Luxembourg",
			i18n.FR: "Luxembourg",
			i18n.DE: "Luxemburg",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Luxembourg's tax system is administered by the Administration de
				l'enregistrement, des domaines et de la TVA (AED). As an EU member
				state, Luxembourg follows the EU VAT Directive and applies the
				lowest standard VAT rate in the Union.

				VAT (taxe sur la valeur ajoutée, TVA) applies at a standard rate of
				17%, an intermediate rate of 14%, a reduced rate of 8%, and a
				super-reduced rate of 3%. A temporary one-point reduction applied
				to the standard, intermediate, and reduced rates during 2023.

				Businesses are identified by their n° TVA (VAT identification
				number) in the format LU followed by 8 digits, of which the last
				two are check digits. Luxembourg supports credit notes for invoice
				corrections. B2G e-invoicing via PEPPOL is supported.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("AED - Value Added Tax (TVA)"),
				URL:   "https://pfi.public.lu/fr/professionnel/tva.html",
			},
			{
				Title: i18n.NewString("Guichet.lu - VAT rates"),
				URL:   "https://guichet.public.lu/en/entreprises/fiscalite/tva/assujettissement/taux-tva.html",
			},
		},
		TimeZone: "Europe/Luxembourg",
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
