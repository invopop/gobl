// Package co handles tax regime data for Colombia.
package co

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// Special keys to use in meta data.
const (
	KeyDIAN                    cbc.Key = "dian"
	KeyDIANCompanyID           cbc.Key = "dian-company-id"
	KeyDIANAdditionalAccountID cbc.Key = "dian-additional-account-id"
)

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:  "CO",
		Currency: "COP",
		Name: i18n.String{
			i18n.EN: "Colombia",
			i18n.ES: "Colombia",
		},
		TimeZone: "America/Bogota",
		Tags: []*tax.TagSet{
			common.InvoiceTags(),
		},
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

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	case *org.Party:
		normalizeParty(obj)
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
