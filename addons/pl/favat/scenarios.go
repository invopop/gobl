package favat

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// FA_VAT specific tags for invoice types
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
		{
			Key: cbc.Key("exempt"),
			Name: i18n.String{
				i18n.EN: "Tax Exempt",
				i18n.PL: "Zwolnienie z VAT",
			},
			Desc: i18n.String{
				i18n.EN: "Marks invoices that are exempt from VAT and requires the pl-favat-exemption code and note.",
				i18n.PL: "Oznacza faktury zwolnione z VAT i wymaga kodu oraz notatki pl-favat-exemption.",
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
			Ext: tax.Extensions{
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
			Ext: tax.Extensions{
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
			Ext: tax.Extensions{
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
			Ext: tax.Extensions{
				ExtKeyInvoiceType: "UPR",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
			Name: i18n.String{
				i18n.EN: "Credit note",
				i18n.PL: "Faktura korygująca",
			},
			Ext: tax.Extensions{
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
			Ext: tax.Extensions{
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
			Ext: tax.Extensions{
				ExtKeyInvoiceType: "KOR_ROZ",
			},
		},
	},
}
