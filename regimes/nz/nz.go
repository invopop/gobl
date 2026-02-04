// Package nz provides the New Zealand tax regime.
package nz

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "NZ",
		Currency:  currency.NZD,
		TaxScheme: tax.CategoryGST,
		Name: i18n.String{
			i18n.EN: "New Zealand",
		},
		TimeZone:   "Pacific/Auckland",
		Validator:  Validate,
		Normalizer: Normalize,
		Scenarios: []*tax.ScenarioSet{
			bill.InvoiceScenarios(),
		},
		Categories: taxCategories,
		Identities: append(identityKeyDefinitions, orgIdentityDefinitions...),
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

func Validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *org.Identity:
		return validateIdentity(obj)
	}
	return nil
}

func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	case *org.Identity:
		normalizeIdentity(obj)
	}
}
