// Package ae provides tax scenarios specific to UAE VAT regulations.
package no

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
				Text: "Reverse Charge",
			},
		},
		// Simplified Tax Invoice
		{
			Tags: []cbc.Key{tax.TagSimplified},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  tax.TagSimplified,
				Text: "Simplified Tax Invoice",
			},
		},
		// Simplified Invoice
		{
			Tags: []cbc.Key{tax.TagSimplified},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  tax.TagSimplified,
				Text: "Simplified Invoice (for transactions below NOK 1,000)",
			},
		},
		// Zero-Rated Export
		{
			Tags: []cbc.Key{tax.RateZero},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  tax.RateZero,
				Text: "Zero-Rated Export",
			},
		},
	},
}
