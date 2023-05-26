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
)

var invoiceTags = []*tax.KeyDefinition{
	// Simplified Invoice
	{
		Key: common.TagSimplified,
		Name: i18n.String{
			i18n.EN: "Simplified Invoice",
			i18n.ES: "Factura Simplificada",
		},
	},
	// Customer rates (mainly for digital goods inside EU)
	{
		Key: common.TagCustomerRates,
		Name: i18n.String{
			i18n.EN: "Customer rates",
			i18n.ES: "Tarifas aplicables al destinatario",
		},
	},
	// Reverse Charge Mechanism
	{
		Key: common.TagReverseCharge,
		Name: i18n.String{
			i18n.EN: "Reverse Charge",
			i18n.ES: "Inversión del sujeto pasivo",
		},
	},
	// Customer issued invoices
	{
		Key: common.TagSelfBilled,
		Name: i18n.String{
			i18n.EN: "Customer issued invoice",
			i18n.ES: "Facturación por el destinatario",
		},
	},
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
}

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		// ** Invoice Document Types **
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard, bill.InvoiceTypeCorrective},
			Codes: cbc.CodeSet{
				KeyFacturaEInvoiceDocumentType: "FC", // default
			},
		},
		{
			Tags: []cbc.Key{common.TagSimplified},
			Codes: cbc.CodeSet{
				KeyFacturaEInvoiceDocumentType: "FA",
			},
		},
		{
			Tags: []cbc.Key{common.TagSelfBilled},
			Codes: cbc.CodeSet{
				KeyFacturaEInvoiceDocumentType: "AF",
			},
		},
		// ** Invoice Class **
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Codes: cbc.CodeSet{
				KeyFacturaEInvoiceClass: "OO", // Original Invoice
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCorrective},
			Codes: cbc.CodeSet{
				KeyFacturaEInvoiceClass: "OR", // Corrective
			},
		},
		{
			Tags: []cbc.Key{TagSummary},
			Codes: cbc.CodeSet{
				KeyFacturaEInvoiceClass: "OC", // Summary
			},
		},
		{
			Tags:  []cbc.Key{TagCopy},
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Codes: cbc.CodeSet{
				KeyFacturaEInvoiceClass: "CO", // Copy of the original
			},
		},
		{
			Tags:  []cbc.Key{TagCopy},
			Types: []cbc.Key{bill.InvoiceTypeCorrective},
			Codes: cbc.CodeSet{
				KeyFacturaEInvoiceClass: "CR", // Copy of the corrective
			},
		},
		{
			Tags:  []cbc.Key{TagCopy, TagSummary},
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Codes: cbc.CodeSet{
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
