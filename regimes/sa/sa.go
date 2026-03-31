// Package sa provides the tax regime definition for the Kingdom of Saudi Arabia.
package sa

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("sa", rules.GOBL.Add("SA"),
		taxIdentityRules(),
		billInvoiceRules(),
	)
}

// New provides the tax regime definition for Saudi Arabia.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "SA",
		Currency:  currency.SAR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Kingdom of Saudi Arabia",
			i18n.AR: "المملكة العربية السعودية",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Saudi Arabia's tax system is administered by the Zakat, Tax and
				Customs Authority (ZATCA). VAT was introduced on January 1, 2018,
				at 5% and increased to 15% on July 1, 2020.

				Businesses must register for VAT if annual taxable supplies exceed
				SAR 375,000, with voluntary registration available above SAR 187,500.
				Registered businesses receive a VAT Identification Number which must
				appear on all tax invoices.

				Simplified tax invoices may be used for B2C transactions. Both credit
				notes and debit notes are supported for invoice corrections.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("ZATCA - VAT Implementing Regulations"),
				URL:   "https://zatca.gov.sa/en/RulesRegulations/VAT/Pages/default.aspx",
			},
		},
		TimeZone: "Asia/Riyadh",
		Scenarios: []*tax.ScenarioSet{
			invoiceScenarios,
		},
		Corrections: correctionDefinitions(),
		Normalizer:  Normalize,
		Categories:  taxCategories(),
	}
}

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	}
}
