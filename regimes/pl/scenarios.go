package pl

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

var invoiceTags = []*tax.KeyDefinition{
	{
		Key: TagSettlement,
		Name: i18n.String{
			i18n.EN: "Settlement Invoice",
			i18n.PL: "Faktura Rozliczeniowa",
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
				KeyFA_VATInvoiceType: "VAT",
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
				KeyFA_VATInvoiceType: "ZAL",
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
				KeyFA_VATInvoiceType: "ROZ",
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
				KeyFA_VATInvoiceType: "UPR",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCorrective},
			Name: i18n.String{
				i18n.EN: "Corrective Invoice",
				i18n.PL: "Faktura Korygująca",
			},
			Codes: cbc.CodeMap{
				KeyFA_VATInvoiceType: "KOR",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCorrective},
			Tags:  []cbc.Key{tax.TagPartial},
			Name: i18n.String{
				i18n.EN: "Corrective Prepayment Invoice",
				i18n.PL: `Faktura korygująca fakturę zaliczkową`,
			},
			Codes: cbc.CodeMap{
				KeyFA_VATInvoiceType: "KOR_ZAL",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCorrective},
			Tags:  []cbc.Key{TagSettlement},
			Name: i18n.String{
				i18n.EN: "Corrective Settlement Invoice",
				i18n.PL: "Faktura korygująca fakturę rozliczeniową",
			},
			Codes: cbc.CodeMap{
				KeyFA_VATInvoiceType: "KOR_ROZ",
			},
		},
	},
}
