// Package ae provides the tax region definition for United Arab Emirates.
package ae

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// Local tax category definition which is not considered standard.
const (
	TaxCategoryExcise cbc.Code = "EXCISE"
)

// Specific tax rate codes.
const (
	TaxRateSmokingProducts  cbc.Key = "smoking-product"
	TaxRateCarbonatedDrinks cbc.Key = "carbonated-drink"
	TaxRateEnergyDrinks     cbc.Key = "energy-drink"
	TaxRateSweetenedDrinks  cbc.Key = "sweetened-drink"
)

// New provides the tax region definition for UAE.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:  "AE",
		Currency: currency.AED,
		Name: i18n.String{
			i18n.EN: "United Arab Emirates",
			i18n.AR: "الإمارات العربية المتحدة",
		},
		TimeZone: "Asia/Dubai",
		Tags: []*tax.TagSet{
			common.InvoiceTags(),
		},
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

// Validate checks the document type and determines if it can be validated.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
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
