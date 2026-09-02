// Package ro provides the Romanian tax regime.
package ro

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

// CountryCode is the tax country code for Romania.
const CountryCode = "RO"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("ro", rules.GOBL.Add(CountryCode), taxIdentityRules())
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(func(id *tax.Identity) { tax.NormalizeIdentity(id) })),
	)
}

// New provides the tax region definition for Romania.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.RON,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Romania",
			i18n.RO: "România",
		},

		// Short description (based on pl format)
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Romania's tax system is administered by the Agenția Națională de
				Administrare Fiscală (ANAF). As an EU member state, Romania follows
				the EU VAT Directive.

				VAT (Taxa pe Valoarea Adăugată, TVA) applies to most goods and
				services. Businesses are identified by their Cod de Identificare Fiscală
				(CIF), also known historically as the Cod Unic de Înregistrare
				(CUI). For VAT-registered entities this is prefixed with RO to
				form the Romanian VAT number (e.g. RO12345678).

				Romania has progressively mandated the RO e-Factura national
				e-invoicing system for B2B and B2G transactions; support for this
				format is out of scope for this base regime and is expected to be
				implemented as a separate addon.
			`),
		},

		// Trustworthy sources about tax systems in Romania (ANAF)
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("ANAF - Cod fiscal 227/2015 Cotele de TVA (VAT rates)"),
				URL:   "https://static.anaf.ro/static/10/Anaf/legislatie/Cod_fiscal_norme_2023.htm",
			},
			{
				Title: i18n.NewString("MODIFICĂRI ADUSE LEGII NR. 227/2015 PRIVIND CODUL FISCAL PRIN LEGEA NR. 141/2025*. VAT 19% -> 21%"),
				URL:   "https://static.anaf.ro/static/10/Brasov/Brasov/tva_2025.pdf",
			},
		},
		TimeZone: "Europe/Bucharest",
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
