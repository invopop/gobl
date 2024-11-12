package verifactu

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
)

const (
	TagSubstitution = "substitution"
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
			{
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
				},
				Tags: []cbc.Key{
					TagSubstitution,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "F3",
				},
			},
			{
				Types: es.InvoiceCorrectionTypes,
				Ext: tax.Extensions{
					ExtKeyDocType: "R1",
				},
			},
			{
				Types: es.InvoiceCorrectionTypes,
				Tags: []cbc.Key{
					tax.TagSimplified,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: "R5",
				},
			},
		},
	},
}
