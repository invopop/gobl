// Package it provides the Italian tax regime.
package it

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New instantiates a new Italian regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "IT",
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Italy",
			i18n.IT: "Italia",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Italy's tax system is administered by the Agenzia delle Entrate (Revenue
				Agency). All invoices must comply with the FatturaPA electronic format,
				transmitted through the Sistema di Interscambio (SDI).

				IVA (Imposta sul Valore Aggiunto) applies at standard, reduced, intermediate,
				and minimum rates covering various categories of goods and services.

				Businesses are identified by their Partita IVA (VAT number), an 11-digit code,
				and by the Codice Fiscale (fiscal code) for individuals and entities. The
				FatturaPA format requires a Codice Destinatario (recipient code) or PEC
				(certified email) for invoice routing. Every supplier must declare a fiscal
				regime (Regime Fiscale, e.g. RF01 Ordinary, RF19 Flat rate) in their invoices.

				The FatturaPA format supports an extensive set of document types (TD01-TD28)
				covering standard invoices, self-billed invoices, and various special cases.
				Line items may require Nature (Natura) codes to explain VAT exemptions or
				reverse charge situations. Stamp duty (Imposta di bollo) applies to certain
				exempt invoices. Withholding taxes (IRPEF, IRES, INPS, ENASARCO, ENPAM) can
				be applied alongside VAT. Both credit notes and debit notes are supported for
				corrections.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Agenzia delle Entrate - Electronic Invoicing"),
				URL:   "https://www.agenziaentrate.gov.it/portale/web/guest/fatturazione-elettronica-e-dati-fatture-transfrontaliere-new",
			},
			{
				Title: i18n.NewString("FatturaPA - Documentation and Schemas"),
				URL:   "https://www.fatturapa.gov.it/it/norme-e-regole/documentazione-fattura-elettronica/formato-fatturapa/",
			},
			{
				Title: i18n.NewString("FatturaPA - Filling Guide"),
				URL:   "https://www.agenziaentrate.gov.it/portale/documents/20143/451259/Guida_compilazione-FE-Esterometro-V_1.9_2024-03-05.pdf",
			},
		},
		TimeZone: "Europe/Rome",
		Identities: identityKeyDefinitions, // identities.go
		Scenarios:  scenarios,              // scenarios.go
		Validator:  Validate,
		Normalizer: Normalize,
		Categories: categories, // categories.go
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
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *org.Identity:
		return validateIdentity(obj)
	}
	return nil
}

// Normalize will perform any regime specific calculations.
func Normalize(doc interface{}) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	case *org.Identity:
		normalizeIdentity(obj)
	case *org.Party:
		normalizeParty(obj)
	case *tax.Combo:
		normalizeTaxCombo(obj)
	}
}
