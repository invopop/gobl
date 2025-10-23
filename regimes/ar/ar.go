// Package ar provides the tax region definition for Argentina.
package ar

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
			i18n.EN: "The Argentine tax system is administered by AFIP (Administración Federal de Ingresos Públicos). Tax identification in Argentina is provided through CUIT (Clave Única de Identificación Tributaria) for businesses and individuals, CUIL (Clave Única de Identificación Laboral) for employees, or CDI (Clave de Identificación) for foreign residents.",
		},
		TimeZone:   "America/Argentina/Buenos_Aires",
		Validator:  Validate,
		Normalizer: Normalize,
		Identities: identityDefinitions(),
		Categories: taxCategories(),
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

// Validate checks the document type and determines if it can be validated.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
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
