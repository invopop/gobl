// Package jp provides the tax regime definition for Japan.
package jp

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

// Local tax category definitions which are not considered standard.
const (
	// TaxCategoryWHT is the code for Japan's Withholding Income Tax
	TaxCategoryWHT cbc.Code = "WHT"
)

// Specific withholding tax rate keys.
const (
	TaxRatePro     cbc.Key = "pro"      // Professional services (≤ ¥1,000,000)
	TaxRateProOver cbc.Key = "pro-over" // Professional services (> ¥1,000,000)
)

// New provides the tax regime definition for Japan.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "JP",
		Currency:  currency.JPY,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Japan",
			i18n.JA: "日本",
		},
		TimeZone:   "Asia/Tokyo",
		Identities: identityDefinitions,
		Scenarios: []*tax.ScenarioSet{
			invoiceScenarios,
		},
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
			},
		},
		Categories: taxCategories,
		Validator:  Validate,
		Normalizer: Normalize,
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *org.Identity:
		return validateRegistrationNumber(obj)
	}
	return nil
}

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	case *org.Identity:
		normalizeRegistrationNumber(obj)
	}
}
