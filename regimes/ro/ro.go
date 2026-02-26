// Package ro provides the Romanian tax regime.
package ro

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

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "RO",
		Currency:  currency.RON,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Romania",
			i18n.RO: "România",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Romania's tax regime is managed by ANAF (Agenția Națională de
				Administrare Fiscală). VAT rates as of August 2025 are: standard
				21%, reduced 11%, and super-reduced 11%.

				Tax identity uses the CUI (Codul Unic de Înregistrare), a 2-10
				digit number with a weighted checksum, assigned to every business
				upon registration with ONRC regardless of VAT status. The "RO"
				prefix indicates VAT registration and is automatically stripped
				during normalization. Non-VAT-registered businesses (below the
				RON 395,000 annual threshold) still have a CUI and issue invoices
				without VAT.

				B2B e-invoicing is mandatory via the RO e-Factura system since
				January 2024 using UBL 2.1 format. Corrections support both
				credit and debit notes. Simplified invoices and reverse charge
				are available via tags.
			`),
		},
		TimeZone:   "Europe/Bucharest",
		Validator:  Validate,
		Normalizer: Normalize,
		Scenarios: []*tax.ScenarioSet{
			invoiceScenarios,
		},
		Categories: taxCategories,
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
			},
		},
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
