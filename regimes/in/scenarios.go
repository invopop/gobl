// Package in provides tax scenarios specific to India GST regulations.
package in

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

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
	},
}
