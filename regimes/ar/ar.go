// Package ar provides the tax region definition for Argentina.
package ar

import (
	"github.com/invopop/gobl/bill"
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
				Argentina's tax system is administered by ARCA (Agencia de Recaudación y Control Aduanero), which oversees the collection of IVA (Impuesto al Valor Agregado), the country's value-added tax.

				Taxpayers are identified using three main types of tax identification numbers: CUIT (Clave Única de Identificación Tributaria) - 11 digits (XX-XXXXXXXX-X) used by companies and legal entities with prefixes 30, 33 (conflict resolution), and 34 (foreign entities); CUIL (Clave Única de Identificación Laboral) - 11 digits (XX-XXXXXXXX-X) used by individuals and employees with prefixes 20 (males), 27 (females), and 23 (conflict resolution); and CDI (Clave de Identificación) for foreign residents without CUIT/CUIL (not validated).

				IVA rates include 27% increased rate for gas, water and telecom services, 21% general rate for most goods and services, and 10.5% reduced rate for essential goods such as construction, medicine, transportation, and food products.

				Argentina applies several retention taxes: IVA Retenido (Retained VAT) with variable rates based on taxpayer category and registration status (reference: AFIP RG 2854/2010 and modifications); Ganancias (Income Tax Withholding) applied to payments for services with rates ranging from 0.5% to 35% depending on service type (reference: AFIP RG 830/2000, RG 4003/2017); and Ingresos Brutos (Gross Income Tax), a provincial tax with rates set by each jurisdiction, typically 1% to 5% depending on province and activity.

				Electronic invoicing through ARCA's system is required for most transactions. Electronic invoices must include CAE/CAI (Código de Autorización Electrónico - Electronic Authorization Code) and Point of Sale (Punto de Venta) for invoice numbering.

				Common invoice types include Tipo A (issued by Responsable Inscripto to another Responsable Inscripto), Tipo B (issued by Responsable Inscripto to Monotributista or final consumer), Tipo C (issued by Monotributista or exempt entities), Tipo E (export invoices), and Credit Notes (Notas de Crédito - corrective documents).

				Argentina has different tax regime classifications: Responsable Inscripto (registered taxpayer with full IVA obligations), Monotributo (simplified tax regime for small businesses), Exento (exempt from IVA), No Responsable (not responsible for IVA collection), and Consumidor Final (final consumer - no tax ID required).
			`),
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
