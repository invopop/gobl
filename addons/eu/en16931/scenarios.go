package en16931

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

// Scenarios provides a list of scenarios related to the UNTDID addon
// that can be used inside other addons.
func Scenarios() []*tax.ScenarioSet {
	return scenarios
}

var scenarios = []*tax.ScenarioSet{
	{
		Schema: bill.ShortSchemaInvoice,
		List: []*tax.Scenario{
			// ** Invoice Document Type Mappings for most common use cases **
			{
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
				},
				Ext: tax.Extensions{
					untdid.ExtKeyDocumentType: "380",
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
			{
				Types: []cbc.Key{
					bill.InvoiceTypeCorrective,
				},
				Ext: tax.Extensions{
					untdid.ExtKeyDocumentType: "384",
				},
			},
			{
				Types: []cbc.Key{
					bill.InvoiceTypeProforma,
				},
				Ext: tax.Extensions{
					untdid.ExtKeyDocumentType: "325",
				},
			},
			{
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
				},
				Tags: []cbc.Key{
					tax.TagPartial,
				},
				Ext: tax.Extensions{
					untdid.ExtKeyDocumentType: "326",
				},
			},
			{
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
				},
				Tags: []cbc.Key{
					tax.TagSelfBilled,
				},
				Ext: tax.Extensions{
					untdid.ExtKeyDocumentType: "389",
				},
			},
			// https://docs.peppol.eu/poacc/self-billing/3.0/
			{
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
				Tags: []cbc.Key{
					tax.TagSelfBilled,
				},
				Ext: tax.Extensions{
					untdid.ExtKeyDocumentType: "261",
				},
			},
			{
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
				},
				Tags: []cbc.Key{
					tax.TagPrepayment,
				},
				Ext: tax.Extensions{
					untdid.ExtKeyDocumentType: "386",
				},
			},
			{
				Types: []cbc.Key{
					bill.InvoiceTypeStandard,
				},
				Tags: []cbc.Key{
					tax.TagFactored,
				},
				Ext: tax.Extensions{
					untdid.ExtKeyDocumentType: "393",
				},
			},
			{
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
				Tags: []cbc.Key{
					tax.TagFactored,
				},
				Ext: tax.Extensions{
					untdid.ExtKeyDocumentType: "396",
				},
			},
		},
	},
}
