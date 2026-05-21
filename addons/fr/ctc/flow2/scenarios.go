package flow2

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

// scenarios is a self-contained UNTDID 1001 document-type mapping for
// Flow 2 invoices. Mirrors the codes the en16931 base profile also
// covers — last-write is harmless because the value is identical.
var scenarios = []*tax.ScenarioSet{
	{
		Schema: bill.ShortSchemaInvoice,
		List: []*tax.Scenario{
			// Simple invoices ---------------------------------------------
			{
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyDocumentType: "380",
				}),
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Tags:  []cbc.Key{tax.TagSelfBilled},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyDocumentType: "389",
				}),
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Tags:  []cbc.Key{tax.TagFactoring},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyDocumentType: "393",
				}),
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Tags:  []cbc.Key{tax.TagSelfBilled, tax.TagFactoring},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyDocumentType: "501",
				}),
			},
			// Deposit invoices ---------------------------------------------
			{
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Tags:  []cbc.Key{tax.TagPrepayment},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyDocumentType: "386",
				}),
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeStandard},
				Tags:  []cbc.Key{tax.TagSelfBilled, tax.TagPrepayment},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyDocumentType: "500",
				}),
			},
			// Corrective invoices ------------------------------------------
			{
				Types: []cbc.Key{bill.InvoiceTypeCorrective},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyDocumentType: "384",
				}),
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeCorrective},
				Tags:  []cbc.Key{tax.TagSelfBilled},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyDocumentType: "471",
				}),
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeCorrective},
				Tags:  []cbc.Key{tax.TagFactoring},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyDocumentType: "472",
				}),
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeCorrective},
				Tags:  []cbc.Key{tax.TagSelfBilled, tax.TagFactoring},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyDocumentType: "473",
				}),
			},
			// Credit notes -------------------------------------------------
			{
				Types: []cbc.Key{bill.InvoiceTypeCreditNote},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyDocumentType: "381",
				}),
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeCreditNote},
				Tags:  []cbc.Key{tax.TagSelfBilled},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyDocumentType: "261",
				}),
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeCreditNote},
				Tags:  []cbc.Key{tax.TagFactoring},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyDocumentType: "396",
				}),
			},
			{
				Types: []cbc.Key{bill.InvoiceTypeCreditNote},
				Tags:  []cbc.Key{tax.TagSelfBilled, tax.TagFactoring},
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyDocumentType: "502",
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

// allowedInvoiceDocumentTypes is the whitelist of UNTDID 1001 codes
// permitted on a Flow 2 invoice. Includes 262 for consolidated credit
// notes (caller sets the extension explicitly; not driven by a
// scenario).
var allowedInvoiceDocumentTypes = []cbc.Code{
	"380", "389", "393", "501",
	"386", "500",
	"384", "471", "472", "473",
	"381", "261", "396", "502", "503",
	"262",
}

// Self-billed document types (BR-FR-21/22).
var selfBilledDocumentTypes = []cbc.Code{
	"389", "501", "500", "471", "473", "261", "502",
}

// Corrective invoice document types (BR-FR-CO-04).
var correctiveInvoiceTypes = []cbc.Code{
	"384", "471", "472", "473",
}

// Credit note document types (BR-FR-CO-05).
var creditNoteTypes = []cbc.Code{
	"261", "381", "396", "502", "503",
}

// advancePaymentDocumentTypes are the UNTDID 1001 codes representing
// advance-payment invoices (forbidden combined with B4/S4/M4 modes).
var advancePaymentDocumentTypes = []cbc.Code{
	"386", "500", "503",
}
