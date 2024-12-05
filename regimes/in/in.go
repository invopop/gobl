// Package in provides models for dealing with India.
package in

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the tax region definition for India.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:  "IN",
		Currency: currency.INR,
		Name: i18n.String{
			i18n.EN: "India",
		},
		TimeZone: "Asia/Kolkata",
		Tags: []*tax.TagSet{
			common.InvoiceTags().Merge(invoiceTags),
		},
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
		Validator:  Validate,
		Normalizer: Normalize,
		Categories: taxCategories,
	}
}

// Validate function assesses the document type to determine if validation is required.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *org.Identity:
		return validateOrgIdentity(obj)
	case *org.Item:
		return validateOrgItem(obj)
	}
	return nil
}

// Normalize attempts to clean up the object passed to it.
func Normalize(doc interface{}) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	case *org.Identity:
		normalizeOrgIdentity(obj)
	}
}
