// Package br provides the tax region definition for Brazil.
package br

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

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:  "BR",
		Currency: currency.BRL,
		Name: i18n.String{
			i18n.EN: "Brazil",
			i18n.PT: "Brasil",
		},
		Description: i18n.String{
			i18n.EN: "Tax identification in Brazil is provided either through a CNPJ for businesses or a CPF for individuals. Both types are valid for the issuance of NFS-e (electronic service invoices).",
		},
		TimeZone:   "America/Sao_Paulo",
		Validator:  Validate,
		Normalizer: Normalize,
		Categories: taxCategories,
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
	case *bill.Invoice:
		return validateInvoice(obj)
	case *org.Party:
		return validateParty(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Normalize will attempt to clean the object passed to it.
func Normalize(doc interface{}) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
