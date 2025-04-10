// Package fi provides the tax region definition for Finland.
package fi

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

// Local tax category definitions specific to Finland
const (
	// Personal income tax
	// source of truth : https://finlex.fi/fi/lainsaadanto/2024/701#sec_1__subsec_1
	TaxCategoryTulovero cbc.Code = "TUL"
)

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "FI",
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Finland",
			i18n.FI: "Suomi",
		},
		TimeZone: "Europe/Helsinki",
		Tags: []*tax.TagSet{
			common.InvoiceTags(),
		},
		Categories: taxCategories,
		Validator:  Validate,
		Normalizer: Normalize,
		Scenarios: []*tax.ScenarioSet{
			common.InvoiceScenarios(),
		},
	}
}

// Validate Finnish invoice and tax identity.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Normalize will attempt to clean the object passed to it.
func Normalize(doc interface{}) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	}
}
