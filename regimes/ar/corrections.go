package ar

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

// Correction types recognized in Argentina according to ARCA regulations.
//
// In Argentina, credit notes (notas de crédito) and debit notes (notas de débito)
// are used to adjust previously issued invoices. These must be linked to the original
// invoice and issued within specific timeframes.
//
// References:
// - AFIP Resolution 4540/2019: Electronic voucher requirements
// - 15 business day limit for issuing corrections
// - Must be linked to original invoice with same recipient
// - Require Electronic Authorization Code (CAE)
//
// Sources:
// - https://www.afip.gob.ar/facturacion/regimen-general/
// - https://www.boletinoficial.gob.ar/detalleAviso/primera/263505/20220531

// InvoiceCorrectionTypes defines the invoice correction types recognized in Argentina
var InvoiceCorrectionTypes = []cbc.Key{
	bill.InvoiceTypeCreditNote,
	bill.InvoiceTypeDebitNote,
}

// correctionDefinitions returns the correction definitions for Argentina
func correctionDefinitions() []*tax.CorrectionDefinition {
	return []*tax.CorrectionDefinition{
		{
			Schema: bill.ShortSchemaInvoice,
			Types:  InvoiceCorrectionTypes,
		},
	}
}
