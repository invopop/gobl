package flow10

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

// Flow 10 accepts a curated subset of UNTDID 1001 document type codes.
// The scenarios below map GOBL invoice types (+ tag combinations) to the
// corresponding UNTDID code via untdid.ExtKeyDocumentType. The list is
// intentionally self-contained so Flow 10 can operate without requiring
// the en16931 addon.
var scenarios = []*tax.ScenarioSet{
	{
		Schema: bill.ShortSchemaInvoice,
		List: []*tax.Scenario{
			// Simple invoices ---------------------------------------------------
			{
				// 380 — Sales invoice
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Ext: tax.ExtensionsOf(tax.ExtMap{
					untdid.ExtKeyDocumentType: "380",
				}),
			},
			{
				// 389 — Self-billed invoice
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Tags:  []cbc.Key{tax.TagSelfBilled},
				Ext: tax.ExtensionsOf(tax.ExtMap{
					untdid.ExtKeyDocumentType: "389",
				}),
			},
			{
				// 393 — Factored invoice
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Tags:  []cbc.Key{tax.TagFactoring},
				Ext: tax.ExtensionsOf(tax.ExtMap{
					untdid.ExtKeyDocumentType: "393",
				}),
			},
			{
				// 501 — Self-invoiced factored invoice
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Tags:  []cbc.Key{tax.TagSelfBilled, tax.TagFactoring},
				Ext: tax.ExtensionsOf(tax.ExtMap{
					untdid.ExtKeyDocumentType: "501",
				}),
			},

			// Deposit invoices --------------------------------------------------
			{
				// 386 — Deposit invoice
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Tags:  []cbc.Key{tax.TagPrepayment},
				Ext: tax.ExtensionsOf(tax.ExtMap{
					untdid.ExtKeyDocumentType: "386",
				}),
			},
			{
				// 500 — Self-invoiced deposit invoice
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Tags:  []cbc.Key{tax.TagSelfBilled, tax.TagPrepayment},
				Ext: tax.ExtensionsOf(tax.ExtMap{
					untdid.ExtKeyDocumentType: "500",
				}),
			},

			// Corrective invoices -----------------------------------------------
			{
				// 384 — Corrective invoice
				Types: []cbc.Key{bill.InvoiceTypeCorrective},
				Ext: tax.ExtensionsOf(tax.ExtMap{
					untdid.ExtKeyDocumentType: "384",
				}),
			},
			{
				// 471 — Self-billed corrective invoice
				Types: []cbc.Key{bill.InvoiceTypeCorrective},
				Tags:  []cbc.Key{tax.TagSelfBilled},
				Ext: tax.ExtensionsOf(tax.ExtMap{
					untdid.ExtKeyDocumentType: "471",
				}),
			},
			{
				// 472 — Factored corrective invoice
				Types: []cbc.Key{bill.InvoiceTypeCorrective},
				Tags:  []cbc.Key{tax.TagFactoring},
				Ext: tax.ExtensionsOf(tax.ExtMap{
					untdid.ExtKeyDocumentType: "472",
				}),
			},
			{
				// 473 — Self-billed factored corrective invoice
				Types: []cbc.Key{bill.InvoiceTypeCorrective},
				Tags:  []cbc.Key{tax.TagSelfBilled, tax.TagFactoring},
				Ext: tax.ExtensionsOf(tax.ExtMap{
					untdid.ExtKeyDocumentType: "473",
				}),
			},

			// Credit memos ------------------------------------------------------
			{
				// 381 — Credit memo
				Types: []cbc.Key{bill.InvoiceTypeCreditNote},
				Ext: tax.ExtensionsOf(tax.ExtMap{
					untdid.ExtKeyDocumentType: "381",
				}),
			},
			{
				// 261 — Self-billed credit memo
				Types: []cbc.Key{bill.InvoiceTypeCreditNote},
				Tags:  []cbc.Key{tax.TagSelfBilled},
				Ext: tax.ExtensionsOf(tax.ExtMap{
					untdid.ExtKeyDocumentType: "261",
				}),
			},
			{
				// 396 — Factored credit memo
				Types: []cbc.Key{bill.InvoiceTypeCreditNote},
				Tags:  []cbc.Key{tax.TagFactoring},
				Ext: tax.ExtensionsOf(tax.ExtMap{
					untdid.ExtKeyDocumentType: "396",
				}),
			},
			{
				// 502 — Self-invoiced and factored credit memo
				Types: []cbc.Key{bill.InvoiceTypeCreditNote},
				Tags:  []cbc.Key{tax.TagSelfBilled, tax.TagFactoring},
				Ext: tax.ExtensionsOf(tax.ExtMap{
					untdid.ExtKeyDocumentType: "502",
				}),
			},
			{
				// 503 — Down-payment invoice credit memo
				Types: []cbc.Key{bill.InvoiceTypeCreditNote},
				Tags:  []cbc.Key{tax.TagPrepayment},
				Ext: tax.ExtensionsOf(tax.ExtMap{
					untdid.ExtKeyDocumentType: "503",
				}),
			},
		},
	},
}

// allowedDocumentTypes is the whitelist of UNTDID 1001 codes permitted on a
// Flow 10 invoice (B2B scope). Kept in sync with the scenarios above.
var allowedDocumentTypes = []cbc.Code{
	"380", "389", "393", "501",
	"386", "500",
	"384", "471", "472", "473",
	"381", "261", "396", "502", "503",
}
