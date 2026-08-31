// Package cr provides the tax regime definition for Costa Rica.
package cr

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

// CountryCode is the tax country code for Costa Rica.
const CountryCode = "CR"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("cr", rules.GOBL.Add(CountryCode),
		taxIdentityRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(normalizeTaxIdentity)),
	)
}

// New provides the tax regime definition for Costa Rica.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.CRC,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Costa Rica",
			i18n.ES: "Costa Rica",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Costa Rica's tax system is administered by the Dirección General
				de Tributación (DGT) of the Ministerio de Hacienda. IVA (Impuesto
				al Valor Agregado) has applied since 1 July 2019, when Title I of
				Law 9635 reformed the former general sales tax law (Law 6826)
				into a value added tax; ISC (Impuesto Selectivo de Consumo) is an
				ad valorem excise on selected goods.

				Taxpayers are identified by the same documents used for civil
				identification (cédula física, cédula jurídica, DIMEX, or NITE),
				none of which carry a check digit. Electronic invoicing is
				mandatory for most taxpayers, and issued documents may only be
				corrected with electronic credit or debit notes.
			`),
			i18n.ES: here.Doc(`
				El sistema tributario de Costa Rica es administrado por la
				Dirección General de Tributación (DGT) del Ministerio de
				Hacienda. El IVA (Impuesto al Valor Agregado) se aplica desde el
				1 de julio de 2019, cuando el Título I de la Ley 9635 reformó la
				anterior ley del impuesto general sobre las ventas (Ley 6826) en
				un impuesto al valor agregado; el ISC (Impuesto Selectivo de
				Consumo) es un impuesto ad valorem sobre mercancías
				seleccionadas.

				Los contribuyentes se identifican con los mismos documentos de
				identificación civil (cédula física, cédula jurídica, DIMEX o
				NITE), ninguno de los cuales incluye dígito verificador. La
				factura electrónica es obligatoria para la mayoría de los
				contribuyentes y los comprobantes emitidos solo pueden corregirse
				mediante notas de crédito o débito electrónicas.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Ministerio de Hacienda - Dirección General de Tributación"),
				URL:   "https://www.hacienda.go.cr/",
			},
			{
				Title: i18n.NewString("OECD - Costa Rica Tax Identification Numbers"),
				URL:   "https://www.oecd.org/content/dam/oecd/en/topics/policy-issue-focus/aeoi/costa-rica-tin.pdf",
			},
			{
				Title: i18n.NewString("Reglamento de Comprobantes Electrónicos para Efectos Tributarios"),
				URL:   "https://sinalevi.go.cr/ResultadosNormativa/Informacion?param1=103206&param2=143152&param3=1",
			},
			{
				Title: i18n.NewString("Anexos y Estructuras de Comprobantes Electrónicos v4.4"),
				URL:   "https://www.hacienda.go.cr/docs/ANEXOS_Y_ESTRUCTURAS_V4.4.pdf",
			},
		},
		TimeZone:   "America/Costa_Rica",
		Categories: taxCategories,
		// Issued documents are immutable and may only be corrected with
		// electronic credit or debit notes (Reglamento de Comprobantes
		// Electrónicos, art. 10).
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
