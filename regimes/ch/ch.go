// Package ch provides the Swiss tax regime.
package ch

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("ch", rules.GOBL.Add("CH"), taxIdentityRules())
}

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "CH",
		Currency:  currency.CHF,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Switzerland",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Switzerland's tax system is administered by the Federal Tax Administration
				(Eidgenössische Steuerverwaltung, ESTV). Although not an EU member, Switzerland
				operates its own VAT system (Mehrwertsteuer, MWST).

				VAT applies at standard, reduced, and special rates. The reduced rate covers
				everyday goods such as food, non-alcoholic beverages, books, newspapers, and
				medicines, while a special rate applies to accommodation services.

				Businesses with annual taxable revenues exceeding CHF 100,000 must register
				for VAT. Tax identification uses the UID (Unternehmens-Identifikationsnummer)
				in the format CHE-XXX.XXX.XXX followed by "MWST" for VAT purposes.

				Switzerland supports credit notes for invoice corrections. E-invoicing is not
				mandatory but is increasingly used, particularly in business-to-government (B2G)
				transactions.
			`),
		},
		TimeZone:   "Europe/Zurich",
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

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	}
}
