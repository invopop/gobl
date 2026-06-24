package nfe

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

var scenarios = []*tax.ScenarioSet{
	invoiceScenarios,
}

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		// Model
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard, bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyModel: ModelNFe,
			}),
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard, bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote},
			Tags:  []cbc.Key{tax.TagSimplified},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyModel: ModelNFCe,
			}),
		},
		// Purpose & operation type: only the normal/outbound combination is set,
		// and only for standard invoices. Other combinations can't be handled
		// cleanly with scenarios; to use them, the invoice must have a
		// non-standard type and set the extensions manually.
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyPurpose:       PurposeNormal,
				ExtKeyOperationType: OperationOutbound,
			}),
		},
	},
}
