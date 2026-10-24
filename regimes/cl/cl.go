// Package cl provides the tax region definition for Chile.
package cl

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the tax country code for Chile.
const CountryCode = "CL"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("cl", rules.GOBL.Add(CountryCode),
		taxIdentityRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(normalizeTaxIdentity)),
	)
}

// New provides the tax region definition for Chile.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.CLP,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Chile",
			i18n.ES: "Chile",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Chile's tax system is administered by the Servicio de Impuestos Internos
				(SII). IVA (Impuesto al Valor Agregado), Chile's value-added tax,
				applies at the general rate defined in Decreto Ley 825.

				Taxpayers are identified by their RUT (Rol Unico Tributario), which
				includes a modulo-11 check digit. This regime normalizes and validates
				RUT values when used as Chilean tax identities.

				Electronic tax documents (Documentos Tributarios Electronicos, DTE)
				include electronic invoices, credit notes, and debit notes. This regime
				provides the core tax metadata and correction types needed by GOBL;
				DTE XML generation, SII transmission, folio authorization, CAF, and
				electronic stamps are format-specific concerns left for a future addon.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("SII - Servicio de Impuestos Internos"),
				URL:   "https://www.sii.cl/",
			},
		},
		TimeZone: "America/Santiago",
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
