// Package bs provides the tax regime definition for Bahamas.
package bs

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the tax country code for Bahamas.
const CountryCode = "BS"

func init() {
	tax.RegisterRegimeDef(New())                      // regimes.go
	rules.Register("bs", rules.GOBL.Add(CountryCode), // rules.go
		taxIdentityRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(normalizeTaxIdentity)),
	)
}

// New provides the tax regime definition for Bahamas.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.BSD,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Commonwealth of The Bahamas",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Value Added Tax (VAT) was introduced in The Bahamas in January 2015
				as part of a broader reform of the tax system. The broad based consumption
				tax is applied to almost all goods and services that are imported, bought
				and sold for use in The Bahamas. Three rates of VAT are applied which are
				the standard rate of 10%, a reduced rate of 5% and a zero rate of 0%. Goods
				exported to customers abroad are classified as exempt from VAT or zero rated.

				Businesses are identified by their Tax Identification Number, a 9 digit
				identifier.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("PRIME MINISTER DAVIS’S CONTRIBUTION ON THE COMMONWEALTH OF THE BAHAMAS 2025/2026 MID-YEAR BUDGET STATEMENT"),
				URL:   "https://opm.gov.bs/prime-minister-daivs-contribution-2025-2026-mid-year-budget-statement/",
			},
			{
				Title: i18n.NewString("ABOUT - VAT, Department of Inland Revenue"),
				URL:   "https://inlandrevenue.finance.gov.bs/value-added-tax/about-vat/",
			},
		},
		TimeZone:   "America/Nassau",
		Categories: taxCategories(),
	}
}
