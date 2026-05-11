package ctc

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

// scenarios is a self-contained UNTDID 1001 document-type mapping for
// French CTC invoices. It does not assume eu-en16931 is also declared
// (en16931 is only mandatory for Flow 2 clearance and is enforced via
// the rule set, not as an addon dependency), so we duplicate the
// document-type rows for the codes that en16931 also covers. When both
// addons are listed the rows merge with the same value — last-write
// is harmless because the value is identical.
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

// allowedInvoiceDocumentTypes is the whitelist of UNTDID 1001 codes
// permitted on a French CTC invoice (covers both Flow 2 and Flow 10).
// Kept as a flat list because the rule that consumes it checks for
// presence/absence rather than the type+tag combination.
var allowedInvoiceDocumentTypes = []cbc.Code{
	"380", "389", "393", "501",
	"386", "500",
	"384", "471", "472", "473",
	"381", "261", "396", "502", "503",
	// Flow 2-only: consolidated credit note. Not driven by a scenario;
	// the caller sets the extension explicitly.
	"262",
}
