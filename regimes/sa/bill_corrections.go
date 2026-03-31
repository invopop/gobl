package sa

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

func correctionDefinitions() []*tax.CorrectionDefinition {
	return []*tax.CorrectionDefinition{
		{
			Schema: bill.ShortSchemaInvoice,
			Types: []cbc.Key{
				bill.InvoiceTypeCreditNote,
				bill.InvoiceTypeDebitNote,
			},
		},
	}
}
