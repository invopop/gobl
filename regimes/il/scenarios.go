// Package il provides tax scenarios specific to Israel VAT regulations.
package il

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
		// Source: Israeli VAT Law 1975, Section 20.
		{
			Tags: []cbc.Key{tax.TagReverseCharge},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  tax.TagReverseCharge,
				Text: "Reverse Charge",
			},
		},
		// Simplified Tax Invoice
		// Source: Israeli VAT Law 1975, Section 46.
		{
			Tags: []cbc.Key{tax.TagSimplified},
			Note: &tax.ScenarioNote{
				Key:  org.NoteKeyLegal,
				Src:  tax.TagSimplified,
				Text: "Simplified Tax Invoice",
			},
		},
	},
}
