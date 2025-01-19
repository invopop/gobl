// Package sg provides tax scenarios specific to Singapore GST regulations.
package sg

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

var invoiceScenarios = &tax.ScenarioSet{
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

		// Simplified Tax Invoice or Reciept
		{
			Tags: []cbc.Key{tax.TagSimplified},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  tax.TagSimplified,
				Text: "Price Payable includes GST",
			},
		},
		{
			Tags: []cbc.Key{TagInvoiceReceipt},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  TagInvoiceReceipt,
				Text: "Price Payable includes GST",
			},
		},
	},
}
