package no

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// invoiceScenarios defines Norwegian-specific invoice scenarios that inject
// legal notes based on tags.
var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		// Reverse charge
		{
			Tags: []cbc.Key{tax.TagReverseCharge},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  tax.TagReverseCharge,
				Text: "Reverse charge / Omvendt avgiftsplikt – Merverdiavgift ikke beregnet.",
			},
		},
	},
}
