package ro

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

// InvoiceCorrectionTypes defines the types of corrections recognized in Romanian law.
//
// According to Romanian legislation, once an invoice is issued and sent to e-Factura,
// it cannot be deleted or simply cancelled. It must be corrected by issuing a specific
// correction document (Credit Note/Debit Note) that references the original.
//
// References:
//
//   - Law 227/2015 (Fiscal Code), Art. 330 "Correction of invoices":
//     https://legislatie.just.ro/Public/DetaliiDocument/171282#id_artA2436_1
//
//   - ANAF e-Factura Technical Guide (Corrections & Storno):
//     https://www.anaf.ro/anaf/internet/ANAF/servicii_online/efactura/ghiduri_tehnice
//
//   - RO_CIUS (National Specification based on EN 16931), Rule BR-RO-K01:
//     "Credit Notes (381) must contain a reference to the billing document being corrected."
var InvoiceCorrectionTypes = []cbc.Key{
	bill.InvoiceTypeCreditNote, // "Factură de stornare" (Negative amounts)
	bill.InvoiceTypeDebitNote,  // "Factură de ajustare" (Positive amounts)
	bill.InvoiceTypeCorrective, // Generic correction
}

func correctionDefinitions() []*tax.CorrectionDefinition {
	return []*tax.CorrectionDefinition{
		{
			Schema: bill.ShortSchemaInvoice,
			Types:  InvoiceCorrectionTypes,
		},
	}
}
