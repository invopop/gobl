package flow10

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

// scenarios maps the bill.Invoice type+tags combinations Flow 10
// accepts to their UNTDID 1001 document type codes.
var scenarios = []*tax.ScenarioSet{
	{
		Schema: bill.ShortSchemaInvoice,
		List: []*tax.Scenario{
			{
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyDocumentType: "380",
				}),
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Tags:  []cbc.Key{tax.TagPrepayment},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyDocumentType: "386",
				}),
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeCorrective},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyDocumentType: "384",
				}),
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeCreditNote},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyDocumentType: "381",
				}),
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeCreditNote},
				Tags:  []cbc.Key{tax.TagPrepayment},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyDocumentType: "503",
				}),
			},
		},
	},
}

// allowedDocumentTypes is the whitelist of UNTDID 1001 codes
// permitted on a Flow 10 invoice.
var allowedDocumentTypes = []cbc.Code{
	"380", "386", "384", "381", "503",
}
