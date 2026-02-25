// Package il provides the tax region definition for Israel.
package il

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the tax region definition for IL.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "IL",
		Currency:  currency.ILS,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Israel",
			i18n.HE: "ישראל",
		},
		TimeZone: "Asia/Jerusalem",
		Scenarios: []*tax.ScenarioSet{
			invoiceScenarios,
		},
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
			},
		},
		Validator:  Validate,
		Normalizer: Normalize,
		Categories: taxCategories,
	}
}

// Validate checks the document type to determine if validation is required.
// Note that, under the IL tax regime, validation of the tax identity verifies
// only the 9-digit numeric format of the Mispar Osek Murshe. Full verification
// of whether a number is registered with the ITA must be performed directly
// through the official government entity register at gov.il.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Normalize attempts to clean up the object passed to it
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
