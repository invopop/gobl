package nz

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
			Note: &tax.Note{
				Key:  tax.KeyReverseCharge,
				Text: "Reverse charge: Recipient to account for GST to Inland Revenue.",
			},
		},
	},
}
