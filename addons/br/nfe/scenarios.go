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
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard, bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote},
			Ext: tax.Extensions{
				ExtKeyModel: ModelNFe,
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard, bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote},
			Tags:  []cbc.Key{tax.TagSimplified},
			Ext: tax.Extensions{
				ExtKeyModel: ModelNFCe,
			},
		},
	},
}
