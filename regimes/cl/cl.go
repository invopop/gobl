// Package cl provides the Chile tax regime.
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
	rules.Register("cl", rules.GOBL.Add(CountryCode), taxIdentityRules())
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(func(id *tax.Identity) { tax.NormalizeIdentity(id) })),
	)
}

// New provides the tax regime definition for Chile.
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
				Chile's tax system is administered by the Servicio de Impuestos
				Internos (SII). Businesses and individuals are identified by their
				RUT (Rol Único Tributario), a numeric code followed by a single
				check digit (0-9 or K) calculated with a modulo 11 algorithm.

				VAT (Impuesto al Valor Agregado, IVA) is a single, flat-rate tax
				applied uniformly to sales and services, with no reduced or
				intermediate rates.

				Chile's electronic invoicing system (Documentos Tributarios
				Electrónicos, DTE) requires invoices to be digitally signed and
				sent to the SII for clearance before being delivered to the
				customer. That clearance flow — including certificate
				management, CAF folio ranges, and the SII's SOAP/XML web
				services — is a separate integration concern from the core tax
				rules modelled here, and is intentionally out of scope for this
				regime. It would be a natural candidate for a future GOBL addon
				(similar to how es/tbai or es/verifactu extend the Spanish
				regime), not part of the base regime itself.

				Both credit notes and debit notes are supported for invoice
				corrections, mirroring the SII's Nota de Crédito Electrónica and
				Nota de Débito Electrónica document types.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "SII - What is the VAT rate?",
					i18n.ES: "SII - ¿Cuál es la tasa del Impuesto al Valor Agregado (IVA)?",
				},
				URL: "https://www.sii.cl/preguntas_frecuentes/impuestos_mensuales/001_130_0572.htm",
			},
		},
		TimeZone: "America/Santiago",
		Scenarios: []*tax.ScenarioSet{
			bill.InvoiceScenarios(),
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
