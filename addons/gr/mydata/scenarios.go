package mydata

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// Document tag keys
const (
	TagIslands  cbc.Key = "islands"
	TagGoods    cbc.Key = "goods"
	TagServices cbc.Key = "services"
	TagExport   cbc.Key = "export"
	TagEU       cbc.Key = "eu"
)

var invoiceTags = &tax.TagSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*cbc.Definition{
		{
			Key: TagIslands,
			Name: i18n.String{
				i18n.EN: "Islands Reduced Rates",
				i18n.EL: "Νησιά μειωμένοι συντελεστές",
			},
		},
		{
			Key: TagGoods,
			Name: i18n.String{
				i18n.EN: "Goods",
			},
		},
		{
			Key: TagServices,
			Name: i18n.String{
				i18n.EN: "Services",
			},
		},
		{
			Key: TagExport,
			Name: i18n.String{
				i18n.EN: "Export",
			},
		},
		{
			Key: TagEU,
			Name: i18n.String{
				i18n.EN: "European Union",
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
		// ** Invoice Types **
		{
			// Default: if no other, more specific scenario matches,
			// a services invoice is assumed.
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Ext: tax.Extensions{
				ExtKeyInvoiceType: "2.1",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{TagGoods},
			Ext: tax.Extensions{
				ExtKeyInvoiceType: "1.1",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{TagGoods, TagExport},
			Ext: tax.Extensions{
				ExtKeyInvoiceType: "1.3",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{TagGoods, TagExport, TagEU},
			Ext: tax.Extensions{
				ExtKeyInvoiceType: "1.2",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{TagGoods, tax.TagSelfBilled},
			Ext: tax.Extensions{
				ExtKeyInvoiceType: "1.4",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{TagServices},
			Ext: tax.Extensions{
				ExtKeyInvoiceType: "2.1",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{TagServices, TagExport},
			Ext: tax.Extensions{
				ExtKeyInvoiceType: "2.3",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{TagServices, TagExport, TagEU},
			Ext: tax.Extensions{
				ExtKeyInvoiceType: "2.2",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
			Ext: tax.Extensions{
				ExtKeyInvoiceType: "5.1",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{tax.TagSimplified},
			Ext: tax.Extensions{
				ExtKeyInvoiceType: "11.3",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{TagGoods, tax.TagSimplified},
			Ext: tax.Extensions{
				ExtKeyInvoiceType: "11.1",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{TagServices, tax.TagSimplified},
			Ext: tax.Extensions{
				ExtKeyInvoiceType: "11.2",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
			Tags:  []cbc.Key{tax.TagSimplified},
			Ext: tax.Extensions{
				ExtKeyInvoiceType: "11.4",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
			Tags:  []cbc.Key{TagGoods, tax.TagSimplified, tax.TagSelfBilled},
			Ext: tax.Extensions{
				ExtKeyInvoiceType: "11.5",
			},
		},
	},
}
