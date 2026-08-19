// Package cy provides tax regime support for Cyprus.
package cy

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the ISO 3166-1 alpha-2 code for Cyprus.
const CountryCode = "CY"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("cy", rules.GOBL.Add(CountryCode),
		taxIdentityRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(func(id *tax.Identity) { tax.NormalizeIdentity(id) })),
	)
}

// New instantiates a new Cyprus regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Cyprus",
			i18n.EL: "Κύπρος",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Cyprus's tax system is administered by the Cyprus Tax Department. As an
				EU member state, Cyprus follows the EU VAT Directive with locally adapted
				rates.

				VAT applies at standard, reduced, super-reduced, zero, and exempt rates.
				Exports and qualifying intra-EU supplies may be zero-rated.

				Taxpayers are identified by a Tax Identification Number (TIN) consisting
				of eight digits followed by a Latin letter.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Cyprus Tax Department - VAT Invoice Requirements"),
				URL:   "https://www.gov.cy/mof-tax/documents/ekdosi-timologioy-f-p-a/",
			},
		},
		TimeZone:   "Asia/Nicosia",
		Categories: taxCategories,
		Scenarios: []*tax.ScenarioSet{
			invoiceScenarios,
		},
	}
}
