package zatca

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
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
			Types: []cbc.Key{
				bill.InvoiceTypeStandard,
			},
			Ext: tax.Extensions{
				untdid.ExtKeyDocumentType: "388",
			},
		},
		{
			Types: []cbc.Key{
				bill.InvoiceTypeStandard,
			},
			Tags: []cbc.Key{tax.TagPrepayment},
			Ext: tax.Extensions{
				untdid.ExtKeyDocumentType: "386",
			},
		},
		{
			Types: []cbc.Key{
				bill.InvoiceTypeCreditNote,
			},
			Ext: tax.Extensions{
				untdid.ExtKeyDocumentType: "381",
			},
		},
		{
			Types: []cbc.Key{
				bill.InvoiceTypeDebitNote,
			},
			Ext: tax.Extensions{
				untdid.ExtKeyDocumentType: "383",
			},
		},
	},
}
