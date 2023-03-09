package es

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		// ** Invoice Document Types **
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Meta:  cbc.Meta{KeyFacturaEInvoiceDocumentType: "FC"},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeSimplified},
			Meta:  cbc.Meta{KeyFacturaEInvoiceDocumentType: "FA"},
		},
		{
			Tags: []cbc.Key{common.TagSelfBilled}, // duplicated with notes
			Meta: cbc.Meta{KeyFacturaEInvoiceDocumentType: "AF"},
		},
		// ** Invoice Class **
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard, bill.InvoiceTypeSimplified},
			Meta: cbc.Meta{
				KeyFacturaEInvoiceClass: "OO", // Original Invoice
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCorrective},
			Meta: cbc.Meta{
				KeyFacturaEInvoiceClass: "OR", // Corrective
			},
		},
		{
			Tags: []cbc.Key{TagSummary},
			Meta: cbc.Meta{
				KeyFacturaEInvoiceClass: "OC", // Summary
			},
		},
		{
			Tags:  []cbc.Key{TagCopy},
			Types: []cbc.Key{bill.InvoiceTypeStandard, bill.InvoiceTypeSimplified},
			Meta: cbc.Meta{
				KeyFacturaEInvoiceClass: "CO", // Copy of the original
			},
		},
		{
			Tags:  []cbc.Key{TagCopy},
			Types: []cbc.Key{bill.InvoiceTypeCorrective},
			Meta: cbc.Meta{
				KeyFacturaEInvoiceClass: "CR", // Copy of the corrective
			},
		},
		{
			Tags:  []cbc.Key{TagCopy, TagSummary},
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Meta: cbc.Meta{
				KeyFacturaEInvoiceClass: "CC", // Copy of the summary
			},
		},
		// ** Special Messages **
		// Reverse Charges
		{
			Tags: []cbc.Key{common.TagReverseCharge},
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  common.TagReverseCharge,
				Text: "Reverse Charge / Inversión del sujeto pasivo.",
			},
		},
		// Simplified Regime (Modules)
		{
			Tags: []cbc.Key{TagSimplified},
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  TagSimplified,
				Text: "Factura expedida por contibuyente en régimen simplificado.",
			},
		},
		// Customer issued invoices
		{
			Tags: []cbc.Key{common.TagSelfBilled},
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  common.TagSelfBilled,
				Text: "Facturación por el destinatario.",
			},
		},
		// Travel agency
		{
			Tags: []cbc.Key{TagTravelAgency},
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  TagTravelAgency,
				Text: "Régimen especial de las agencias de viajes.",
			},
		},
		// Secondhand stuff
		{
			Tags: []cbc.Key{TagSecondHandGoods},
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  TagSecondHandGoods,
				Text: "Régimen especial de los bienes usados.",
			},
		},
		// Art
		{
			Tags: []cbc.Key{TagArt},
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  TagArt,
				Text: "Régimen especial de los objetos de arte.",
			},
		},
		// Antiques
		{
			Tags: []cbc.Key{TagAntiques},
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  TagAntiques,
				Text: "Régimen especial de las antigüedades y objetos de colección.",
			},
		},
		// Special Regime of "Cash Criteria"
		{
			Tags: []cbc.Key{TagCashBasis},
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  TagCashBasis,
				Text: "Régimen especial del criterio de caja.",
			},
		},
	},
}
