package nfe

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

var corrections = tax.CorrectionSet{
	{
		Schema: bill.ShortSchemaInvoice,
		Types:  []cbc.Key{bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote},
		Stamps: []cbc.Key{StampProviderSEFAZKey},
	},
}
