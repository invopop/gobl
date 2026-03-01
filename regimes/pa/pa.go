// Package pa provides the tax regime definition for Panama.
package pa

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New instantiates a new Panamanian tax regime definition.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   l10n.PA.Tax(),
		Currency:  currency.PAB,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Panama",
			i18n.ES: "Panamá",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Panama's tax system is administered by the DGI (Dirección General de Ingresos).
				Electronic invoicing is managed through the SFEP (Sistema de Factura Electrónica
				de Panamá), which requires XML documents signed with XAdES-BES and submitted
				through an authorized PAC (Proveedor Autorizado de Certificación).

				Businesses and individuals are identified by their RUC (Registro Único de
				Contribuyente) paired with a DV (Dígito Verificador), a 2-digit check digit
				pre-assigned by the DGI. The RUC format varies by taxpayer type: natural persons
				use their cédula number, foreigners use an E-prefix, and legal entities follow
				numeric sequences from the Public Registry.

				ITBMS (Impuesto de Transferencia de Bienes Muebles y Servicios) is Panama's
				principal consumption tax, functioning as a VAT. A selective consumption tax
				(ISC) applies to specific goods and must be calculated before the ITBMS.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("DGI - SFEP Portal"),
				URL:   "https://dgi.mef.gob.pa/facturaelectronica",
			},
			{
				Title: i18n.NewString("DGI - SFEP Technical Documentation"),
				URL:   "https://dgi.mef.gob.pa/FacturaElectronica/Documentacion.html",
			},
			{
				Title: i18n.NewString("Executive Decree 766 (2020) - SFEP Operational Rules"),
				URL:   "https://www.gacetaoficial.gob.pa/pdfTemp/29187_A/82818.pdf",
			},
		},
		TimeZone:   "America/Panama",
		Categories: taxCategories,
		// SFEP supports credit and debit notes per Executive Decree 766 (2020).
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
			},
		},
		Validator:  Validate,
		Normalizer: Normalize,
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

// Normalize will perform any regime specific normalization.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	}
}
