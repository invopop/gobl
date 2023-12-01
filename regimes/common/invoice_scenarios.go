package common

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
				Text: "Reverse charge: Customer to account for VAT to the relevant tax authority.",
			},
		},
	},
}

// InvoiceScenarios provides a standard set of scenarios to either be extended
// or overridden by the regime.
func InvoiceScenarios() *tax.ScenarioSet {
	return invoiceScenarios
}
