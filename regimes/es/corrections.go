package es

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

var correctionTypes = []cbc.Key{
	bill.InvoiceTypeCorrective,
	bill.InvoiceTypeCreditNote,
	bill.InvoiceTypeDebitNote,
}

var correctionDefinitions = []*tax.CorrectionDefinition{
	{
		Schema: bill.ShortSchemaInvoice,
		Types:  correctionTypes,
		Extensions: []cbc.Key{
			ExtKeyFacturaECorrection,
			ExtKeyTBAICorrection,
		},
	},
}
