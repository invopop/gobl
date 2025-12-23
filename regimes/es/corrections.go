package es

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

// InvoiceCorrectionTypes defines the types of corrections recognized in Spanish law
var InvoiceCorrectionTypes = []cbc.Key{
	bill.InvoiceTypeCreditNote,
	bill.InvoiceTypeDebitNote,
}

func correctionDefinitions() []*tax.CorrectionDefinition {
	return []*tax.CorrectionDefinition{
		{
			Schema: bill.ShortSchemaInvoice,
			Types:  InvoiceCorrectionTypes,
		},
	}
}
