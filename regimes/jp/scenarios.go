package jp

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		// Reverse Charge
		{
			Tags: []cbc.Key{tax.TagReverseCharge},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  tax.TagReverseCharge,
				Text: "Reverse Charge: The recipient is liable for the consumption tax on this transaction",
			},
		},
		// Simplified Invoice
		{
			Tags: []cbc.Key{tax.TagSimplified},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  tax.TagSimplified,
				Text: "Simplified Invoice",
			},
		},
	},
}
