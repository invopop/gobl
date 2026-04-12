// Package uy provides the tax regime definition for Uruguay.
package uy

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the tax region definition for Uruguay.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "UY",
		Currency:  currency.UYU,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Uruguay",
			i18n.ES: "Uruguay",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Uruguay's tax system is administered by the DGI (Dirección General
				Impositiva). The primary indirect tax is the IVA (Impuesto al Valor
				Agregado), which applies at a standard rate and a reduced rate known
				as the "tasa mínima".

				Taxpayers are identified by their RUT (Registro Único Tributario),
				a 12-digit number that includes a check digit calculated using a
				modulo 11 algorithm. The RUT serves as both the general taxpayer
				identification and the IVA registration number.

				Exports are zero-rated. Certain goods and services are exempt from
				IVA, including some financial services and agricultural products.

				Electronic invoicing (Comprobantes Fiscales Electrónicos, CFE) is
				mandatory for all IVA taxpayers, administered through the DGI.
				Both credit notes and debit notes are supported for invoice
				corrections.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("DGI - RUT numbering"),
				URL:   "https://www.gub.uy/direccion-general-impositiva/comunicacion/noticias/nueva-numeracion-del-rut",
			},
			{
				Title: i18n.NewString("OECD - Tax Identification Numbers: Uruguay"),
				URL:   "https://www.oecd.org/content/dam/oecd/en/topics/policy-issue-focus/aeoi/uruguay-tin.pdf",
			},
			{
				Title: i18n.NewString("python-stdnum RUT validation"),
				URL:   "https://arthurdejong.org/python-stdnum/doc/1.20/stdnum.uy.rut",
			},
		},
		TimeZone:   "America/Montevideo",
		Validator:  Validate,
		Normalizer: Normalize,
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
