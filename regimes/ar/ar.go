// Package ar provides the tax region definition for Argentina.
package ar

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
			i18n.EN: here.Doc(`
				Argentina's tax system is administered by ARCA (Agencia de Recaudación
				y Control Aduanero), which oversees the collection of IVA (Impuesto al
				Valor Agregado), the country's value-added tax.

				Taxpayers are identified using three main types of tax identification
				numbers: CUIT (Clave Única de Identificación Tributaria), 11 digits
				(XX-XXXXXXXX-X) used by companies and legal entities with prefixes 30,
				33, and 34; CUIL (Clave Única de Identificación Laboral), 11 digits used
				by individuals with prefixes 20, 27, and 23; and CDI (Clave de
				Identificación) for foreign residents without CUIT/CUIL.

				IVA has increased, general, and reduced rates. Argentina also
				applies several retention taxes: IVA Retenido (Retained VAT) with
				variable rates based on taxpayer category; Ganancias (Income Tax
				Withholding) applied to payments for services; and Ingresos Brutos
				(Gross Income Tax), a provincial tax with rates set by each
				jurisdiction.

				Electronic invoicing through ARCA is required for most transactions.
				Invoices must include a CAE (Código de Autorización Electrónico) and
				Point of Sale (Punto de Venta) number.

				Common invoice types include Tipo A (between Responsable Inscripto
				parties), Tipo B (to Monotributista or final consumer), Tipo C (from
				Monotributista or exempt entities), Tipo E (exports), and Credit Notes
				(Notas de Crédito).

				Tax regime classifications include Responsable Inscripto (full IVA
				obligations), Monotributo (simplified regime for small businesses),
				Exento (exempt from IVA), and Consumidor Final (no tax ID required).
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("AFIP - Administración Federal de Ingresos Públicos"),
				URL:   "https://www.afip.gob.ar/",
			},
			{
				Title: i18n.NewString("Invoice type and mandatory information"),
				URL:   "https://www.argentina.gob.ar/normativa/recurso/54461/259-98a/htm",
			},
			{
				Title: i18n.NewString("AFIP SIRE - Percepciones y Retenciones"),
				URL:   "https://www.afip.gob.ar/sire/percepciones-retenciones/",
			},
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
