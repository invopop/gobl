// Package sa provides the tax regime definition for the Kingdom of Saudi Arabia.
package sa

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the tax country code for Saudi Arabia.
const CountryCode = "SA"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("sa", rules.GOBL.Add(CountryCode),
		orgIdentityRules(),
		taxIdentityRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(normalizeTaxIdentity)),
	)
}

// New provides the tax regime definition for Saudi Arabia.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.SAR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Kingdom of Saudi Arabia",
			i18n.AR: "المملكة العربية السعودية",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Saudi Arabia's tax system is administered by the Zakat, Tax and
				Customs Authority (ZATCA).

				B2C invoices follow the reporting flow while B2B invoices follow 
				the clearance flow. Both credit notes and debit notes are supported
				for invoice corrections and follow the same flow as the invoices 
				they correct.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("ZATCA - VAT Implementing Regulations"),
				URL:   "https://zatca.gov.sa/en/RulesRegulations/VAT/Pages/default.aspx",
			},
		},
		TimeZone:   "Asia/Riyadh",
		Categories: taxCategories(),
	}
}
