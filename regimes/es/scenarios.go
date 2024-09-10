package es

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

// Universal tax tags
const (
	TagCopy             cbc.Key = "copy"
	TagSummary          cbc.Key = "summary"
	TagSimplifiedScheme cbc.Key = "simplified-scheme"
	TagCustomerIssued   cbc.Key = "customer-issued"
	TagTravelAgency     cbc.Key = "travel-agency"
	TagSecondHandGoods  cbc.Key = "second-hand-goods"
	TagArt              cbc.Key = "art"
	TagAntiques         cbc.Key = "antiques"
	TagCashBasis        cbc.Key = "cash-basis"
	TagFacturaE         cbc.Key = "facturae"
	TagTicketBAI        cbc.Key = "ticketbai"
)

var invoiceTags = common.InvoiceTagsWith([]*cbc.KeyDefinition{
	// Copy of the original document
	{
		Key: TagCopy,
		Name: i18n.String{
			i18n.EN: "Copy",
			i18n.ES: "Copia",
		},
	},
	// Summary document
	{
		Key: TagSummary,
		Name: i18n.String{
			i18n.EN: "Summary",
			i18n.ES: "Recapitulativa",
		},
	},
	// Simplified Scheme (Modules)
	{
		Key: TagSimplifiedScheme,
		Name: i18n.String{
			i18n.EN: "Simplified tax scheme",
			i18n.ES: "Contribuyente en régimen simplificado",
		},
	},
	{
		Key: TagFacturaE,
		Name: i18n.String{
			i18n.EN: "FacturaE",
			i18n.ES: "FacturaE",
		},
		Desc: i18n.String{
			i18n.EN: "Invoice must comply with the FacturaE standard.",
			i18n.ES: "La factura debe cumplir con el estándar FacturaE.",
		},
	},
	{
		Key: TagTicketBAI,
		Name: i18n.String{
			i18n.EN: "TicketBAI",
			i18n.ES: "TicketBAI",
		},
		Desc: i18n.String{
			i18n.EN: "Invoice must comply with the TicketBAI specifications.",
			i18n.ES: "La factura debe cumplir con el estándar TicketBAI.",
		},
	},

	// Travel agency
	{
		Key: TagTravelAgency,
		Name: i18n.String{
			i18n.EN: "Special scheme for travel agencies",
			i18n.ES: "Régimen especial de las agencias de viajes",
		},
	},
	// Secondhand stuff
	{
		Key: TagSecondHandGoods,
		Name: i18n.String{
			i18n.EN: "Special scheme for second-hand goods",
			i18n.ES: "Régimen especial de los bienes usados",
		},
	},
	// Art
	{
		Key: TagArt,
		Name: i18n.String{
			i18n.EN: "Special scheme of works of art",
			i18n.ES: "Régimen especial de los objetos de arte",
		},
	},
	// Antiques
	{
		Key: TagAntiques,
		Name: i18n.String{
			i18n.EN: "Special scheme of antiques and collectables",
			i18n.ES: "Régimen especial de las antigüedades y objetos de colección",
		},
	},
	// Special Regime of "Cash Criteria"
	{
		Key: TagCashBasis,
		Name: i18n.String{
			i18n.EN: "Special scheme on cash basis",
			i18n.ES: "Régimen especial del criterio de caja",
		},
	},
})

var invoiceScenarios = &tax.ScenarioSet{
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
			Tags: []cbc.Key{TagFacturaE},
			Ext: tax.Extensions{
				ExtKeyFacturaEDocType: "FC", // default
			},
		},
		{
			Tags: []cbc.Key{
				TagFacturaE,
				tax.TagSimplified,
			},
			Ext: tax.Extensions{
				ExtKeyFacturaEDocType: "FA",
			},
		},
		{
			Tags: []cbc.Key{
				TagFacturaE,
				tax.TagSelfBilled,
			},
			Ext: tax.Extensions{
				ExtKeyFacturaEDocType: "AF",
			},
		},
		// ** Invoice Class **
		{
			Types: []cbc.Key{
				bill.InvoiceTypeStandard,
			},
			Tags: []cbc.Key{TagFacturaE},
			Ext: tax.Extensions{
				ExtKeyFacturaEInvoiceClass: "OO", // Original Invoice
			},
		},
		{
			Types: []cbc.Key{
				bill.InvoiceTypeCorrective,
				bill.InvoiceTypeCreditNote,
				bill.InvoiceTypeDebitNote,
			},
			Tags: []cbc.Key{TagFacturaE},
			Ext: tax.Extensions{
				ExtKeyFacturaEInvoiceClass: "OR", // Corrective
			},
		},
		{
			Tags: []cbc.Key{TagFacturaE, TagSummary},
			Ext: tax.Extensions{
				ExtKeyFacturaEInvoiceClass: "OC", // Summary
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{TagFacturaE, TagCopy},
			Ext: tax.Extensions{
				ExtKeyFacturaEInvoiceClass: "CO", // Copy of the original
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCorrective},
			Tags:  []cbc.Key{TagFacturaE, TagCopy},
			Ext: tax.Extensions{
				ExtKeyFacturaEInvoiceClass: "CR", // Copy of the corrective
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{TagFacturaE, TagCopy, TagSummary},
			Ext: tax.Extensions{
				ExtKeyFacturaEInvoiceClass: "CC", // Copy of the summary
			},
		},
		// ** Special Messages **
		// Reverse Charges
		{
			Tags: []cbc.Key{tax.TagReverseCharge},
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  tax.TagReverseCharge,
				Text: "Reverse Charge / Inversión del sujeto pasivo.",
			},
		},
		// Simplified Scheme (Modules)
		{
			Tags: []cbc.Key{TagSimplifiedScheme},
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  TagSimplifiedScheme,
				Text: "Factura expedida por contibuyente en régimen simplificado.",
			},
		},
		// Customer issued invoices
		{
			Tags: []cbc.Key{tax.TagSelfBilled},
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  tax.TagSelfBilled,
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
