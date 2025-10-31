package ar

import (
	"github.com/invopop/gobl/org"
)

// InvoiceLegalNoteExamples defines legal notes that may be required or commonly used
// in Argentine invoices according to AFIP regulations and Argentine tax law.
//
// These examples are provided to help developers understand common legal note requirements.
// Users should consult with local accountants or tax advisors to ensure compliance with
// specific business situations and current regulations.
//
// References:
// - IVA Law 23.349 and modifications
// - AFIP Resolution 4540/2019 (Electronic invoicing)
// - Provincial tax codes for Ingresos Brutos
//
// Sources:
// - https://www.afip.gob.ar/facturacion/regimen-general/
// - https://www.afip.gob.ar/monotributo/
var InvoiceLegalNoteExamples = map[string]*org.Note{
	// VAT Exempt operation
	"vat-exempt": {
		Key:  org.NoteKeyLegal,
		Text: "Operación exenta de IVA conforme al Artículo [indicar artículo aplicable] de la Ley 23.349.",
	},

	// Monotributo - not valid as tax credit
	"monotributo": {
		Key:  org.NoteKeyLegal,
		Text: "Comprobante emitido bajo el Régimen Simplificado para Pequeños Contribuyentes (Monotributo). No válido como crédito fiscal.",
	},

	// Export operation
	"export": {
		Key:  org.NoteKeyLegal,
		Text: "Operación de exportación. IVA tasa cero conforme al Artículo 43 de la Ley 23.349.",
	},

	// Services to foreign customers
	"export-services": {
		Key:  org.NoteKeyLegal,
		Text: "Exportación de servicios. Operación exenta de IVA conforme al Artículo 43 de la Ley 23.349.",
	},

	// VAT retention applied
	"vat-retention": {
		Key:  org.NoteKeyLegal,
		Text: "Se aplicó retención de IVA conforme a la Resolución General AFIP 2854/2010 y modificatorias.",
	},

	// Income tax withholding applied
	"income-withholding": {
		Key:  org.NoteKeyLegal,
		Text: "Se aplicó retención de Impuesto a las Ganancias conforme a la Resolución General AFIP 830/2000 y modificatorias.",
	},

	// Gross income tax (provincial)
	"gross-income": {
		Key:  org.NoteKeyLegal,
		Text: "Sujeto a Impuesto sobre los Ingresos Brutos - Provincia de [indicar provincia]. Alícuota aplicable: [indicar porcentaje]%.",
	},

	// Final consumer - no tax credit
	"final-consumer": {
		Key:  org.NoteKeyLegal,
		Text: "IVA incluido. Documento no válido como crédito fiscal.",
	},

	// Reverse charge
	"reverse-charge": {
		Key:  org.NoteKeyLegal,
		Text: "Inversión del sujeto pasivo. El cliente es responsable del pago del IVA.",
	},

	// Credit note reference
	"credit-note": {
		Key:  org.NoteKeyLegal,
		Text: "Nota de crédito por [motivo: devolución/descuento/bonificación/anulación]. Comprobante original: [indicar tipo y número].",
	},

	// Debit note reference
	"debit-note": {
		Key:  org.NoteKeyLegal,
		Text: "Nota de débito por [motivo: intereses/gastos adicionales/diferencias de precio]. Comprobante original: [indicar tipo y número].",
	},

	// Electronic authorization code
	"cae": {
		Key:  org.NoteKeyLegal,
		Text: "Comprobante Autorizado. CAE N°: [número]. Fecha de vencimiento: [fecha].",
	},

	// Responsible taxpayer declaration
	"responsable-inscripto": {
		Key:  org.NoteKeyLegal,
		Text: "IVA Responsable Inscripto.",
	},

	// VAT exempt entity
	"exento": {
		Key:  org.NoteKeyLegal,
		Text: "IVA Sujeto Exento.",
	},

	// Not responsible for VAT
	"no-responsable": {
		Key:  org.NoteKeyLegal,
		Text: "No Responsable de IVA.",
	},

	// Simplified invoice disclaimer
	"simplified": {
		Key:  org.NoteKeyLegal,
		Text: "Comprobante simplificado conforme a normativa AFIP. Válido como comprobante no fiscal.",
	},

	// Book exemption (common for educational materials)
	"books-exempt": {
		Key:  org.NoteKeyLegal,
		Text: "Operación exenta de IVA - Libros, diarios, revistas y publicaciones periódicas según Ley 23.349 Art. 7° inc. h).",
	},

	// Medical services exemption
	"medical-exempt": {
		Key:  org.NoteKeyLegal,
		Text: "Operación exenta de IVA - Servicios de medicina, incluidos laboratorios de análisis clínicos según Ley 23.349 Art. 7° inc. e).",
	},

	// Educational services exemption
	"education-exempt": {
		Key:  org.NoteKeyLegal,
		Text: "Operación exenta de IVA - Servicios de enseñanza según Ley 23.349 Art. 7° inc. h).",
	},

	// Financial services exemption
	"financial-exempt": {
		Key:  org.NoteKeyLegal,
		Text: "Operación exenta de IVA - Servicios financieros según Ley 23.349 Art. 7° inc. h).",
	},

	// Interest on loans exemption
	"interest-exempt": {
		Key:  org.NoteKeyLegal,
		Text: "Operación exenta de IVA - Intereses de préstamos según Ley 23.349 Art. 7° inc. h).",
	},

	// Insurance exemption
	"insurance-exempt": {
		Key:  org.NoteKeyLegal,
		Text: "Operación exenta de IVA - Servicios de seguros según Ley 23.349 Art. 7° inc. h).",
	},
}
