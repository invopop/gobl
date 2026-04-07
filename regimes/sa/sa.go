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

const countryCode = "SA"

// Identification keys used for additional codes not covered by the standard fields
const (
	IdentityTypeTIN      cbc.Code = "TIN" // Tax Identification Number
	IdentityTypeCRN      cbc.Code = "CRN" // Commercial Registration Number
	IdentityTypeMom      cbc.Code = "MOM" // Ministry of Municipal, Rural Affairs and Housing Number
	IdentityTypeMLS      cbc.Code = "MLS" // Ministry of Human Resources and Social Development Number
	IdentityType700      cbc.Code = "700" // 700 Number
	IdentityTypeSAG      cbc.Code = "SAG" // Saudi Arabian General Authority Number
	IdentityTypeNational cbc.Code = "NAT" // National ID
	IdentityTypeGcc      cbc.Code = "GCC" // GCC ID
	IdentityTypeIqa      cbc.Code = "IQA" // Iqama Number (Resident ID)
	IdentityTypePassport cbc.Code = "PAS" // Passport Number
	IdentityTypeOTH      cbc.Code = "OTH" // Other ID
)

var supplierValidIdentities = []cbc.Code{
	IdentityTypeCRN,
	IdentityTypeMom,
	IdentityTypeMLS,
	IdentityType700,
	IdentityTypeSAG,
	IdentityTypeOTH,
}

var customerValidIdentities = []cbc.Code{
	IdentityTypeTIN,
	IdentityTypeCRN,
	IdentityTypeMom,
	IdentityTypeMLS,
	IdentityType700,
	IdentityTypeSAG,
	IdentityTypeNational,
	IdentityTypeGcc,
	IdentityTypeIqa,
	IdentityTypePassport,
	IdentityTypeOTH,
}

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("sa", rules.GOBL.Add(countryCode),
		billInvoiceRules(),
		orgIdentityRules(),
		taxIdentityRules(),
	)
}

// New provides the tax regime definition for Saudi Arabia.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   countryCode,
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
		Normalizer: Normalize,
		Categories: taxCategories(),
	}
}

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	}
}
