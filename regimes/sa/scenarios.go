package sa

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		// Simplified Tax Invoice
		{
			Tags: []cbc.Key{tax.TagSimplified},
			Note: &tax.Note{
				Category: tax.CategoryVAT,
				Key:      tax.TagSimplified,
				Text:     "Simplified Tax Invoice",
			},
		},
	},
}
