package favat

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// Document tag keys
const (
	TagSettlement cbc.Key = "settlement"
)

var invoiceTags = &tax.TagSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*cbc.Definition{
		{
			Key: TagSettlement,
			Name: i18n.String{
				i18n.EN: "Settlement Invoice",
				i18n.PL: "Faktura Rozliczeniowa",
			},
		},
	},
}

var scenarios = []*tax.ScenarioSet{
	invoiceScenarios,
}

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		// **** Invoice Type ****
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Name: i18n.String{
				i18n.EN: "Regular Invoice",
				i18n.PL: "Faktura Podstawowa",
			},
			Codes: cbc.CodeMap{
				ExtKeyInvoiceType: "VAT",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{tax.TagPartial},
			Name: i18n.String{
				i18n.EN: "Prepayment Invoice",
				i18n.PL: `Faktura Zaliczkowa`,
			},
			Codes: cbc.CodeMap{
				ExtKeyInvoiceType: "ZAL",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{TagSettlement},
			Name: i18n.String{
				i18n.EN: "Settlement Invoice",
				i18n.PL: "Faktura Rozliczeniowa",
			},
			Codes: cbc.CodeMap{
				ExtKeyInvoiceType: "ROZ",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{tax.TagSimplified},
			Name: i18n.String{
				i18n.EN: "Simplified Invoice",
				i18n.PL: "Faktura Uproszczona",
			},
			Codes: cbc.CodeMap{
				ExtKeyInvoiceType: "UPR",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
			Name: i18n.String{
				i18n.EN: "Credit note",
				i18n.PL: "Faktura korygująca",
			},
			Codes: cbc.CodeMap{
				ExtKeyInvoiceType: "KOR",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
			Tags:  []cbc.Key{tax.TagPartial},
			Name: i18n.String{
				i18n.EN: "Prepayment credit note",
				i18n.PL: `Faktura korygująca fakturę zaliczkową`,
			},
			Codes: cbc.CodeMap{
				ExtKeyInvoiceType: "KOR_ZAL",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
			Tags:  []cbc.Key{TagSettlement},
			Name: i18n.String{
				i18n.EN: "Settlement credit note",
				i18n.PL: "Faktura korygująca fakturę rozliczeniową",
			},
			Codes: cbc.CodeMap{
				ExtKeyInvoiceType: "KOR_ROZ",
			},
		},
	},
}
