package xrechnung

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
					bill.InvoiceTypeCorrective,
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
			},
			{
				Tags: []cbc.Key{
					tax.TagSelfBilled,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: tax.ExtValue(invoiceTypeSelfBilled),
				},
			},
			{
				Tags: []cbc.Key{
					tax.TagPartial,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: tax.ExtValue(invoiceTypePartial),
				},
			},
			{
				Tags: []cbc.Key{
					tax.TagPartial,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: tax.ExtValue(invoiceTypePartialConstruction),
				},
			},
			{
				Tags: []cbc.Key{
					tax.TagPartial,
				},
				Ext: tax.Extensions{
					ExtKeyDocType: tax.ExtValue(invoiceTypePartialFinalConstruction),
				},
			},
			{
				Ext: tax.Extensions{
					ExtKeyDocType: tax.ExtValue(invoiceTypeFinalConstruction),
				},
			},
			// ** Tax Rates **
			{
				Ext: tax.Extensions{
					ExtKeyTaxRate: "S",
				},
			},
			{
				Ext: tax.Extensions{
					ExtKeyTaxRate: "Z",
				},
			},
			{
				Ext: tax.Extensions{
					ExtKeyTaxRate: "E",
				},
			},
			{
				Tags: []cbc.Key{
					tax.TagReverseCharge,
				},
				Ext: tax.Extensions{
					ExtKeyTaxRate: "AE",
				},
			},
			{
				Ext: tax.Extensions{
					ExtKeyTaxRate: "K",
				},
			},
			{
				Ext: tax.Extensions{
					ExtKeyTaxRate: "G",
				},
			},
			{
				Ext: tax.Extensions{
					ExtKeyTaxRate: "O",
				},
			},
		},
	},
}
