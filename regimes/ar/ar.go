// Package ar provides the tax region definition for Argentina.
package ar

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the tax region definition for Argentina
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:  "AR",
		Currency: currency.ARS,
		Name: i18n.String{
			i18n.EN: "Argentina",
			i18n.ES: "Argentina",
		},
		Description: i18n.String{
			i18n.EN: "Argentina's tax system is administered by ARCA (Agencia de Recaudación y Control Aduanero), which oversees the collection of IVA (Impuesto al Valor Agregado), the country's value-added tax. Taxpayers are identified using CUIT (Clave Única de Identificación Tributaria) numbers, which serve as unique tax identifiers for both individuals and businesses.",
		},
		TimeZone:    "America/Argentina/Buenos_Aires",
		Validator:   Validate,
		Normalizer:  Normalize,
		Categories:  taxCategories(),
		Corrections: correctionDefinitions(),
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

// Normalize will attempt to clean the object passed to it.
func Normalize(doc interface{}) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	}
}
