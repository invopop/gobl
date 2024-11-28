// Package in provides tax scenarios specific to India GST regulations.
package in

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// Tax tags that can be applied in India.
const (
	TagBillOfSupply           cbc.Key = "bill-of-supply"
	TagInvoiceCumBillOfSupply cbc.Key = "invoice-cum-bill-of-supply"
)

var invoiceTags = &tax.TagSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*cbc.KeyDefinition{
		{
			Key: TagBillOfSupply,
			Name: i18n.String{
				i18n.EN: "Bill of Supply",
				i18n.HI: "आपूर्ति का बिल",
			},
		},
		{
			Key: TagInvoiceCumBillOfSupply,
			Name: i18n.String{
				i18n.EN: "Invoice-cum-bill of supply",
				i18n.HI: "चालान-सह-आपूर्ति का बिल",
			},
		},
	},
}

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		// Reverse Charges
		{
			Tags: []cbc.Key{tax.TagReverseCharge},
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  tax.TagReverseCharge,
				Text: "Reverse Charge",
			},
		},
		// Simplified Tax Invoice
		{
			Tags: []cbc.Key{tax.TagSimplified},
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  tax.TagSimplified,
				Text: "Simplified Tax Invoice",
			},
		},
		// Bill of Supply
		{
			Tags: []cbc.Key{TagBillOfSupply},
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  TagBillOfSupply,
				Text: "Bill Of Supply",
			},
		},
		// Invoice-cum-bill of Supply
		{
			Tags: []cbc.Key{TagInvoiceCumBillOfSupply},
			Note: &cbc.Note{
				Key:  cbc.NoteKeyLegal,
				Src:  TagInvoiceCumBillOfSupply,
				Text: "Invoice-cum-bill Of Supply",
			},
		},
	},
}
