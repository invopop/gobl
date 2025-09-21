// Package au provides the Australian tax regime.
package au

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Name: i18n.String{
			i18n.EN: "Australia",
		},
		Country:    l10n.TaxCountryCode("AU"),
		Currency:   currency.AUD,
		TaxScheme:  tax.CategoryGST,
		TimeZone:   "Australia/Sydney",
		Categories: taxCategories,
		Scenarios: []*tax.ScenarioSet{
			bill.InvoiceScenarios(),
		},
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
			},
		},
		Normalizer: Normalize,
		Validator:  Validate,
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc any) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Normalize will perform any regime specific calculations.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
