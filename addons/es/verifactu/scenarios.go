package verifactu

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

var scenarios = []*tax.ScenarioSet{
	{
		Schema: bill.ShortSchemaInvoice,
		List: []*tax.Scenario{
			// ** Invoice Document Types **
			{
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "F1",
				},
			},
			{
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
				},
				Tags: []cbc.Key{
					tax.TagSimplified,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "F2",
				},
			},
		},
	},
}
