// Package ro implements the tax regime for Romania.
package ro

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"	
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "RO",
		Currency:  currency.RON,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Romania",
			i18n.RO: "România",
		},			
		TimeZone:  "Europe/Bucharest",				
		Validator:  Validate,
		Normalizer: Normalize,
		Categories: taxCategories,
	}
}

// Normalize will perform any regime specific calculations.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc any) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *bill.Invoice:
		return validateInvoice(obj)
	}
	return nil
}
