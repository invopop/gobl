// Package co handles tax regime data for Colombia.
package co

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegime(New())
}

// DIAN official codes to include in stamps.
const (
	StampProviderDIANCUDE cbc.Key = "dian-cude"
	StampProviderDIANQR   cbc.Key = "dian-qr"
)

// Special keys to use in meta data.
const (
	KeyDIAN                    cbc.Key = "dian"
	KeyDIANCompanyID           cbc.Key = "dian-company-id"
	KeyDIANAdditionalAccountID cbc.Key = "dian-additional-account-id"
)

// New provides the tax region definition
func New() *tax.Regime {
	return &tax.Regime{
		Country:  "CO",
		Currency: "COP",
		Name: i18n.String{
			i18n.EN: "Colombia",
			i18n.ES: "Colombia",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				The Colombian tax regime is based on the DIAN (Dirección de Impuestos y Aduanas Nacionales)
				specifications for electronic invoicing.
			`),
		},
		TimeZone:     "America/Bogota",
		Tags:         common.InvoiceTags(),
		Validator:    Validate,
		Calculator:   Calculate,
		IdentityKeys: identityKeyDefs, // see identities.go
		Extensions:   extensionKeys,   // see extensions.go
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
				Extensions: []cbc.Key{
					ExtKeyDIANCreditCode,
					ExtKeyDIANDebitCode,
				},
				ReasonRequired: true,
				Stamps: []cbc.Key{
					StampProviderDIANCUDE,
				},
			},
		},
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

// Calculate will attempt to clean the object passed to it.
func Calculate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return normalizeTaxIdentity(obj)
	case *org.Party:
		return normalizeParty(obj)
	}
	return nil
}
