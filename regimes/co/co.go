// Package co handles tax regime data for Colombia.
package co

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:  "CO",
		Currency: "COP",
		Name: i18n.String{
			i18n.EN: "Colombia",
			i18n.ES: "Colombia",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Colombia's tax system is administered by the DIAN (Dirección de Impuestos y
				Aduanas Nacionales). Electronic invoicing is mandatory for most businesses
				through the DIAN's e-invoicing platform.

				Businesses are identified by their NIT (Número de Identificación Tributaria),
				which includes a check digit. For B2C transactions (using the simplified tag),
				customers may be identified using various document types including Registro
				civil, Tarjeta de identidad, Cédula de ciudadanía, Tarjeta de extranjería,
				Cédula de extranjería, Pasaporte, PEP, or NUIP. If no customer identity is
				provided for simplified invoices, the reserved final consumer code is used
				automatically.

				IVA (Impuesto sobre el Valor Agregado) rates include a 19% general rate for
				most goods and services, a 5% reduced rate for certain goods, and 0% for
				exports and certain basic goods. Some goods and services are excluded or
				exempt from IVA.

				Invoice series must be pre-registered with the DIAN. Municipality codes are
				required for addresses. Both credit notes and debit notes are supported for
				invoice corrections, each requiring a specific correction cause code from the
				DIAN (e.g. partial refund, revoked, discount, adjustment for credit notes;
				interest, pending charges, change in value for debit notes).
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("DIAN - Invoice numbering authorization"),
				URL:   "https://www.dian.gov.co/impuestos/sociedades/presentacionclientes/Solicitud_de_Autorizacion_de_Numeracion_de_Facturacion.pdf",
			},
			{
				Title: i18n.NewString("DIAN - Municipality codes"),
				URL:   "https://www.dian.gov.co/atencionciudadano/formulariosinstructivos/Formularios/2007/Codigos_municipios_2007.pdf",
			},
		},
		TimeZone: "America/Bogota",
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
