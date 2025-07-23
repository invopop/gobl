// Package ae provides the tax region definition for United Arab Emirates.
package ae

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

// New provides the tax region definition for AE.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "AE",
		Currency:  currency.AED,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "United Arab Emirates",
			i18n.AR: "الإمارات العربية المتحدة",
		},
		TimeZone: "Asia/Dubai",
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

// Validate function assesses the document type to determine if validation is required.
// Note that, under the AE tax regime, validation of the supplier's tax ID is not necessary if it does not meet the specified threshold (refer to the README section for more details).
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Normalize attempts to clean up the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)

	}
}
