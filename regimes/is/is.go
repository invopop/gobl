// Package is provides the tax regime definition for Iceland.
package is

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the tax country code for Iceland.
const CountryCode = "IS"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("is", rules.GOBL.Add(CountryCode),
		taxIdentityRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(normalizeTaxIdentity)),
	)
}

// New provides the tax regime definition for Iceland.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.ISK,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Iceland",
			i18n.IS: "Ísland",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Iceland's Value Added Tax (virðisaukaskattur, VSK), administered by
				Skatturinn (Iceland Revenue and Customs), is applied to most goods
				and services at a standard rate, with a reduced rate for
				accommodation, books, food and certain cultural services.

				Businesses are identified by their Kennitala, a 10-digit national
				identifier whose validity is confirmed with a weighted modulus-11
				check digit; the same number serves as the tax identity for both
				individuals and legal entities.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Skatturinn - Value Added Tax"),
				URL:   "https://www.skatturinn.is/english/companies/value-added-tax/",
			},
		},
		TimeZone:   "Atlantic/Reykjavik",
		Categories: taxCategories(),
	}
}
