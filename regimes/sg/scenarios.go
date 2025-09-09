// Package sg provides tax scenarios specific to Singapore GST regulations.
package sg

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// Reference: https://www.iras.gov.sg/media/docs/default-source/e-tax/etaxguide_gst_gst-general-guide-for-businesses(1).pdf?sfvrsn=8a66716d_97 (pg 26-27)

// Invoice type tags
const (
	TagInvoiceReceipt cbc.Key = "receipt"
)

func invoiceTags() *tax.TagSet {
	return &tax.TagSet{
		Schema: bill.ShortSchemaInvoice,
		List: []*cbc.Definition{
			{
				Key: TagInvoiceReceipt,
				Name: i18n.String{
					i18n.EN: "Receipt",
				},
			},
		},
	}
}

func invoiceScenarios() *tax.ScenarioSet {
	return &tax.ScenarioSet{
		Schema: bill.ShortSchemaInvoice,
		List: []*tax.Scenario{
			// Reverse Charges
			{
				Tags: []cbc.Key{tax.TagReverseCharge},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  tax.TagReverseCharge,
					Text: "This supply is subject to reverse charge. GST to be accounted for by the recipient.",
				},
			},
			// Simplified Tax Invoice
			{
				Tags: []cbc.Key{tax.TagSimplified},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  tax.TagSimplified,
					Text: "Price Payable includes GST",
				},
			},
			// Receipt
			{
				Tags: []cbc.Key{TagInvoiceReceipt},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  TagInvoiceReceipt,
					Text: "Price Payable includes GST",
				},
			},
			// Self-billed
			{
				Tags: []cbc.Key{tax.TagSelfBilled},
				Note: &tax.ScenarioNote{
					Key:  org.NoteKeyLegal,
					Src:  tax.TagSelfBilled,
					Text: "Self-billed",
				},
			},
		},
	}
}
